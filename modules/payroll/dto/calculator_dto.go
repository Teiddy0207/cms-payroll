package dto

import "github.com/google/uuid"

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
	AchievedLevel  int       `json:"achieved_level"`
	RequiredLevel  int       `json:"required_level"`
	Weight         int       `json:"weight"`
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
}
