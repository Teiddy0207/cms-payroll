package entity

import (
	"time"

	"github.com/google/uuid"
)

type PayrollRecord struct {
	ID          uuid.UUID `db:"id" json:"id"`
	PeriodID    uuid.UUID `db:"period_id" json:"period_id"`
	EmployeeID  uuid.UUID `db:"employee_id" json:"employee_id"`
	P1Value     float64   `db:"p1_value" json:"p1_value"`
	P2Value     float64   `db:"p2_value" json:"p2_value"`
	P3Value     float64   `db:"p3_value" json:"p3_value"`
	GrossSalary float64   `db:"gross_salary" json:"gross_salary"`
	Tax         float64   `db:"tax" json:"tax"`
	NetSalary   float64   `db:"net_salary" json:"net_salary"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
