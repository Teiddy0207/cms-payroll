package entity

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeFaceTemplate struct {
	ID           uuid.UUID `db:"id"`
	EmployeeCode string    `db:"employee_code"`
	FaceData     string    `db:"face_data"`
	CreatedAt    time.Time `db:"created_at"`
}
