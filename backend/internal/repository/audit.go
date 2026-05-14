package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/model"
)

type AuditRepository struct {
	db *sqlx.DB
}

func NewAuditRepository(db *sqlx.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Log(action model.AuditAction, resource string, id int64, detail string) error {
	_, err := r.db.Exec(
		"INSERT INTO audit_logs (action, resource, resource_id, detail) VALUES (?, ?, ?, ?)",
		string(action), resource, id, detail,
	)
	return err
}
