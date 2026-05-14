package com.pingan.monitor.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public class MetricEvent {

    @JsonProperty("hostname")
    private String hostname;

    @JsonProperty("ts")
    private long timestamp;

    @JsonProperty("metrics")
    private Map<String, Object> metrics;

    public String getHostname() { return hostname; }
    public void setHostname(String hostname) { this.hostname = hostname; }

    public long getTimestamp() { return timestamp; }
    public void setTimestamp(long timestamp) { this.timestamp = timestamp; }

    public Map<String, Object> getMetrics() { return metrics; }
    public void setMetrics(Map<String, Object> metrics) { this.metrics = metrics; }

    @SuppressWarnings("unchecked")
    public double extractMetricValue(String metricPath) {
        String[] parts = metricPath.split("\\.");
        Object current = metrics;
        for (String part : parts) {
            if (current instanceof Map) {
                current = ((Map<String, Object>) current).get(part);
            } else {
                return Double.NaN;
            }
        }
        if (current instanceof Number) {
            return ((Number) current).doubleValue();
        }
        return Double.NaN;
    }
}
