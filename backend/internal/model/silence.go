package model

import "time"

type Silence struct {
	ID        int64     `db:"id" json:"id"`
	Matchers  string    `db:"matchers" json:"matchers"` // JSON: [{"name":"hostname","value":"prod-.*","isRegex":true}]
	StartsAt  time.Time `db:"starts_at" json:"starts_at"`
	EndsAt    time.Time `db:"ends_at" json:"ends_at"`
	Comment   string    `db:"comment" json:"comment"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
