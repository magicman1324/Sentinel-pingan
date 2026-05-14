-- Monitor Database Schema for TiDB
CREATE DATABASE IF NOT EXISTS monitor;
USE monitor;

-- Alert rules table
CREATE TABLE IF NOT EXISTS rules (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    rule_type   ENUM('atomic','composite') NOT NULL DEFAULT 'atomic',
    metric      VARCHAR(255) NOT NULL COMMENT 'e.g. cpu.percent_used',
    operator    VARCHAR(10) NOT NULL DEFAULT '>' COMMENT '>, >=, <, <=, ==, !=',
    threshold   DOUBLE NOT NULL,
    duration_sec INT NOT NULL DEFAULT 0 COMMENT '持续秒数',
    severity    ENUM('critical','warning','info') NOT NULL DEFAULT 'warning',
    expression  TEXT COMMENT '复合规则表达式',
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_enabled (enabled)
);

-- Alerts history table
CREATE TABLE IF NOT EXISTS alerts (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    rule_id     BIGINT NOT NULL,
    hostname    VARCHAR(255) NOT NULL,
    severity    ENUM('critical','warning','info') NOT NULL DEFAULT 'warning',
    metric      VARCHAR(255) NOT NULL,
    value       DOUBLE NOT NULL,
    threshold   DOUBLE NOT NULL,
    message     TEXT,
    status      ENUM('firing','resolved') NOT NULL DEFAULT 'firing',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_at DATETIME,
    INDEX idx_status (status),
    INDEX idx_hostname (hostname),
    INDEX idx_created_at (created_at)
);

-- Alert channels configuration
CREATE TABLE IF NOT EXISTS channels (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    type        ENUM('dingtalk','email','sms') NOT NULL,
    config      JSON NOT NULL COMMENT 'channel-specific configuration',
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Relation table: rule <-> channel
CREATE TABLE IF NOT EXISTS rule_channels (
    rule_id    BIGINT NOT NULL,
    channel_id BIGINT NOT NULL,
    PRIMARY KEY (rule_id, channel_id)
);
