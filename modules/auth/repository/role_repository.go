package repository

import (
	"cal-salary/core/constants"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/auth/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (r *AuthRepository) PrivateCreateRole(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	query := `
		INSERT INTO roles (name, slug, description, is_system, is_active)
		VALUES (:name, :slug, :description, :is_system, :is_active)
		RETURNING id, name, slug, description, is_system, is_active, created_at, updated_at
	`
	
	rows, err := r.DB.NamedQueryContext(ctx, query, role)
	if err != nil {
		logger.Error("AuthRepository:CreateRole:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdRole entity.Role
		if err := rows.StructScan(&createdRole); err != nil {
			logger.Error("AuthRepository:CreateRole:Scan:Error %v", err)
			return nil, err
		}
		return &createdRole, nil
	}

	return nil, errors.New("failed to create role")
}

func (r *AuthRepository) PrivateGetRoles(ctx context.Context, params params.QueryParams) (*entity.PaginatedRoleEntity, error) {
	// Tính offset cho pagination
	offset := (params.PageNumber - 1) * params.PageSize

	// Base query để lấy roles
	baseQuery := `FROM roles r`

	// Thêm điều kiện search nếu có
	var whereClause string
	var args []interface{}

	conditions := []string{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(r.name ILIKE $%d OR r.slug ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Query để đếm tổng số records 
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause

	var totalItems int
	err := r.DB.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		logger.Error("AuthRepository:PrivateGetRoles - Count", err)
		return nil, err
	}

	// Query để lấy data với pagination
	var orderClause string
	if params.OrderBy != "" {
		allowedColumns := map[string]string{
			"name":       "r.name",
			"slug":       "r.slug",
			"created_at": "r.created_at",
			"updated_at": "r.updated_at",
			"is_active":  "r.is_active",
			"is_system":  "r.is_system",
		}
		
		if col, ok := allowedColumns[params.OrderBy]; ok {
			order := "ASC"
			if params.OrderDesc {
				order = "DESC"
			}
			orderClause = fmt.Sprintf(" ORDER BY %s %s", col, order)
		} else {
			// Nếu order_by không hợp lệ, dùng default
			orderClause = " ORDER BY r.name ASC, r.created_at DESC"
		}
	} else {
		orderClause = " ORDER BY r.name ASC, r.created_at DESC"
	}

	dataQuery := `
		SELECT 
			r.id, 
			r.name, 
			r.slug, 
			r.description, 
			r.is_system,
			r.is_active,
			r.created_at, 
			r.updated_at
	` + baseQuery + whereClause + orderClause + `
		OFFSET $` + fmt.Sprintf("%d", argIndex) + ` ROWS FETCH NEXT $` + fmt.Sprintf("%d", argIndex+1) + ` ROWS ONLY`

	// Thêm params cho pagination
	args = append(args, offset, params.PageSize)

	var roles []entity.Role
	err = r.DB.SelectContext(ctx, &roles, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("AuthRepository:PrivateGetRoles - Select", err)
		return nil, err
	}

	// Tạo response pagination
	response := &entity.PaginatedRoleEntity{
		Items:      roles,
		TotalItems: totalItems,
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
	}

	return response, nil
}

func (r *AuthRepository) PrivateGetRoleByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	query := `SELECT * FROM roles WHERE id = $1`
	err := r.DB.GetContext(ctx, &role, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("AuthRepository:PrivateGetRoleByID - Select", err)
		return nil, err
	}
	return &role, nil
}

func (r *AuthRepository) PrivateUpdateRole(ctx context.Context, id uuid.UUID, role *entity.Role) error {
	query := `
		UPDATE roles
		SET name = :name, slug = :slug, description = :description, is_system = :is_system, is_active = :is_active
		WHERE id = :id
	`

	role.ID = id

	result, err := r.DB.NamedExecContext(ctx, query, role)
	if err != nil {
		logger.Error("AuthRepository:PrivateUpdateRole:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("AuthRepository:PrivateUpdateRole:Error %v", err)
		return err
	}

	if rowsAffected == constants.DefaultZeroValue {
		return errors.New("role not found")
	}

	return nil
}

func (r *AuthRepository) PrivateDeleteRole(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id = $1`
	err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("AuthRepository:PrivateDeleteRole:Error %v", err)
		return err
	}
	return nil
}

// GetRoleByNameOrSlug lấy role theo name hoặc slug
func (r *AuthRepository) GetRoleByNameOrSlug(ctx context.Context, name, slug string) (*entity.Role, error) {
	var role entity.Role
	query := `SELECT * FROM roles WHERE name = $1 OR slug = $2 LIMIT 1`
	err := r.DB.GetContext(ctx, &role, query, name, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("AuthRepository:GetRoleByNameOrSlug:Error %v", err)
		return nil, err
	}
	return &role, nil
}

// GetUserIDsByRoleID lấy danh sách user_id được gán role_id tương ứng
func (r *AuthRepository) GetUserIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM user_roles WHERE role_id = $1 AND is_active = true AND (expires_at IS NULL OR expires_at > NOW())`
	var userIDs []uuid.UUID
	err := r.DB.SelectContext(ctx, &userIDs, query, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []uuid.UUID{}, nil
		}
		logger.Error("AuthRepository:GetUserIDsByRoleID - Select", err)
		return nil, err
	}
	if userIDs == nil {
		return []uuid.UUID{}, nil
	}
	return userIDs, nil
}

// PrivateAssignRoleToUser gán role cho user
func (r *AuthRepository) PrivateAssignRoleToUser(ctx context.Context, req *entity.UserRole) error {
	query := `
		INSERT INTO user_roles (id, code, name, description, user_id, role_id, assigned_by, assigned_at, expires_at, is_active, created_at, updated_at)
		VALUES (:id, :code, :name, :description, :user_id, :role_id, :assigned_by, :assigned_at, :expires_at, :is_active, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			user_id = EXCLUDED.user_id,
			role_id = EXCLUDED.role_id,
			assigned_by = EXCLUDED.assigned_by,
			assigned_at = EXCLUDED.assigned_at,
			expires_at = EXCLUDED.expires_at,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
	`
	// Generate new UUID for ID if it is empty
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	// Gán base entity fields
	if req.BaseEntity.ID == uuid.Nil {
		req.BaseEntity.ID = uuid.MustParse(req.ID)
	}

	_, err := r.DB.NamedExecContext(ctx, query, req)
	if err != nil {
		logger.Error("AuthRepository:PrivateAssignRoleToUser:Error %v", err)
		return err
	}
	return nil
}