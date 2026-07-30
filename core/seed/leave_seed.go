package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SeedLeaveData seeds initial leave requests and balances for employees.
func SeedLeaveData(ctx context.Context, db database.Database) error {
	// 1. Get all employees
	var profiles []Profile
	err := db.SelectContext(ctx, &profiles, `SELECT id, code FROM user_profiles`)
	if err != nil {
		return fmt.Errorf("failed to get user profiles for leave seeding: %w", err)
	}

	if len(profiles) == 0 {
		logger.Warn("SeedLeaveData: No user profiles found to seed leaves for")
		return nil
	}

	// 2. Clear old leave mock data
	_ = db.ExecContext(ctx, `DELETE FROM leave_requests`)
	_ = db.ExecContext(ctx, `DELETE FROM leave_balances`)

	logger.Info(fmt.Sprintf("SeedLeaveData: Seeding leave balances and requests for %d employees for year 2026...", len(profiles)))

	currentYear := 2026
	lastAccrualMonth := 7 // July

	for _, p := range profiles {
		// Create a leave balance for each employee
		balanceID := uuid.New()
		queryBal := `
			INSERT INTO leave_balances (id, employee_id, year, accrued_days, used_days, balance, last_accrual_month, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		`
		// Grant 7 days accrued (Jan to Jul), 1 day used, remaining 6 days balance
		_, errBal := db.SQLx().ExecContext(ctx, queryBal, balanceID, p.ID, currentYear, 7.0, 1.0, 6.0, lastAccrualMonth)
		if errBal != nil {
			logger.Error("SeedLeaveData: Failed to insert leave balance", "employee", p.Code, "error", errBal)
			continue
		}

		// Create an approved past leave request (1 day in June)
		reqID1 := uuid.New()
		queryReq1 := `
			INSERT INTO leave_requests (id, employee_id, start_date, end_date, days_requested, reason, approved_by, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, NOW(), NOW())
		`
		_, errReq1 := db.SQLx().ExecContext(ctx, queryReq1, reqID1, p.ID, "2026-06-15", "2026-06-15", 1.0, "Nghỉ giải quyết việc gia đình", "APPROVED")
		if errReq1 != nil {
			logger.Error("SeedLeaveData: Failed to insert approved leave request", "employee", p.Code, "error", errReq1)
		}

		// Create a pending future leave request (1 day in August)
		reqID2 := uuid.New()
		queryReq2 := `
			INSERT INTO leave_requests (id, employee_id, start_date, end_date, days_requested, reason, approved_by, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, NOW(), NOW())
		`
		_, errReq2 := db.SQLx().ExecContext(ctx, queryReq2, reqID2, p.ID, "2026-08-10", "2026-08-10", 1.0, "Nghỉ phép đi du lịch hè", "PENDING")
		if errReq2 != nil {
			logger.Error("SeedLeaveData: Failed to insert pending leave request", "employee", p.Code, "error", errReq2)
		}
	}

	logger.Info("SeedLeaveData completed successfully!")
	return nil
}
