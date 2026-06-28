package dto

import "github.com/google/uuid"

type AssignCompetencyRequest struct {
	CompetencyID  uuid.UUID `json:"competency_id" validate:"required"`
	RequiredLevel int       `json:"required_level" validate:"required,min=1"`
	Weight        int       `json:"weight" validate:"required,min=1"`
}

type JobPositionCompetencyResponse struct {
	JobDescriptionID uuid.UUID `json:"job_description_id"`
	CompetencyID     uuid.UUID `json:"competency_id"`
	CompetencyName   string    `json:"competency_name"`
	CompetencyCode   string    `json:"competency_code"`
	RequiredLevel    int       `json:"required_level"`
	Weight           int       `json:"weight"`
	PointValue       int       `json:"point_value"`
}
