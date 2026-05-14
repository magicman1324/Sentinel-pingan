package com.pingan.monitor.model;

public class AlertEvent {

    private long ruleId;
    private String ruleName;
    private String hostname;
    private String severity;
    private String metric;
    private double value;
    private double threshold;
    private String message;
    private long timestamp;

    public AlertEvent() {}

    public AlertEvent(long ruleId, String ruleName, String hostname, String severity,
                      String metric, double value, double threshold, String message, long timestamp) {
        this.ruleId = ruleId;
        this.ruleName = ruleName;
        this.hostname = hostname;
        this.severity = severity;
        this.metric = metric;
        this.value = value;
        this.threshold = threshold;
        this.message = message;
        this.timestamp = timestamp;
    }

    // Getters and setters
    public long getRuleId() { return ruleId; }
    public void setRuleId(long ruleId) { this.ruleId = ruleId; }
    public String getRuleName() { return ruleName; }
    public void setRuleName(String ruleName) { this.ruleName = ruleName; }
    public String getHostname() { return hostname; }
    public void setHostname(String hostname) { this.hostname = hostname; }
    public String getSeverity() { return severity; }
    public void setSeverity(String severity) { this.severity = severity; }
    public String getMetric() { return metric; }
    public void setMetric(String metric) { this.metric = metric; }
    public double getValue() { return value; }
    public void setValue(double value) { this.value = value; }
    public double getThreshold() { return threshold; }
    public void setThreshold(double threshold) { this.threshold = threshold; }
    public String getMessage() { return message; }
    public void setMessage(String message) { this.message = message; }
    public long getTimestamp() { return timestamp; }
    public void setTimestamp(long timestamp) { this.timestamp = timestamp; }
}
