package dto

import "time"

type CreateExplanationRequest struct {
	Date   time.Time `json:"date"`
	Reason string    `json:"reason"`
}

type ExplanationRequestResponse struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	EmployeeCode string    `json:"employee_code"`
	FullName     string    `json:"full_name"`
	Date         string    `json:"date"`
	Reason       string    `json:"reason"`
	ApprovedBy   *string   `json:"approved_by,omitempty"`
	ApprovedName *string   `json:"approved_name,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateOTRequest struct {
	Date           time.Time `json:"date"`
	HoursRequested float64   `json:"hours_requested"`
	IsNightOT      bool      `json:"is_night_ot"`
	IsHolidayOT    bool      `json:"is_holiday_ot"`
}

type OTRequestResponse struct {
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employee_id"`
	EmployeeCode   string    `json:"employee_code"`
	FullName       string    `json:"full_name"`
	Date           string    `json:"date"`
	HoursRequested float64   `json:"hours_requested"`
	IsNightOT      bool      `json:"is_night_ot"`
	IsHolidayOT    bool      `json:"is_holiday_ot"`
	ApprovedBy     *string   `json:"approved_by,omitempty"`
	ApprovedName   *string   `json:"approved_name,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateRequestStatus struct {
	Status string `json:"status"` // APPROVED or REJECTED
}
