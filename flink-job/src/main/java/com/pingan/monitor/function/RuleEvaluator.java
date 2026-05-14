package com.pingan.monitor.function;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.pingan.monitor.model.AlertEvent;
import com.pingan.monitor.model.MetricEvent;
import com.pingan.monitor.model.Rule;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import redis.clients.jedis.Jedis;

import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

public class RuleEvaluator extends KeyedProcessFunction<String, MetricEvent, AlertEvent> {

    private static final ObjectMapper mapper = new ObjectMapper();
    private static final String RULES_KEY = "monitor:rules";

    private final String redisHost;
    private final int redisPort;
    private transient Jedis jedis;
    private final AtomicReference<List<Rule>> cachedRules = new AtomicReference<>();
    private long lastSync = 0;
    private static final long SYNC_INTERVAL_MS = 5000;

    public RuleEvaluator(String redisHost, int redisPort) {
        this.redisHost = redisHost;
        this.redisPort = redisPort;
    }

    @Override
    public void open(Configuration parameters) {
        jedis = new Jedis(redisHost, redisPort);
    }

    @Override
    public void processElement(MetricEvent event, Context ctx, Collector<AlertEvent> out) throws Exception {
        List<Rule> rules = getRules();
        if (rules == null || rules.isEmpty()) return;

        long now = System.currentTimeMillis();

        for (Rule rule : rules) {
            if (!rule.isEnabled()) continue;

            double value = event.extractMetricValue(rule.getMetric());
            if (Double.isNaN(value)) continue;

            boolean triggered = evaluate(rule.getOperator(), value, rule.getThreshold());

            if (triggered) {
                out.collect(new AlertEvent(
                        rule.getId(),
                        rule.getName(),
                        event.getHostname(),
                        rule.getSeverity(),
                        rule.getMetric(),
                        value,
                        rule.getThreshold(),
                        buildMessage(rule, value),
                        now
                ));
            }
        }
    }

    private List<Rule> getRules() {
        long now = System.currentTimeMillis();
        if (now - lastSync > SYNC_INTERVAL_MS) {
            try {
                String json = jedis.get(RULES_KEY);
                if (json != null) {
                    List<Rule> rules = mapper.readValue(json, new TypeReference<List<Rule>>() {});
                    cachedRules.set(rules);
                }
            } catch (Exception e) {
                // fallback to cached rules
            }
            lastSync = now;
        }
        return cachedRules.get();
    }

    private boolean evaluate(String operator, double value, double threshold) {
        return switch (operator) {
            case ">" -> value > threshold;
            case ">=" -> value >= threshold;
            case "<" -> value < threshold;
            case "<=" -> value <= threshold;
            case "==" -> Math.abs(value - threshold) < 0.0001;
            case "!=" -> Math.abs(value - threshold) >= 0.0001;
            default -> false;
        };
    }

    private String buildMessage(Rule rule, double currentValue) {
        return String.format("[%s] %s: %s=%.2f (threshold=%.2f)",
                rule.getSeverity(), rule.getName(), rule.getMetric(), currentValue, rule.getThreshold());
    }

    @Override
    public void close() {
        if (jedis != null) jedis.close();
    }
}
