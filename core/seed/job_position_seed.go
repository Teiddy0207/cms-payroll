package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func SeedJobPositions(ctx context.Context, db database.Database) error {
	jobPositions := []struct {
		ID              string
		Code            string
		Name            string
		Description     string
		EScore          float64
		CScore          float64
		RScore          float64
		WEWeight        float64
		WCWeight        float64
		WRWeight        float64
		SalarySpread    float64
		JobScore        float64
		Midpoint        float64
		MinSalary       float64
		MaxSalary       float64
		IsBenchmark     bool
		MarketSalary    float64
		SearchKeyword   string
		P2BaseAllowance float64
		P2Cap           float64
	}{
		{
			ID:              "11111111-1111-1111-1111-111111111111",
			Code:            "BACKEND_DEV",
			Name:            "Backend Developer (Golang)",
			Description:     "Phát triển hệ thống xử lý phía máy chủ sử dụng ngôn ngữ Go.",
			EScore:          3.0,
			CScore:          3.0,
			RScore:          2.0,
			WEWeight:        0.4,
			WCWeight:        0.4,
			WRWeight:        0.2,
			SalarySpread:    0.3,
			JobScore:        80.0,
			Midpoint:        15000000.0,
			MinSalary:       10000000.0,
			MaxSalary:       20000000.0,
			IsBenchmark:     true,
			MarketSalary:    18000000.0,
			SearchKeyword:   "Golang Backend",
			P2BaseAllowance: 2500000.0,
			P2Cap:           5000000.0,
		},
		{
			ID:              "22222222-2222-2222-2222-222222222222",
			Code:            "FRONTEND_DEV",
			Name:            "Frontend Developer (React)",
			Description:     "Phát triển giao diện web tương tác sử dụng ReactJS.",
			EScore:          3.0,
			CScore:          2.0,
			RScore:          2.0,
			WEWeight:        0.4,
			WCWeight:        0.3,
			WRWeight:        0.3,
			SalarySpread:    0.3,
			JobScore:        70.0,
			Midpoint:        13000000.0,
			MinSalary:       9000000.0,
			MaxSalary:       17000000.0,
			IsBenchmark:     true,
			MarketSalary:    15000000.0,
			SearchKeyword:   "React Frontend",
			P2BaseAllowance: 2000000.0,
			P2Cap:           4000000.0,
		},
		{
			ID:              "33333333-3333-3333-3333-333333333333",
			Code:            "PROJECT_MGR",
			Name:            "Project Manager",
			Description:     "Quản lý tiến độ, nguồn lực và điều phối dự án phát triển phần mềm.",
			EScore:          4.0,
			CScore:          4.0,
			RScore:          3.0,
			WEWeight:        0.3,
			WCWeight:        0.4,
			WRWeight:        0.3,
			SalarySpread:    0.4,
			JobScore:        95.0,
			Midpoint:        25000000.0,
			MinSalary:       18000000.0,
			MaxSalary:       32000000.0,
			IsBenchmark:     true,
			MarketSalary:    28000000.0,
			SearchKeyword:   "Project Manager PM",
			P2BaseAllowance: 4000000.0,
			P2Cap:           8000000.0,
		},
		{
			ID:              "44444444-4444-4444-4444-444444444444",
			Code:            "ACCOUNTANT",
			Name:            "Kế toán viên",
			Description:     "Thực hiện công tác kế toán, quản lý sổ sách và xử lý chứng từ.",
			EScore:          2.0,
			CScore:          3.0,
			RScore:          2.0,
			WEWeight:        0.3,
			WCWeight:        0.4,
			WRWeight:        0.3,
			SalarySpread:    0.25,
			JobScore:        60.0,
			Midpoint:        11000000.0,
			MinSalary:       8000000.0,
			MaxSalary:       14000000.0,
			IsBenchmark:     false,
			MarketSalary:    11000000.0,
			SearchKeyword:   "Kế toán",
			P2BaseAllowance: 1500000.0,
			P2Cap:           3000000.0,
		},
	}

	successCount := 0
	errorCount := 0
	skipCount := 0

	for _, jp := range jobPositions {
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM job_descriptions WHERE code = $1 LIMIT 1`, jp.Code)
		if err == nil && existingID != uuid.Nil {
			skipCount++
			continue
		} else if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedJobPositions: Error checking existing position", "code", jp.Code, "error", err)
			errorCount++
			continue
		}

		posID := uuid.MustParse(jp.ID)

		query := `
			INSERT INTO job_descriptions (
				id, code, name, description,
				e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
				salary_spread, job_score, midpoint, min_salary, max_salary,
				is_benchmark, market_salary, search_keyword,
				p2_base_allowance, p2_cap, created_at, updated_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, NOW(), NOW())
		`

		err = db.ExecContext(ctx, query,
			posID, jp.Code, jp.Name, jp.Description,
			jp.EScore, jp.CScore, jp.RScore, jp.WEWeight, jp.WCWeight, jp.WRWeight,
			jp.SalarySpread, jp.JobScore, jp.Midpoint, jp.MinSalary, jp.MaxSalary,
			jp.IsBenchmark, jp.MarketSalary, jp.SearchKeyword,
			jp.P2BaseAllowance, jp.P2Cap,
		)
		if err != nil {
			logger.Error("SeedJobPositions: Failed to create job position", "code", jp.Code, "error", err)
			errorCount++
			continue
		}

		successCount++
	}

	if errorCount > 0 {
		return fmt.Errorf("seed job positions completed with %d errors", errorCount)
	}

	logger.Info("Seeded job positions successfully!", "added", successCount, "skipped", skipCount)
	return nil
}
