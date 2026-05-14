package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/model"
)

type AlertRepository struct {
	db *sqlx.DB
}

func NewAlertRepository(db *sqlx.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) List(limit, offset int) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.Select(&alerts, "SELECT * FROM alerts ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	return alerts, err
}

func (r *AlertRepository) Create(alert *model.Alert) error {
	result, err := r.db.Exec(`INSERT INTO alerts
		(rule_id, hostname, severity, metric, value, threshold, message, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.RuleID, alert.Hostname, alert.Severity, alert.Metric,
		alert.Value, alert.Threshold, alert.Message, alert.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	alert.ID = id
	return nil
}

func (r *AlertRepository) Resolve(id int64) error {
	_, err := r.db.Exec("UPDATE alerts SET status='resolved', resolved_at=NOW() WHERE id=?", id)
	return err
}
