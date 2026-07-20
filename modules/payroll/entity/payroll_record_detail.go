package entity

import (
	"time"

	"github.com/google/uuid"
)

type PayrollRecordDetail struct {
	ID          uuid.UUID `db:"id" json:"id"`
	RecordID    uuid.UUID `db:"record_id" json:"record_id"`
	Component   string    `db:"component" json:"component"`
	Description string    `db:"description" json:"description"`
	Source      string    `db:"source" json:"source"`
	Amount      float64   `db:"amount" json:"amount"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
