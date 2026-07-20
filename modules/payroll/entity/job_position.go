package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type JobPosition struct {
	Code         string     `db:"code" json:"code"`
	Name         string     `db:"name" json:"name"`
	Description  *string    `db:"description" json:"description"`
	DepartmentID *uuid.UUID `db:"department_id" json:"department_id"`
	EScore       float64    `db:"e_score" json:"e_score"`
	CScore       float64    `db:"c_score" json:"c_score"`
	RScore       float64    `db:"r_score" json:"r_score"`
	WEWeight     float64    `db:"we_weight" json:"we_weight"`
	WCWeight     float64    `db:"wc_weight" json:"wc_weight"`
	WRWeight     float64    `db:"wr_weight" json:"wr_weight"`
	SalarySpread float64    `db:"salary_spread" json:"salary_spread"`
	JobScore     float64    `db:"job_score" json:"job_score"`
	Midpoint     float64    `db:"midpoint" json:"midpoint"`
	MinSalary    float64    `db:"min_salary" json:"min_salary"`
	MaxSalary    float64    `db:"max_salary" json:"max_salary"`
	IsBenchmark  bool       `db:"is_benchmark" json:"is_benchmark"`
	MarketSalary float64    `db:"market_salary" json:"market_salary"`
	SearchKeyword *string    `db:"search_keyword" json:"search_keyword"`
	entity.BaseEntity
}
type PaginatedJobPositionEntity = entity.Pagination[JobPosition]

