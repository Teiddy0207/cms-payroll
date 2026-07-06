package repository

import (
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (r *PayrollRepository) CreateUserProfile(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	query := `
		INSERT INTO user_profiles (id, user_id, code, full_name, phone, avatar, date_of_birth, gender, position_id, department_id, created_at, updated_at)
		VALUES (:id, :user_id, :code, :full_name, :phone, :avatar, :date_of_birth, :gender, :position_id, :department_id, NOW(), NOW())
		RETURNING id, user_id, code, full_name, phone, avatar, date_of_birth, gender, position_id, department_id, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, profile)
	if err != nil {
		logger.Error("PayrollRepository:CreateUserProfile:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdProfile entity.UserProfile
		if err := rows.StructScan(&createdProfile); err != nil {
			logger.Error("PayrollRepository:CreateUserProfile:Scan:Error %v", err)
			return nil, err
		}
		return &createdProfile, nil
	}
	return nil, errors.New("failed to create user profile")
}

func (r *PayrollRepository) GetUserProfiles(ctx context.Context, qp params.QueryParams) ([]entity.UserProfile, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `
		FROM user_profiles u
		LEFT JOIN job_descriptions jd ON u.position_id = jd.id
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR u.code ILIKE $%d OR u.phone ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+qp.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalItems int
	err := r.DB.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		logger.Error("PayrollRepository:GetUserProfiles - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT
			u.id, u.user_id, u.code, u.full_name, u.phone, u.avatar,
			u.date_of_birth, u.gender, u.position_id, u.department_id,
			u.created_at, u.updated_at,
			jd.id AS jp_id, jd.code AS jp_code, jd.name AS jp_name,
			jd.description AS jp_description, jd.department_id AS jp_department_id
	` + baseQuery + whereClause + ` ORDER BY u.full_name ASC, u.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	type profileRow struct {
		entity.UserProfile
		JpID           *uuid.UUID `db:"jp_id"`
		JpCode         *string    `db:"jp_code"`
		JpName         *string    `db:"jp_name"`
		JpDescription  *string    `db:"jp_description"`
		JpDepartmentID *uuid.UUID `db:"jp_department_id"`
	}

	var rows []profileRow
	err = r.DB.SelectContext(ctx, &rows, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.UserProfile{}, 0, nil
		}
		logger.Error("PayrollRepository:GetUserProfiles - Select", err)
		return nil, 0, err
	}

	profiles := make([]entity.UserProfile, len(rows))
	for i, row := range rows {
		profiles[i] = row.UserProfile
		if row.JpID != nil {
			profiles[i].JobPosition = &entity.JobPosition{}
			profiles[i].JobPosition.ID = *row.JpID
			if row.JpCode != nil {
				profiles[i].JobPosition.Code = *row.JpCode
			}
			if row.JpName != nil {
				profiles[i].JobPosition.Name = *row.JpName
			}
			profiles[i].JobPosition.Description = row.JpDescription
			profiles[i].JobPosition.DepartmentID = row.JpDepartmentID
		}
	}

	return profiles, totalItems, nil
}


func (r *PayrollRepository) GetUserProfileByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error) {
	query := `
		SELECT
			u.id, u.user_id, u.code, u.full_name, u.phone, u.avatar,
			u.date_of_birth, u.gender, u.position_id, u.department_id,
			u.created_at, u.updated_at,
			jd.id AS jp_id, jd.code AS jp_code, jd.name AS jp_name,
			jd.description AS jp_description, jd.department_id AS jp_department_id
		FROM user_profiles u
		LEFT JOIN job_descriptions jd ON u.position_id = jd.id
		WHERE u.id = $1
	`

	type profileRow struct {
		entity.UserProfile
		JpID           *uuid.UUID `db:"jp_id"`
		JpCode         *string    `db:"jp_code"`
		JpName         *string    `db:"jp_name"`
		JpDescription  *string    `db:"jp_description"`
		JpDepartmentID *uuid.UUID `db:"jp_department_id"`
	}

	var row profileRow
	err := r.DB.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetUserProfileByID - Select", err)
		return nil, err
	}

	profile := row.UserProfile
	if row.JpID != nil {
		profile.JobPosition = &entity.JobPosition{}
		profile.JobPosition.ID = *row.JpID
		if row.JpCode != nil {
			profile.JobPosition.Code = *row.JpCode
		}
		if row.JpName != nil {
			profile.JobPosition.Name = *row.JpName
		}
		profile.JobPosition.Description = row.JpDescription
		profile.JobPosition.DepartmentID = row.JpDepartmentID
	}

	return &profile, nil
}

func (r *PayrollRepository) UpdateUserProfile(ctx context.Context, id uuid.UUID, profile *entity.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET full_name = :full_name, phone = :phone, avatar = :avatar, date_of_birth = :date_of_birth, gender = :gender, position_id = :position_id, department_id = :department_id, updated_at = NOW()
		WHERE id = :id
	`
	profile.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, profile)
	if err != nil {
		logger.Error("PayrollRepository:UpdateUserProfile:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user profile not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteUserProfile(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_profiles WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteUserProfile:Error %v", err)
		return err
	}
	return nil
}
