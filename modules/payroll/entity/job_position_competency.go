package entity

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeCompetency struct {
	UserProfileID uuid.UUID `db:"user_profile_id" json:"user_profile_id"`
	CompetencyID  uuid.UUID `db:"competency_id" json:"competency_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type EmployeeCompetencyDetail struct {
	UserProfileID  uuid.UUID `db:"user_profile_id" json:"user_profile_id"`
	CompetencyID   uuid.UUID `db:"competency_id" json:"competency_id"`
	CompetencyName string    `db:"competency_name" json:"competency_name"`
	CompetencyCode string    `db:"competency_code" json:"competency_code"`
	PointValue     int       `db:"point_value" json:"point_value"`
}
