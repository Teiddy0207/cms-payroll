package dto

import (
	coredto "cal-salary/core/dto"

	"github.com/google/uuid"
)

type JobStandardBreakdown struct {
	StandardID     uuid.UUID `json:"standard_id"`
	StandardCode   string    `json:"standard_code"`
	StandardName   string    `json:"standard_name"`
	AllowanceValue float64   `json:"allowance_value"`
}

type CompetencyScoreBreakdown struct {
	CompetencyID   uuid.UUID `json:"competency_id"`
	CompetencyName string    `json:"competency_name"`
	CompetencyCode string    `json:"competency_code"`
	PointValue     int       `json:"point_value"`
	IsAchieved     bool      `json:"is_achieved"`
	EarnedPoints   int       `json:"earned_points"`
}

type SalaryPreviewResponse struct {
	EmployeeID     uuid.UUID                  `json:"employee_id"`
	FullName       string                     `json:"full_name"`
	Period         string                     `json:"period"`
	P1Score        float64                    `json:"p1_score"`
	P1Total        float64                    `json:"p1_total"`
	P2Score        float64                    `json:"p2_score"`
	P2Total        float64                    `json:"p2_total"`
	SystemRate     float64                    `json:"system_rate"`
	SubtotalP1P2   float64                    `json:"subtotal_p1_p2"`
	P1Standards    []JobStandardBreakdown     `json:"p1_standards"`
	P2Competencies []CompetencyScoreBreakdown `json:"p2_competencies"`
	Note           string                     `json:"note"`

	// Compatibility fields for frontend
	P1          float64        `json:"p1"`
	P2          float64        `json:"p2"`
	Total       float64        `json:"total"`
	SalaryP1    float64        `json:"salary_p1"`
	SalaryP2    float64        `json:"salary_p2"`
	TotalSalary float64        `json:"total_salary"`
	Breakdown   map[string]any `json:"breakdown,omitempty"`
}

type SalaryCalculationRequest struct {
	Period        string      `json:"period" validate:"required"`
	EmployeeIDs   []uuid.UUID `json:"employee_ids,omitempty"`
	DepartmentIDs []uuid.UUID `json:"department_ids,omitempty"`
}

type SalaryCalculationItem struct {
	EmployeeID   uuid.UUID              `json:"employee_id"`
	FullName     string                 `json:"full_name"`
	DepartmentID *uuid.UUID             `json:"department_id,omitempty"`
	PositionID   *uuid.UUID             `json:"position_id,omitempty"`
	Status       string                 `json:"status"`
	Error        string                 `json:"error,omitempty"`
	Preview      *SalaryPreviewResponse `json:"preview,omitempty"`
}

type SalaryCalculationResponse struct {
	Period           string                  `json:"period"`
	TotalEmployees   int                     `json:"total_employees"`
	SuccessEmployees int                     `json:"success_employees"`
	FailedEmployees  int                     `json:"failed_employees"`
	Items            []SalaryCalculationItem `json:"items"`
	Note             string                  `json:"note"`
}

type SalaryCalculationItemPagination = coredto.Pagination[SalaryCalculationItem]
