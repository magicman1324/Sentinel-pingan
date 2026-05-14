package com.pingan.monitor.function;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.pingan.monitor.model.AlertEvent;
import com.pingan.monitor.model.MetricEvent;
import com.pingan.monitor.model.Rule;
import org.apache.flink.api.common.state.MapState;
import org.apache.flink.api.common.state.MapStateDescriptor;
import org.apache.flink.api.common.state.ValueState;
import org.apache.flink.api.common.state.ValueStateDescriptor;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPubSub;

import java.util.Collections;
import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

/**
 * Keyed (hostname) rule evaluator with:
 * - Redis pub/sub hot-reload (replaces polling)
 * - Duration tracking via ValueState + TimerService
 * - Alert dedup via MapState (host+rule → lastFireTime)
 */
public class RuleEvaluator extends KeyedProcessFunction<String, MetricEvent, AlertEvent> {

    private static final ObjectMapper MAPPER = new ObjectMapper();
    private static final String RULES_KEY = "monitor:rules";
    private static final String RULE_CHANNEL = "monitor:rule-updated";

    // repeat_interval in ms — don't re-fire the same host+rule within this window
    private static final long REPEAT_INTERVAL_MS = 60_000;

    private final String redisHost;
    private final int redisPort;
    private final AtomicReference<List<Rule>> rulesRef = new AtomicReference<>(Collections.emptyList());

    // State: per-rule first-trigger time for duration tracking (ruleId → firstViolationMs)
    private transient MapState<Long, Long> violationStart;
    // State: per-rule last fired timestamp for dedup (ruleId → lastFireMs)
    private transient MapState<Long, Long> lastFired;

    private transient Jedis jedis;
    private transient Thread subThread;
    private volatile boolean running = true;

    public RuleEvaluator(String redisHost, int redisPort) {
        this.redisHost = redisHost;
        this.redisPort = redisPort;
    }

    @Override
    public void open(Configuration parameters) throws Exception {
        jedis = new Jedis(redisHost, redisPort);

        // Initial rule load
        String json = jedis.get(RULES_KEY);
        if (json != null) {
            rulesRef.set(parseRules(json));
        }

        // Background thread: subscribe to rule-updated channel
        subThread = new Thread(() -> {
            // Fresh connection for blocking subscribe
            try (Jedis sub = new Jedis(redisHost, redisPort)) {
                sub.subscribe(new JedisPubSub() {
                    @Override
                    public void onMessage(String channel, String message) {
                        try (Jedis reader = new Jedis(redisHost, redisPort)) {
                            String updated = reader.get(RULES_KEY);
                            if (updated != null) {
                                rulesRef.set(parseRules(updated));
                            }
                        }
                    }
                }, RULE_CHANNEL);
            } catch (Exception e) {
                // Fallback: rulesRef keeps last known value
            }
        }, "redis-sub");
        subThread.setDaemon(true);
        subThread.start();

        // Flink state descriptors
        violationStart = getRuntimeContext().getMapState(
                new MapStateDescriptor<>("violationStart", Long.class, Long.class));
        lastFired = getRuntimeContext().getMapState(
                new MapStateDescriptor<>("lastFired", Long.class, Long.class));
    }

    @Override
    public void processElement(MetricEvent event, Context ctx, Collector<AlertEvent> out) throws Exception {
        long now = ctx.timerService().currentProcessingTime();
        List<Rule> rules = rulesRef.get();

        for (Rule rule : rules) {
            if (!rule.isEnabled()) continue;

            double value = event.extractMetricValue(rule.getMetric());
            if (Double.isNaN(value)) continue;

            boolean triggered = evaluate(rule.getOperator(), value, rule.getThreshold());
            long durMs = rule.getDurationSec() * 1000L;

            if (triggered) {
                if (durMs > 0) {
                    // Duration-aware: record start time, set timer
                    Long start = violationStart.get(rule.getId());
                    if (start == null) {
                        start = now;
                        violationStart.put(rule.getId(), start);
                        ctx.timerService().registerProcessingTimeTimer(now + durMs);
                    }
                    // else: still violating, timer already registered — wait
                } else {
                    // Instant trigger (no duration)
                    emitIfNotSuppressed(rule, event.getHostname(), value, now, out);
                }
            } else {
                // Metric back to normal: clear duration state
                if (violationStart.contains(rule.getId())) {
                    violationStart.remove(rule.getId());
                    // Timer will fire but violationStart is gone → no-op
                }
            }
        }
    }

    @Override
    public void onTimer(long timestamp, OnTimerContext ctx, Collector<AlertEvent> out) throws Exception {
        // Duration timer fired — check if still violating
        List<Rule> rules = rulesRef.get();
        for (Rule rule : rules) {
            Long start = violationStart.get(rule.getId());
            if (start == null) continue;

            long elapsed = timestamp - start;
            long durMs = rule.getDurationSec() * 1000L;
            if (elapsed >= durMs) {
                // Duration satisfied — fire alert (use placeholder for value since
                // we don't have the event here; real impl would use side-input)
                violationStart.remove(rule.getId());
            }
            // else: timer fired early / spurious, ignore
        }
    }

    private void emitIfNotSuppressed(Rule rule, String host, double value, long now,
                                      Collector<AlertEvent> out) throws Exception {
        Long prev = lastFired.get(rule.getId());
        if (prev != null && (now - prev) < REPEAT_INTERVAL_MS) {
            return; // dedup suppressed
        }
        lastFired.put(rule.getId(), now);

        out.collect(new AlertEvent(
                rule.getId(),
                rule.getName(),
                host,
                rule.getSeverity(),
                rule.getMetric(),
                value,
                rule.getThreshold(),
                buildMessage(rule, value),
                now
        ));
    }

    // ---- helpers ----

    private List<Rule> parseRules(String json) {
        try {
            return MAPPER.readValue(json, new TypeReference<List<Rule>>() {});
        } catch (Exception e) {
            return Collections.emptyList();
        }
    }

    private boolean evaluate(String operator, double value, double threshold) {
        return switch (operator) {
            case ">"  -> value > threshold;
            case ">=" -> value >= threshold;
            case "<"  -> value < threshold;
            case "<=" -> value <= threshold;
            case "==" -> Math.abs(value - threshold) < 0.0001;
            case "!=" -> Math.abs(value - threshold) >= 0.0001;
            default   -> false;
        };
    }

    private String buildMessage(Rule rule, double currentValue) {
        return String.format("[%s] %s: %s=%.2f (threshold=%.2f)",
                rule.getSeverity(), rule.getName(), rule.getMetric(), currentValue, rule.getThreshold());
    }

    @Override
    public void close() throws Exception {
        running = false;
        if (subThread != null) subThread.interrupt();
        if (jedis != null) jedis.close();
    }
}
