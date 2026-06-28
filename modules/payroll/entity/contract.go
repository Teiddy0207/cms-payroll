package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type Contract struct {
	EmployeeID       uuid.UUID  `db:"employee_id" json:"employee_id"`
	ContractCode     string     `db:"contract_code" json:"contract_code"`
	PositionBaseRate float64    `db:"position_base_rate" json:"position_base_rate"`
	StartDate        time.Time  `db:"start_date" json:"start_date"`
	EndDate          *time.Time `db:"end_date" json:"end_date"`
	Status           string     `db:"status" json:"status"` // ACTIVE, EXPIRED, TERMINATED
	entity.BaseEntity
}

type PaginatedContractEntity = entity.Pagination[Contract]
