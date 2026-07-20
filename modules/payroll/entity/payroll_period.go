package entity

import (
	"time"

	"github.com/google/uuid"
)

type PayrollPeriod struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Month       int        `db:"month" json:"month"`
	Year        int        `db:"year" json:"year"`
	Status      string     `db:"status" json:"status"`
	LockedAt    *time.Time `db:"locked_at" json:"locked_at"`
	DisbursedAt *time.Time `db:"disbursed_at" json:"disbursed_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
