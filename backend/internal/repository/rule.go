package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/model"
)

type RuleRepository struct {
	db *sqlx.DB
}

func NewRuleRepository(db *sqlx.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) ListAll() ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Select(&rules, "SELECT * FROM rules ORDER BY id")
	return rules, err
}

func (r *RuleRepository) ListEnabled() ([]model.Rule, error) {
	var rules []model.Rule
	err := r.db.Select(&rules, "SELECT * FROM rules WHERE enabled = 1 ORDER BY id")
	return rules, err
}

func (r *RuleRepository) GetByID(id int64) (*model.Rule, error) {
	var rule model.Rule
	err := r.db.Get(&rule, "SELECT * FROM rules WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RuleRepository) Create(rule *model.Rule) error {
	result, err := r.db.Exec(`INSERT INTO rules
		(name, description, rule_type, metric, operator, threshold, duration_sec, severity, expression, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.Name, rule.Description, rule.RuleType, rule.Metric,
		rule.Operator, rule.Threshold, rule.Duration, rule.Severity, rule.Expression, rule.Enabled)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	rule.ID = id
	return nil
}

func (r *RuleRepository) Update(rule *model.Rule) error {
	_, err := r.db.Exec(`UPDATE rules SET
		name=?, description=?, metric=?, operator=?, threshold=?,
		duration_sec=?, severity=?, expression=?, enabled=?
		WHERE id=?`,
		rule.Name, rule.Description, rule.Metric, rule.Operator,
		rule.Threshold, rule.Duration, rule.Severity, rule.Expression, rule.Enabled, rule.ID)
	return err
}

func (r *RuleRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM rules WHERE id = ?", id)
	return err
}
