package model

import "time"

type AuditAction string

const (
	AuditCreate AuditAction = "create"
	AuditUpdate AuditAction = "update"
	AuditDelete AuditAction = "delete"
)

type AuditLog struct {
	ID         int64       `db:"id" json:"id"`
	Action     AuditAction `db:"action" json:"action"`
	Resource   string      `db:"resource" json:"resource"`
	ResourceID int64       `db:"resource_id" json:"resource_id"`
	Detail     string      `db:"detail" json:"detail"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
}
