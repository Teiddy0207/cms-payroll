package dto

import (
	"time"

	"github.com/google/uuid"
)

type CheckinRequest struct {
	EmployeeCode string    `json:"employee_code"`
	FaceData     string    `json:"face_data"`
	Timestamp    time.Time `json:"timestamp"`
	LocationGPS  string    `json:"location_gps"`
	DeviceId     string    `json:"device_id"`
}

type CheckinResponse struct {
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	EmployeeCode string    `json:"employee_code,omitempty"`
	FullName     string    `json:"full_name,omitempty"`
}

type RegisterFaceRequest struct {
	EmployeeCode string `json:"employee_code"`
	FaceData     string `json:"face_data"`
}

type FaceTemplateResponse struct {
	EmployeeCode string    `json:"employee_code"`
	EmployeeName string    `json:"employee_name"`
	FaceData     string    `json:"face_data"`
	CreatedAt    time.Time `json:"created_at"`
}

type AttendanceLogResponse struct {
	ID             uuid.UUID `json:"id"`
	EmployeeCode   string    `json:"employee_code"`
	FullName       string    `json:"full_name"`
	DepartmentName string    `json:"department_name"`
	Timestamp      time.Time `json:"timestamp"`
	LocationGPS    string    `json:"location_gps"`
	DeviceId       string    `json:"device_id"`
}

type PaginatedAttendanceLogResponse struct {
	Items      []AttendanceLogResponse `json:"items"`
	TotalItems int                     `json:"total_items"`
}
