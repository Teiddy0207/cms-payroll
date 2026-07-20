package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// SeedDepartments seeds initial department data
func SeedDepartments(ctx context.Context, db database.Database) error {
	departments := []struct {
		Code        string
		Name        string
		Description string
	}{
		{Code: "BOD", Name: "Ban Giám đốc", Description: "Ban điều hành công ty"},
		{Code: "IT", Name: "Phòng Công nghệ thông tin", Description: "Phát triển và vận hành hệ thống phần mềm"},
		{Code: "HR", Name: "Phòng Nhân sự", Description: "Tuyển dụng, đào tạo và quản lý nhân sự"},
		{Code: "ACC", Name: "Phòng Kế toán", Description: "Kế toán, tài chính và tiền lương"},
		{Code: "SALES", Name: "Phòng Kinh doanh", Description: "Bán hàng và phát triển khách hàng"},
		{Code: "MKT", Name: "Phòng Marketing", Description: "Truyền thông và marketing"},
		{Code: "ADMIN", Name: "Phòng Hành chính", Description: "Hành chính, văn phòng"},
	}

	successCount := 0
	errorCount := 0
	skipCount := 0

	for _, dept := range departments {
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM departments WHERE code = $1 LIMIT 1`, dept.Code)
		if err == nil && existingID != uuid.Nil {
			skipCount++
			continue
		} else if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedDepartments: Error checking existing department", "code", dept.Code, "error", err)
			errorCount++
			continue
		}

		deptID := uuid.New()

		query := `INSERT INTO departments (id, code, name, description, created_at, updated_at)
		         VALUES ($1, $2, $3, $4, NOW(), NOW())`

		err = db.ExecContext(ctx, query, deptID, dept.Code, dept.Name, dept.Description)
		if err != nil {
			logger.Error("SeedDepartments: Failed to create department",
				"code", dept.Code,
				"name", dept.Name,
				"id", deptID,
				"error", err)
			errorCount++
			continue
		}

		successCount++
	}

	if errorCount > 0 {
		return fmt.Errorf("seed departments completed with %d errors", errorCount)
	}

	logger.Info("Seeded departments successfully!", "added", successCount, "skipped", skipCount)
	return nil
}
