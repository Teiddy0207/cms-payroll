package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type DBUserProfile struct {
	ID   uuid.UUID `db:"id"`
	Code string    `db:"code"`
}

type DBDepartment struct {
	ID   uuid.UUID `db:"id"`
	Code string    `db:"code"`
}

type DBJobPosition struct {
	ID        uuid.UUID `db:"id"`
	Code      string    `db:"code"`
	MinSalary float64   `db:"min_salary"`
	MaxSalary float64   `db:"max_salary"`
}

type DBCompetency struct {
	ID   uuid.UUID `db:"id"`
	Code string    `db:"code"`
}

func SeedPayrollData(ctx context.Context, db database.Database) error {
	logger.Info("Starting SeedPayrollData...")

	// 1. Lấy danh sách phòng ban
	var depts []DBDepartment
	err := db.SelectContext(ctx, &depts, `SELECT id, code FROM departments`)
	if err != nil {
		return fmt.Errorf("failed to fetch departments: %w", err)
	}
	deptMap := make(map[string]uuid.UUID)
	for _, d := range depts {
		deptMap[d.Code] = d.ID
	}

	// 2. Lấy danh sách vị trí công việc
	var positions []DBJobPosition
	err = db.SelectContext(ctx, &positions, `SELECT id, code, min_salary, max_salary FROM job_descriptions`)
	if err != nil {
		return fmt.Errorf("failed to fetch job positions: %w", err)
	}
	posMap := make(map[string]DBJobPosition)
	for _, p := range positions {
		posMap[p.Code] = p
	}

	// 3. Lấy danh sách nhân sự
	var profiles []DBUserProfile
	err = db.SelectContext(ctx, &profiles, `SELECT id, code FROM user_profiles`)
	if err != nil {
		return fmt.Errorf("failed to fetch user profiles: %w", err)
	}

	if len(profiles) == 0 {
		logger.Warn("SeedPayrollData: No user profiles found. Skipping payroll seeding.")
		return nil
	}

	// Xóa dữ liệu cũ của contracts và employee_competencies để seed mới sạch sẽ
	_ = db.ExecContext(ctx, `DELETE FROM contracts`)
	_ = db.ExecContext(ctx, `DELETE FROM employee_competencies`)

	r := rand.New(rand.NewSource(42)) // Dùng seed cố định để dữ liệu ổn định

	// 4. Cập nhật phòng ban, vị trí, tạo hợp đồng lao động cho từng nhân sự
	for _, p := range profiles {
		var deptCode string
		var posCode string

		// Phân bổ phòng ban và vị trí dựa theo mã nhân sự
		switch p.Code {
		case "2510", "2529", "2519", "2592":
			deptCode = "IT"
			posCode = "BACKEND_DEV"
		case "2506", "2897", "2892", "1908":
			deptCode = "IT"
			posCode = "FRONTEND_DEV"
		case "2825", "2622", "3656":
			deptCode = "IT"
			posCode = "PROJECT_MGR"
		case "3442", "3453", "3908", "3968":
			deptCode = "ACC"
			posCode = "ACCOUNTANT"
		case "4005", "4232", "4456", "4461":
			deptCode = "HR"
			posCode = "ACCOUNTANT"
		default:
			deptCode = "SALES"
			posCode = "ACCOUNTANT"
		}

		deptID, hasDept := deptMap[deptCode]
		posInfo, hasPos := posMap[posCode]

		if hasDept && hasPos {
			// Cập nhật vị trí và phòng ban cho Profile
			updateQuery := `UPDATE user_profiles SET department_id = $1, position_id = $2, updated_at = NOW() WHERE id = $3`
			err = db.ExecContext(ctx, updateQuery, deptID, posInfo.ID, p.ID)
			if err != nil {
				logger.Error("SeedPayrollData: Failed to update user profile", "code", p.Code, "error", err)
			}

			// Tạo hợp đồng lao động mẫu cho nhân sự
			contractID := uuid.New()
			contractCode := fmt.Sprintf("HDLD-%s", p.Code)

			// Mức lương cơ bản ngẫu nhiên trong khoảng Min - Max của Vị trí
			baseRate := posInfo.MinSalary + float64(r.Intn(int(posInfo.MaxSalary-posInfo.MinSalary)/1000000))*1000000
			if baseRate <= 0 {
				baseRate = 12000000.0 // Default 12 million
			}

			startDate := "2025-01-01"
			insertContract := `
				INSERT INTO contracts (id, employee_id, contract_code, position_base_rate, start_date, status, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, 'ACTIVE', NOW(), NOW())
			`
			err = db.ExecContext(ctx, insertContract, contractID, p.ID, contractCode, baseRate, startDate)
			if err != nil {
				logger.Error("SeedPayrollData: Failed to insert contract", "code", p.Code, "error", err)
			}
		}
	}

	logger.Info("Seeded departments, job positions and contracts successfully.")

	// 5. Cập nhật point_value cho Competency Dictionaries và gán cho Nhân viên
	var competencies []DBCompetency
	err = db.SelectContext(ctx, &competencies, `SELECT id, code FROM competency_dictionaries`)
	if err == nil && len(competencies) > 0 {
		// Thiết lập điểm point_value mẫu cho từng Competency (từ 500,000đ đến 2,000,000đ)
		for _, comp := range competencies {
			pointVal := 500000 + r.Intn(4)*500000 // 500k, 1M, 1.5M, 2M
			_ = db.ExecContext(ctx, `UPDATE competency_dictionaries SET point_value = $1 WHERE id = $2`, pointVal, comp.ID)
		}

		// Gán 3 năng lực ngẫu nhiên cho mỗi nhân viên
		for _, p := range profiles {
			assignedIndices := make(map[int]bool)
			for i := 0; i < 3; i++ {
				idx := r.Intn(len(competencies))
				if assignedIndices[idx] {
					continue
				}
				assignedIndices[idx] = true
				comp := competencies[idx]

				insertEmpComp := `
					INSERT INTO employee_competencies (user_profile_id, competency_id, created_at)
					VALUES ($1, $2, NOW())
					ON CONFLICT (user_profile_id, competency_id) DO NOTHING
				`
				_ = db.ExecContext(ctx, insertEmpComp, p.ID, comp.ID)
			}
		}
		logger.Info("Seeded competency point values and employee competencies successfully.")
	}

	// 6. Tạo chu kỳ tính lương (payroll_periods) cho tháng 7/2026
	periodID := uuid.New()
	insertPeriod := `
		INSERT INTO payroll_periods (id, month, year, status, created_at, updated_at)
		VALUES ($1, 7, 2026, 'DRAFT', NOW(), NOW())
		ON CONFLICT (month, year) DO NOTHING
	`
	_ = db.ExecContext(ctx, insertPeriod, periodID)
	logger.Info("Seeded payroll period for July 2026.")

	// 7. Cấu hình công thức tính lương mẫu (payroll_formulas)
	var formulaCount int
	_ = db.GetContext(ctx, &formulaCount, `SELECT COUNT(*) FROM payroll_formulas`)
	if formulaCount == 0 {
		formulas := []struct {
			VarName string
			Expr    string
			Desc    string
		}{
			{VarName: "GROSS_SALARY", Expr: "P1 + P2 + P3", Desc: "Tổng thu nhập chịu thuế = P1 + P2 + P3"},
			{VarName: "TAX", Expr: "GROSS_SALARY * 0.1", Desc: "Thuế thu nhập cá nhân 10%"},
			{VarName: "NET_SALARY", Expr: "GROSS_SALARY - TAX", Desc: "Lương thực nhận chuyển khoản = Gross - Thuế"},
		}

		for _, f := range formulas {
			formulaID := uuid.New()
			insertFormula := `
				INSERT INTO payroll_formulas (id, variable_name, expression, start_date, description, created_at, updated_at)
				VALUES ($1, $2, $3, '2026-07-01', $4, NOW(), NOW())
			`
			_ = db.ExecContext(ctx, insertFormula, formulaID, f.VarName, f.Expr, f.Desc)
		}
		logger.Info("Seeded default payroll formulas successfully.")
	}

	logger.Info("SeedPayrollData completed successfully!")
	return nil
}
