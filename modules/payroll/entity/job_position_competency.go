package entity

import (
	"time"

	"github.com/google/uuid"
)

type JobPositionCompetency struct {
	JobDescriptionID uuid.UUID `db:"job_description_id" json:"job_description_id"`
	CompetencyID     uuid.UUID `db:"competency_id" json:"competency_id"`
	RequiredLevel    int       `db:"required_level" json:"required_level"`
	Weight           int       `db:"weight" json:"weight"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type JobPositionCompetencyDetail struct {
	JobDescriptionID uuid.UUID `db:"job_description_id" json:"job_description_id"`
	CompetencyID     uuid.UUID `db:"competency_id" json:"competency_id"`
	CompetencyName   string    `db:"competency_name" json:"competency_name"`
	CompetencyCode   string    `db:"competency_code" json:"competency_code"`
	RequiredLevel    int       `db:"required_level" json:"required_level"`
	Weight           int       `db:"weight" json:"weight"`
	PointValue       int       `db:"point_value" json:"point_value"`
}
