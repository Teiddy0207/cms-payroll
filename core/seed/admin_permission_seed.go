package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"

	"github.com/google/uuid"
)

// SeedAdminRolePermissions grants every existing permission to the "admin"
// role. Roles/permissions are only linked through the Phân quyền UI, which
// itself requires a permission to view — without this, a freshly seeded
// admin role has zero permissions and can never reach that screen to grant
// itself any.
func SeedAdminRolePermissions(ctx context.Context, db database.Database) error {
	var adminRoleID uuid.UUID
	if err := db.GetContext(ctx, &adminRoleID, `SELECT id FROM roles WHERE slug = 'admin' LIMIT 1`); err != nil {
		logger.Warn("SeedAdminRolePermissions: admin role not found, skipping", "error", err)
		return nil
	}

	err := db.ExecContext(ctx, `
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT $1, p.id FROM permissions p
		ON CONFLICT (role_id, permission_id) DO NOTHING`,
		adminRoleID,
	)
	if err != nil {
		logger.Error("SeedAdminRolePermissions: failed to grant permissions", "error", err)
		return err
	}

	logger.Info("SeedAdminRolePermissions: Completed granting all permissions to admin role")
	return nil
}
