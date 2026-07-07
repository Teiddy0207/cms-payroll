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
