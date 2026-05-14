package com.pingan.monitor;

import com.pingan.monitor.deserializer.MetricDeserializer;
import com.pingan.monitor.function.RuleEvaluator;
import com.pingan.monitor.function.AlertMapper;
import com.pingan.monitor.model.AlertEvent;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

import java.time.Duration;

public class MonitoringJob {

    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        String kafkaBootstrap = getEnv("KAFKA_BROKERS", "localhost:9092");
        String inputTopic = getEnv("KAFKA_INPUT_TOPIC", "metrics");
        String outputTopic = getEnv("KAFKA_OUTPUT_TOPIC", "alerts");
        String redisHost = getEnv("REDIS_HOST", "localhost");
        int redisPort = Integer.parseInt(getEnv("REDIS_PORT", "6379"));

        // Kafka source for raw metrics
        KafkaSource<com.pingan.monitor.model.MetricEvent> source = KafkaSource
                .<com.pingan.monitor.model.MetricEvent>builder()
                .setBootstrapServers(kafkaBootstrap)
                .setTopics(inputTopic)
                .setGroupId("flink-monitor")
                .setStartingOffsets(OffsetsInitializer.latest())
                .setValueOnlyDeserializer(new MetricDeserializer())
                .build();

        DataStream<com.pingan.monitor.model.MetricEvent> metrics =
                env.fromSource(source, WatermarkStrategy.noWatermarks(), "kafka-metrics");

        // Rule evaluation with Redis-backed hot-reload
        DataStream<AlertEvent> alerts = metrics
                .keyBy(MetricEvent::getHostname)
                .process(new RuleEvaluator(redisHost, redisPort))
                .name("rule-evaluator");

        // Enrich and sink to output Kafka topic
        DataStream<String> output = alerts
                .map(new AlertMapper())
                .name("alert-mapper");

        output.sinkTo(org.apache.flink.connector.base.DeliveryGuarantee.AT_MOST_ONCE,
                org.apache.flink.streaming.api.functions.sink.PrintSinkFunction::new)
                .name("alert-sink");

        env.execute("pingan-monitoring-job");
    }

    private static String getEnv(String key, String fallback) {
        String v = System.getenv(key);
        return v != null ? v : fallback;
    }
}
