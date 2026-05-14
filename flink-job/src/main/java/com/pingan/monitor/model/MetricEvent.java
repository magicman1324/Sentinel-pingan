package com.pingan.monitor.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Map;

/**
 * Mirrors agent/internal/model/MetricPayload JSON.
 * Flat structure — no "metrics" wrapper:
 * {"hostname":"prod-01","ts":1715700000000,"cpu":{"percent_used":42.5,"cores":8},"memory":{...}}
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class MetricEvent {

    @JsonProperty("hostname")
    private String hostname;

    @JsonProperty("ts")
    private long timestamp;

    @JsonProperty("cpu")
    private Map<String, Object> cpu;

    @JsonProperty("memory")
    private Map<String, Object> memory;

    @JsonProperty("disk")
    private Object disk;

    @JsonProperty("network")
    private Object network;

    public String getHostname()                  { return hostname; }
    public void setHostname(String hostname)     { this.hostname = hostname; }
    public long getTimestamp()                   { return timestamp; }
    public void setTimestamp(long timestamp)     { this.timestamp = timestamp; }
    public Map<String, Object> getCpu()          { return cpu; }
    public Map<String, Object> getMemory()       { return memory; }

    /**
     * Extracts a metric value by dotted path.
     * "cpu.percent_used" → cpu map → percent_used
     * "memory.used_bytes" → memory map → used_bytes
     */
    @SuppressWarnings("unchecked")
    public double extractMetricValue(String metricPath) {
        String[] parts = metricPath.split("\\.", 2);
        if (parts.length != 2) return Double.NaN;

        Map<String, Object> category = switch (parts[0]) {
            case "cpu"    -> cpu;
            case "memory" -> memory;
            default       -> null;
        };

        if (category == null) return Double.NaN;

        Object val = category.get(parts[1]);
        if (val instanceof Number) return ((Number) val).doubleValue();
        return Double.NaN;
    }
}
