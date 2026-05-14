package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/model"
)

type SilenceRepository struct {
	db *sqlx.DB
}

func NewSilenceRepository(db *sqlx.DB) *SilenceRepository {
	return &SilenceRepository{db: db}
}

func (r *SilenceRepository) ListActive() ([]model.Silence, error) {
	var silences []model.Silence
	err := r.db.Select(&silences,
		"SELECT * FROM silences WHERE ends_at > NOW() ORDER BY starts_at DESC")
	return silences, err
}

func (r *SilenceRepository) Create(s *model.Silence) error {
	result, err := r.db.Exec(
		"INSERT INTO silences (matchers, starts_at, ends_at, comment, created_by) VALUES (?, ?, ?, ?, ?)",
		s.Matchers, s.StartsAt, s.EndsAt, s.Comment, s.CreatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	s.ID = id
	return nil
}

func (r *SilenceRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM silences WHERE id = ?", id)
	return err
}
