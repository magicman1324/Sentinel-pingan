package model

import "time"

type AlertStatus string

const (
	AlertFiring   AlertStatus = "firing"
	AlertResolved AlertStatus = "resolved"
)

type Alert struct {
	ID         int64       `db:"id" json:"id"`
	RuleID     int64       `db:"rule_id" json:"rule_id"`
	Hostname   string      `db:"hostname" json:"hostname"`
	Severity   Severity    `db:"severity" json:"severity"`
	Metric     string      `db:"metric" json:"metric"`
	Value      float64     `db:"value" json:"value"`
	Threshold  float64     `db:"threshold" json:"threshold"`
	Message    string      `db:"message" json:"message"`
	Status     AlertStatus `db:"status" json:"status"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	ResolvedAt *time.Time  `db:"resolved_at" json:"resolved_at,omitempty"`
}
