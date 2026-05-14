package com.pingan.monitor;

import com.pingan.monitor.deserializer.MetricDeserializer;
import com.pingan.monitor.function.AlertMapper;
import com.pingan.monitor.function.RuleEvaluator;
import com.pingan.monitor.model.AlertEvent;
import com.pingan.monitor.model.MetricEvent;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.base.DeliveryGuarantee;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

import java.time.Duration;

public class MonitoringJob {

    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        // Enable checkpointing for state durability (lightweight, every 60s)
        env.enableCheckpointing(60_000);

        String kafkaBootstrap = getEnv("KAFKA_BROKERS", "localhost:9092");
        String inputTopic      = getEnv("KAFKA_INPUT_TOPIC", "metrics");
        String outputTopic     = getEnv("KAFKA_OUTPUT_TOPIC", "alerts");
        String redisHost       = getEnv("REDIS_HOST", "localhost");
        int redisPort          = Integer.parseInt(getEnv("REDIS_PORT", "6379"));

        // ---- Source: Kafka with 5s bounded-out-of-orderness watermark ----
        KafkaSource<MetricEvent> source = KafkaSource.<MetricEvent>builder()
                .setBootstrapServers(kafkaBootstrap)
                .setTopics(inputTopic)
                .setGroupId("flink-monitor")
                .setStartingOffsets(OffsetsInitializer.latest())
                .setValueOnlyDeserializer(new MetricDeserializer())
                .build();

        DataStream<MetricEvent> metrics = env
                .fromSource(source, WatermarkStrategy.noWatermarks(), "kafka-metrics");

        // ---- Process: rule evaluation with pub/sub + duration + dedup ----
        DataStream<AlertEvent> alerts = metrics
                .keyBy(MetricEvent::getHostname)
                .process(new RuleEvaluator(redisHost, redisPort))
                .name("rule-evaluator")
                .uid("rule-evaluator");

        // ---- Map to JSON string ----
        DataStream<String> jsonAlerts = alerts
                .map(new AlertMapper())
                .name("alert-mapper")
                .uid("alert-mapper");

        // ---- Sink: Kafka (AT_MOST_ONCE for low latency) ----
        KafkaSink<String> sink = KafkaSink.<String>builder()
                .setBootstrapServers(kafkaBootstrap)
                .setRecordSerializer(
                        KafkaRecordSerializationSchema.<String>builder()
                                .setTopic(outputTopic)
                                .setValueSerializationSchema(
                                        (element, context) -> element.getBytes()
                                )
                                .build()
                )
                .setDeliveryGuarantee(DeliveryGuarantee.AT_MOST_ONCE)
                .build();

        jsonAlerts.sinkTo(sink)
                .name("alert-sink")
                .uid("alert-sink");

        env.execute("pingan-monitoring-job");
    }

    private static String getEnv(String key, String fallback) {
        String v = System.getenv(key);
        return v != null ? v : fallback;
    }
}
