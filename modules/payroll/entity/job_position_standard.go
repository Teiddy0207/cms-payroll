package entity

import (
	"time"

	"github.com/google/uuid"
)

type JobPositionStandard struct {
	JobDescriptionID uuid.UUID `db:"job_description_id" json:"job_description_id"`
	JobStandardID    uuid.UUID `db:"job_standard_id" json:"job_standard_id"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type JobPositionStandardDetail struct {
	JobDescriptionID uuid.UUID `db:"job_description_id" json:"job_description_id"`
	JobStandardID    uuid.UUID `db:"job_standard_id" json:"job_standard_id"`
	StandardCode     string    `db:"standard_code" json:"standard_code"`
	StandardName     string    `db:"standard_name" json:"standard_name"`
	AllowanceValue   float64   `db:"allowance_value" json:"allowance_value"`
}
