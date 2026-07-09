package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// JobPosition DTOs
// ==========================================

type CreateJobPositionRequest struct {
	Code         string     `json:"code" validate:"required,min=2,max=50"`
	Name         string     `json:"name" validate:"required,min=2,max=255"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
	IsBenchmark  bool       `json:"is_benchmark"`
	MarketSalary float64    `json:"market_salary"`
	SearchKeyword string     `json:"search_keyword"`
}

type UpdateJobPositionRequest struct {
	Name         string     `json:"name" validate:"required,min=2,max=255"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
	IsBenchmark  bool       `json:"is_benchmark"`
	MarketSalary float64    `json:"market_salary"`
	SearchKeyword string     `json:"search_keyword"`
}

type JobPositionResponse struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
	JobScore     float64    `json:"job_score"`
	Midpoint     float64    `json:"midpoint"`
	MinSalary    float64    `json:"min_salary"`
	MaxSalary    float64    `json:"max_salary"`
	IsBenchmark  bool       `json:"is_benchmark"`
	MarketSalary float64    `json:"market_salary"`
	SearchKeyword string     `json:"search_keyword"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PaginatedJobPositionDTO = dto.Pagination[JobPositionResponse]

type ScrapeMarketSalaryRequest struct {
	Source  string `json:"source"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
}

type ScrapedJobItem struct {
	Title        string  `json:"title"`
	Company      string  `json:"company"`
	SalaryRange  string  `json:"salary_range"`
	SalaryMin    float64 `json:"salary_min"`
	SalaryMax    float64 `json:"salary_max"`
	SalaryAverage float64 `json:"salary_average"`
}

type ScrapeMarketSalaryResponse struct {
	AverageSalary float64          `json:"average_salary"`
	Jobs          []ScrapedJobItem `json:"jobs"`
	Logs          []string         `json:"logs"`
}

type CalculateKResponse struct {
	NewKFactor    float64               `json:"new_k_factor"`
	BenchmarkJobs []JobPositionResponse `json:"benchmark_jobs"`
}

