package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"cal-salary/core/utils"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const defaultEmployeePassword = "Employee@123"

// SeedEmployeeUsers creates a login account for each seeded user_profile that
// doesn't have one yet, and links user_profiles.user_id back to it. Without
// this link, non-admin features that resolve "the current user's employee
// profile" (OT requests, explanation requests, ...) can't find a profile for
// any seeded employee and fail with "Không tìm thấy hồ sơ nhân sự".
func SeedEmployeeUsers(ctx context.Context, db database.Database) error {
	var employeeRoleID uuid.UUID
	if err := db.GetContext(ctx, &employeeRoleID, `SELECT id FROM roles WHERE slug = 'employee' LIMIT 1`); err != nil {
		logger.Warn("SeedEmployeeUsers: employee role not found, skipping", "error", err)
		return nil
	}

	type profileRow struct {
		ID   uuid.UUID `db:"id"`
		Code string    `db:"code"`
	}
	var profiles []profileRow
	if err := db.SelectContext(ctx, &profiles, `SELECT id, code FROM user_profiles WHERE user_id IS NULL`); err != nil {
		return fmt.Errorf("failed to load user profiles without linked user: %w", err)
	}

	hashedPassword, err := utils.HashPassword(defaultEmployeePassword)
	if err != nil {
		return fmt.Errorf("failed to hash default employee password: %w", err)
	}

	successCount := 0
	errorCount := 0

	for _, p := range profiles {
		username := p.Code

		var userID uuid.UUID
		var existingUserID uuid.UUID
		err := db.GetContext(ctx, &existingUserID, `SELECT id FROM users WHERE username = $1 LIMIT 1`, username)
		if err == nil && existingUserID != uuid.Nil {
			userID = existingUserID
		} else if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedEmployeeUsers: error checking existing user", "username", username, "error", err)
			errorCount++
			continue
		} else {
			userID = uuid.New()
			insertErr := db.ExecContext(ctx,
				`INSERT INTO users (id, username, password, is_active, created_at, updated_at)
				 VALUES ($1, $2, $3, true, NOW(), NOW())`,
				userID, username, hashedPassword,
			)
			if insertErr != nil {
				logger.Error("SeedEmployeeUsers: failed to create user", "username", username, "error", insertErr)
				errorCount++
				continue
			}
		}

		if updErr := db.ExecContext(ctx, `UPDATE user_profiles SET user_id = $1, updated_at = NOW() WHERE id = $2`, userID, p.ID); updErr != nil {
			logger.Error("SeedEmployeeUsers: failed to link user_profile", "code", p.Code, "error", updErr)
			errorCount++
			continue
		}

		var existingRoleID uuid.UUID
		roleCheckErr := db.GetContext(ctx, &existingRoleID, `SELECT id FROM user_roles WHERE user_id = $1 AND role_id = $2 LIMIT 1`, userID, employeeRoleID)
		if roleCheckErr != nil {
			if roleInsertErr := db.ExecContext(ctx,
				`INSERT INTO user_roles (id, user_id, role_id, is_active, created_at, updated_at)
				 VALUES ($1, $2, $3, true, NOW(), NOW())`,
				uuid.New(), userID, employeeRoleID,
			); roleInsertErr != nil {
				logger.Error("SeedEmployeeUsers: failed to assign employee role", "username", username, "error", roleInsertErr)
				errorCount++
				continue
			}
		}

		successCount++
	}

	if errorCount > 0 {
		return fmt.Errorf("seed employee users completed with %d errors", errorCount)
	}

	logger.Info("Seeded employee user accounts successfully!", "linked", successCount, "default_password", defaultEmployeePassword)
	return nil
}
