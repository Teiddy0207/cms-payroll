package dto

import "github.com/google/uuid"

type AssignStandardRequest struct {
	JobStandardID uuid.UUID `json:"job_standard_id" validate:"required"`
}

type JobPositionStandardResponse struct {
	JobDescriptionID uuid.UUID `json:"job_description_id"`
	JobStandardID    uuid.UUID `json:"job_standard_id"`
	StandardCode     string    `json:"standard_code"`
	StandardName     string    `json:"standard_name"`
	AllowanceValue   float64   `json:"allowance_value"`
}
