package entity

import (
	"time"
)

type SystemSetting struct {
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	Description *string   `db:"description" json:"description"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
