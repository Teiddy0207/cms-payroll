package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type roleSeedItem struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	IsActive    bool
}

// SeedRolesAndAdminUserRole seeds roles and assigns admin role to the admin user
func SeedRolesAndAdminUserRole(ctx context.Context, db database.Database) error {
	adminRoleID := uuid.MustParse("10000000-0000-0000-0000-000000000001")
	managerRoleID := uuid.MustParse("10000000-0000-0000-0000-000000000002")
	employeeRoleID := uuid.MustParse("10000000-0000-0000-0000-000000000003")

	roles := []roleSeedItem{
		{
			ID:          adminRoleID,
			Name:        "Quản trị viên",
			Slug:        "admin",
			Description: "Quyền quản trị toàn hệ thống",
			IsActive:    true,
		},
		{
			ID:          managerRoleID,
			Name:        "Quản lý",
			Slug:        "manager",
			Description: "Quản lý phòng ban",
			IsActive:    true,
		},
		{
			ID:          employeeRoleID,
			Name:        "Nhân viên",
			Slug:        "employee",
			Description: "Nhân viên thông thường",
			IsActive:    true,
		},
	}

	// Seed roles
	insertedRoles := 0
	skippedRoles := 0
	for _, r := range roles {
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM roles WHERE slug = $1 LIMIT 1`, r.Slug)
		if err == nil && existingID != uuid.Nil {
			skippedRoles++
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedRoles: Error checking existing role", "slug", r.Slug, "error", err)
			continue
		}

		err = db.ExecContext(ctx,
			`INSERT INTO roles (id, name, slug, description, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			r.ID, r.Name, r.Slug, r.Description, r.IsActive,
		)
		if err != nil {
			logger.Error("SeedRoles: Failed to insert role", "slug", r.Slug, "error", err)
			continue
		}
		insertedRoles++
		logger.Info("SeedRoles: Created role", "slug", r.Slug)
	}

	logger.Info("SeedRoles: Completed roles",
		"inserted", insertedRoles,
		"skipped", skippedRoles,
		"total", len(roles),
	)

	// Tìm TẤT CẢ user có username = 'admin'
	var adminUserIDs []uuid.UUID
	err := db.SelectContext(ctx, &adminUserIDs, `SELECT id FROM users WHERE username = 'admin'`)
	if err != nil || len(adminUserIDs) == 0 {
		logger.Warn("SeedRoles: No admin users found, skipping user_roles seed")
		return nil
	}

	// Lấy ID thực của admin role
	var realAdminRoleID uuid.UUID
	err2 := db.GetContext(ctx, &realAdminRoleID, `SELECT id FROM roles WHERE slug = 'admin' LIMIT 1`)
	if err2 != nil {
		logger.Warn("SeedRoles: Admin role not found in DB, skipping user_roles seed")
		return nil
	}

	for _, adminUserID := range adminUserIDs {
		// Kiểm tra xem đã có role chưa
		var existingUserRoleID uuid.UUID
		check := db.GetContext(ctx, &existingUserRoleID,
			`SELECT id FROM user_roles WHERE user_id = $1 AND role_id = $2 LIMIT 1`,
			adminUserID, realAdminRoleID,
		)
		if check == nil && existingUserRoleID != uuid.Nil {
			logger.Info("SeedRoles: Admin user_role already exists, skipping", "user_id", adminUserID)
			continue
		}

		insertErr := db.ExecContext(ctx,
			`INSERT INTO user_roles (id, user_id, role_id, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, true, NOW(), NOW())`,
			uuid.New(), adminUserID, realAdminRoleID,
		)
		if insertErr != nil {
			logger.Error("SeedRoles: Failed to assign admin role", "user_id", adminUserID, "error", insertErr)
			continue
		}
		logger.Info("SeedRoles: Assigned admin role to user", "user_id", adminUserID, "role_id", realAdminRoleID)
	}
	return nil
}
