package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceLog struct {
	ID           uuid.UUID `db:"id" json:"id"`
	EventID      uuid.UUID `db:"event_id" json:"event_id"`
	EmployeeCode string    `db:"employee_code" json:"employee_code"`
	Timestamp    time.Time `db:"timestamp" json:"timestamp"`
	LocationGPS  string    `db:"location_gps" json:"location_gps"`
	DeviceId     string    `db:"device_id" json:"device_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
