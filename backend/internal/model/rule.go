package model

import "time"

type RuleType string

const (
	RuleAtomic    RuleType = "atomic"
	RuleComposite RuleType = "composite"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

type Rule struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	RuleType    RuleType  `db:"rule_type" json:"rule_type"`
	Metric      string    `db:"metric" json:"metric"`
	Operator    string    `db:"operator" json:"operator"`
	Threshold   float64   `db:"threshold" json:"threshold"`
	Duration    int       `db:"duration_sec" json:"duration_sec"`
	Severity    Severity  `db:"severity" json:"severity"`
	Expression  string    `db:"expression" json:"expression,omitempty"`
	Enabled     bool      `db:"enabled" json:"enabled"`
	ChannelIDs  []string  `db:"-" json:"channel_ids"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
