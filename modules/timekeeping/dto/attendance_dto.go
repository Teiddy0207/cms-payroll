package dto

import "time"

type DailyAttendanceResponse struct {
	ID             string     `json:"id"`
	EmployeeID     string     `json:"employee_id"`
	EmployeeCode   string     `json:"employee_code"`
	FullName       string     `json:"full_name"`
	DepartmentName string     `json:"department_name"`
	Date           string     `json:"date"`
	CheckIn        *time.Time `json:"check_in"`
	CheckOut       *time.Time `json:"check_out"`
	ActualWorkDay  float64    `json:"actual_work_day"`
	OTHours        float64    `json:"ot_hours"`
	Status         string     `json:"status"`
}

type CalculateTimesheetRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
