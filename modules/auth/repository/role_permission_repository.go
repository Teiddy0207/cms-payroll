package repository

import (
	"context"
	"database/sql"

	"cal-salary/core/logger"
	"cal-salary/modules/auth/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (repo *AuthRepository) PrivateGetPermissionsByUserID(ctx context.Context, userID uuid.UUID) (*[]entity.Permission, error) {
	query := `
		SELECT DISTINCT p.slug, p.resource, p.action
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1
		  AND ur.is_active = true
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY p.resource, p.action
	`

	var permissions []entity.Permission
	err := repo.DB.SelectContext(ctx, &permissions, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("AuthRepository:PrivateGetPermissionsByUserID:NoRows:", err)
			return nil, nil
		}
		logger.Error("AuthRepository:PrivateGetPermissionsByUserID:Error:", err)
		return nil, err
	}

	return &permissions, nil
}

func (repo *AuthRepository) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) (*[]entity.Permission, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.slug,
			p.resource,
			p.action,
			p.description,
			p.is_system,
			p.created_at,
			p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	var permissions []entity.Permission
	err := repo.DB.SelectContext(ctx, &permissions, query, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &[]entity.Permission{}, nil
		}
		logger.Error("AuthRepository:GetPermissionsByRoleID:Error:", err)
		return nil, err
	}

	return &permissions, nil
}

func (repo *AuthRepository) PrivateAssignPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, grantedBy uuid.UUID) error {
	tx, err := repo.DB.SQLx().BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("AuthRepository:PrivateAssignPermissionToRole:BeginTx:Error:", err)
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	deleteQuery := `DELETE FROM role_permissions WHERE role_id = $1`
	if _, err = tx.ExecContext(ctx, deleteQuery, roleID); err != nil {
		logger.Error("AuthRepository:PrivateAssignPermissionToRole:Delete:Error:", err)
		return err
	}

	if len(permissionIDs) > 0 {
		uniquePermissionIDs := make([]uuid.UUID, 0, len(permissionIDs))
		seen := make(map[uuid.UUID]struct{}, len(permissionIDs))
		for _, permissionID := range permissionIDs {
			if _, exists := seen[permissionID]; exists {
				continue
			}
			seen[permissionID] = struct{}{}
			uniquePermissionIDs = append(uniquePermissionIDs, permissionID)
		}

		insertQuery := `
			INSERT INTO role_permissions (role_id, permission_id, granted_by, granted_at)
			SELECT $1, permission_id, $3, NOW()
			FROM UNNEST($2::uuid[]) AS permission_id
		`
		if _, err = tx.ExecContext(ctx, insertQuery, roleID, pq.Array(uniquePermissionIDs), grantedBy); err != nil {
			logger.Error("AuthRepository:PrivateAssignPermissionToRole:InsertBatch:Error:", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		logger.Error("AuthRepository:PrivateAssignPermissionToRole:Commit:Error:", err)
		return err
	}

	return nil
}
