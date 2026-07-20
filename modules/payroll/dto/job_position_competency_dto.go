package dto

import "github.com/google/uuid"

// Request gán mảng năng lực cho một nhân sự
type AssignEmployeeCompetenciesRequest struct {
	CompetencyIDs []uuid.UUID `json:"competency_ids" validate:"required,min=1"`
}

type EmployeeCompetencyResponse struct {
	UserProfileID  uuid.UUID `json:"user_profile_id"`
	CompetencyID   uuid.UUID `json:"competency_id"`
	CompetencyName string    `json:"competency_name"`
	CompetencyCode string    `json:"competency_code"`
	PointValue     int       `json:"point_value"`
}
