package dto

import "time"

type CheckinRequest struct {
	EmployeeCode string    `json:"employee_code"`
	Timestamp    time.Time `json:"timestamp"`
	LocationGPS  string    `json:"location_gps"`
	DeviceId     string    `json:"device_id"`
}

type CheckinResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
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
