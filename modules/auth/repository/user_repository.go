package repository

import (
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/auth/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *AuthRepository) PrivateCreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (email, username, password, position_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, password, email_verified_at, is_active, position_id, created_at, updated_at
	`

	var createdUser entity.User
	err := r.DB.GetContext(ctx, &createdUser, query, user.Email, user.Username, user.Password, user.PositionID)
	if err != nil {
		logger.Error("AuthRepository:PrivateCreateUser - Exec", err)
		return nil, err
	}

	return &createdUser, nil
}

func (r *AuthRepository) PrivateUpdatePasswordUser(ctx context.Context, userID uuid.UUID, password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	err := r.DB.ExecContext(ctx, query, password, userID)
	if err != nil {
		logger.Error("AuthRepository:PrivateUpdatePasswordUser - Exec", err)
		return err
	}
	return nil
}

func (r *AuthRepository) PrivateGetUsers(ctx context.Context, params params.QueryParams) (*entity.PaginatedUserWithProfileEntity, error) {
	// Tính offset cho pagination
	offset := (params.PageNumber - 1) * params.PageSize

	baseQuery := `FROM users u
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN job_descriptions jd ON up.position_id = jd.id`

	// Thêm điều kiện search nếu có
	var whereClause string
	var args []interface{}

	conditions := []string{}
	argIndex := 1

	if params.Search != "" {
		search := "%" + params.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`
		(
			u.email ILIKE $%d OR 
			u.username ILIKE $%d
		)
	`, argIndex, argIndex))
		args = append(args, search)
		argIndex++
	}

	// if officeID, exists := params.Filters["office_id"]; exists && officeID != "" {
	// 	conditions = append(conditions, fmt.Sprintf("up.official_id = $%d", argIndex))
	// 	args = append(args, officeID)
	// 	argIndex++
	// }

	// Filter theo role_ids nếu có
	if roleIdsStr, exists := params.Filters["role_ids"]; exists && roleIdsStr != "" {
		// Parse comma-separated role IDs
		roleIds := strings.Split(roleIdsStr, ",")
		trimmedRoleIds := make([]string, 0, len(roleIds))
		for _, id := range roleIds {
			trimmed := strings.TrimSpace(id)
			if trimmed != "" {
				// Validate UUID format
				if _, err := uuid.Parse(trimmed); err == nil {
					trimmedRoleIds = append(trimmedRoleIds, trimmed)
				}
			}
		}

		if len(trimmedRoleIds) > 0 {
			// Sử dụng subquery để tránh duplicate khi JOIN
			placeholders := make([]string, len(trimmedRoleIds))
			for i := range trimmedRoleIds {
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				args = append(args, trimmedRoleIds[i])
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf(`u.id IN (
				SELECT DISTINCT ur.user_id
				FROM user_roles ur
				WHERE ur.role_id IN (%s)
				  AND ur.is_active = true
				  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
			)`, strings.Join(placeholders, ",")))
		}
	}

	// Filter users chưa liên kết với user_profiles
	if ExcludeProfile, ok := params.Filters["exclude_profile"]; ok && ExcludeProfile == "true" {
		conditions = append(conditions, "up.id IS NULL")
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Query để đếm tổng số records (dùng DISTINCT để tránh đếm duplicate khi có nhiều profiles)
	countQuery := "SELECT COUNT(DISTINCT u.id) " + baseQuery + whereClause

	var totalItems int
	err := r.DB.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		logger.Error("AuthRepository:PrivateGetUsers - Count", err)
		return nil, err
	}

	// Query để lấy data với pagination, dùng json_agg để lấy roles
	dataQuery := `
		SELECT 
			u.id, 
			u.email, 
			u.username, 
			u.password,
			u.email_verified_at,
			u.is_active,
			u.position_id,
			u.created_at, 
			u.updated_at,
			up.id AS profile_id,
			up.code AS profile_code,
			up.full_name AS profile_full_name,
			up.phone AS profile_phone,
			up.gender AS profile_gender,
			up.avatar AS profile_avatar,
			up.position_id,
			jd.code AS position_code,
			jd.name AS position_name,
			(
				SELECT r2.id
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_id,
			(
				SELECT r2.slug
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_code,
			(
				SELECT r2.name
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_name,
			(
				SELECT r2.description
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_description
	` + baseQuery + whereClause + `
		GROUP BY u.id, u.email, u.username, u.password, u.email_verified_at, u.is_active, u.position_id, u.created_at, u.updated_at,
			up.id, up.full_name, up.phone, up.gender,
			jd.code, jd.name`

	var orderByClause string
	if params.OrderBy != "" {
		allowedColumns := map[string]string{
			"full_name":  "up.full_name",
			"email":      "u.email",
			"username":   "u.username",
			"created_at": "u.created_at",
			"updated_at": "u.updated_at",
		}

		if col, ok := allowedColumns[params.OrderBy]; ok {
			order := "ASC"
			if params.OrderDesc {
				order = "DESC"
			}
			orderByClause = fmt.Sprintf(" ORDER BY %s %s", col, order)
		} else {
			orderByClause = " ORDER BY u.created_at DESC"
		}
	} else {
		orderByClause = " ORDER BY u.created_at DESC"
	}

	dataQuery += orderByClause + `
		OFFSET $` + fmt.Sprintf("%d", argIndex+1) + ` ROWS FETCH NEXT $` + fmt.Sprintf("%d", argIndex) + ` ROWS ONLY`

	// Thêm params cho pagination
	args = append(args, params.PageSize, offset)

	var usersWithProfile []entity.UserWithProfile
	err = r.DB.SelectContext(ctx, &usersWithProfile, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("AuthRepository:PrivateGetUsers - Select", err)
		return nil, err
	}

	// Tạo response pagination với UserWithProfile
	response := &entity.PaginatedUserWithProfileEntity{
		Items:      usersWithProfile,
		TotalItems: totalItems,
		PageNumber: params.PageNumber,
		PageSize:   params.PageSize,
	}

	return response, nil
}

func (r *AuthRepository) PrivateGetUser(ctx context.Context, id uuid.UUID) (*entity.UserDetail, error) {
	query := `
		SELECT
			u.id                     AS id,
			u.email,
			u.username,
			u.is_active,
			u.created_at,
			COALESCE(up.phone, '') AS phone,
			up.full_name,
			up.avatar,
			up.date_of_birth,
			up.gender,
			(
				SELECT r2.id
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_id,
			(
				SELECT r2.slug
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_code,
			(
				SELECT r2.name
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_name,
			(
				SELECT r2.description
				FROM user_roles ur2
				INNER JOIN roles r2 ON ur2.role_id = r2.id
				WHERE ur2.user_id = u.id
				  AND ur2.is_active = true
				  AND r2.is_active = true
				  AND (ur2.expires_at IS NULL OR ur2.expires_at > NOW())
				ORDER BY r2.name
				LIMIT 1
			) AS role_description,
			(
				SELECT COALESCE(json_agg(
					json_build_object(
						'id', p.id,
						'name', p.name,
						'slug', p.slug,
						'resource', p.resource,
						'action', p.action,
						'description', p.description,
						'is_system', p.is_system,
						'created_at', p.created_at,
						'updated_at', p.updated_at
					)
				), '[]'::json)
				FROM permissions p
				WHERE p.id IN (
					-- Permissions từ roles của user
					SELECT DISTINCT rp.permission_id
					FROM role_permissions rp
					INNER JOIN user_roles ur ON rp.role_id = ur.role_id
					INNER JOIN roles r ON ur.role_id = r.id
					WHERE ur.user_id = u.id 
					  AND ur.is_active = true 
					  AND r.is_active = true
					  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
					
					UNION
					
					-- Permissions trực tiếp được gán cho user
					SELECT up.permission_id
					FROM user_permissions up
					WHERE up.user_id = u.id 
					  AND up.granted = true 
					  AND (up.expires_at IS NULL OR up.expires_at > NOW())
				)
			) AS permissions
		FROM users u
		LEFT JOIN user_profiles up
			ON u.id = up.user_id
		WHERE u.id = $1;
	`
	var userDetail entity.UserDetail
	err := r.DB.GetContext(ctx, &userDetail, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("AuthRepository:PrivateGetUser:UserNotFound:", id)
			return nil, nil // Return nil instead of error for not found case
		}
		logger.Error("AuthRepository:PrivateGetUser:Error:", err)
		return nil, err
	}

	return &userDetail, nil
}

func (r *AuthRepository) PrivateUpdateUser(ctx context.Context, user *entity.User, userId uuid.UUID) error {
	query := `
		UPDATE users 
		SET 
			email = :email,
			username = :username,
			password = :password,
			email_verified_at = :email_verified_at,
			position_id = :position_id,
			is_active = :is_active,
			updated_at = NOW()
		WHERE id = :id
	`

	user.ID = userId

	result, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		logger.Error("AuthRepository:PrivateUpdateUser:Error:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("AuthRepository:PrivateUpdateUser:RowsAffected:Error:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET 
			email = :email,
			phone = :phone,
			username = :username,
			password = :password,
			email_verified_at = :email_verified_at,
			phone_verified_at = :phone_verified_at,
			locked_until = :locked_until,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`
	result, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		logger.Error("AuthRepository:UpdateUser:Error:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("AuthRepository:UpdateUser:RowsAffected:Error:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *AuthRepository) GetUserByIdentifier(ctx context.Context, identifier string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE email = $1 OR id::text = $1 OR username = $1`
	err := r.DB.GetContext(ctx, &user, query, identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found is not an error, return nil
			return nil, nil
		}
		logger.Error("AuthRepository:GetUserByIdentifier:Error:", err)
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (email, username, password)
		VALUES (:email, :username, :password)
		RETURNING id, email, username, password, email_verified_at, position_id, is_active, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, user)
	if err != nil {
		logger.Error("AuthRepository:CreateUser:Error:", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var inserted entity.User
		if err := rows.StructScan(&inserted); err != nil {
			logger.Error("AuthRepository:CreateUser:StructScan", err)
			return nil, err
		}
		return &inserted, nil
	}
	return nil, fmt.Errorf("insert user failed")
}

func (r *AuthRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error) {
	// Query để lấy tất cả permissions của user từ 2 nguồn:
	// 1. Permissions trực tiếp được gán cho user (user_permissions)
	// 2. Permissions từ các roles mà user được gán (user_roles -> role_permissions)
	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.resource, p.action, p.description, p.is_system, p.created_at, p.updated_at
		FROM permissions p
		WHERE p.id IN (
			-- Permissions từ roles của user
			SELECT rp.permission_id
			FROM role_permissions rp
			INNER JOIN user_roles ur ON rp.role_id = ur.role_id
			INNER JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 
			  AND ur.is_active = true 
			  AND r.is_active = true
			  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
			
			UNION
			
			-- Permissions trực tiếp được gán cho user
			SELECT up.permission_id
			FROM user_permissions up
			WHERE up.user_id = $1 
			  AND up.granted = true 
			  AND (up.expires_at IS NULL OR up.expires_at > NOW())
		)
		ORDER BY p.resource, p.action
	`

	var permissions []entity.Permission
	err := r.DB.SelectContext(ctx, &permissions, query, userID)
	if err != nil {
		logger.Error("AuthRepository:GetUserPermissions:Error:", err)
		return nil, err
	}

	return permissions, nil
}

func (r *AuthRepository) PrivateDeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = :id`

	_, err := r.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		logger.Error("AuthRepository:PrivateDeleteUser - Delete", err)
		return err
	}
	return nil
}

func (r *AuthRepository) PrivateGetUserbyIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.UserWithProfile, error) {
	if len(ids) == 0 {
		return []*entity.UserWithProfile{}, nil
	}

	query := `
        SELECT 
			u.id, 
			u.email, 
			u.username, 
			u.password,
			u.email_verified_at,
			u.is_active,
			u.position_id,
			u.created_at, 
			u.updated_at,

			-- Profile
			up.id AS profile_id,
			up.code AS profile_code,
			up.full_name AS profile_full_name,
			up.phone AS profile_phone,
			up.gender AS profile_gender,
			up.avatar AS profile_avatar,
			up.position_id AS position_id,

			-- Job description
			jd.code AS position_code,
			jd.name AS position_name

		FROM users u
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN job_descriptions jd ON up.position_id = jd.id

		WHERE u.id = ANY($1)
		GROUP BY 
			u.id, 
			u.email, 
			u.username, 
			u.password,
			u.email_verified_at,
			u.is_active,
			u.position_id,
			u.created_at, 
			u.updated_at,
			up.id,
			up.code,
			up.full_name,
			up.phone,
			up.gender,
			up.avatar,
			up.position_id,
			jd.code,
			jd.name
		ORDER BY u.created_at DESC;
    `
	var users []*entity.UserWithProfile

	err := r.DB.SelectContext(ctx, &users, query, pq.Array(ids))
	if err != nil {
		logger.Error("PrivateGetUserbyIDs query error:", err)
		return nil, err
	}
	return users, nil
}
