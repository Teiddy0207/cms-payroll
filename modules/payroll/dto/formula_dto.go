package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateFormulaRequest struct {
	VariableName string     `json:"variable_name"`
	Expression   string     `json:"expression"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Description  string     `json:"description"`
}

type UpdateFormulaRequest struct {
	VariableName string     `json:"variable_name"`
	Expression   string     `json:"expression"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Description  string     `json:"description"`
}

type FormulaResponse struct {
	ID           uuid.UUID  `json:"id"`
	VariableName string     `json:"variable_name"`
	Expression   string     `json:"expression"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
