package repository

import (
	"cal-salary/core/logger"
	"cal-salary/modules/payroll/entity"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// ==========================================
// Job Position Standards Repository
// ==========================================

func (r *PayrollRepository) AssignStandardToPosition(ctx context.Context, jps *entity.JobPositionStandard) error {
	query := `
		INSERT INTO job_position_standards (job_description_id, job_standard_id, created_at)
		VALUES (:job_description_id, :job_standard_id, NOW())
		ON CONFLICT (job_description_id, job_standard_id) DO NOTHING
	`
	_, err := r.DB.NamedExecContext(ctx, query, jps)
	if err != nil {
		logger.Error("PayrollRepository:AssignStandardToPosition:Error %v", err)
		return err
	}
	return nil
}

func (r *PayrollRepository) GetStandardsByPosition(ctx context.Context, positionID uuid.UUID) ([]entity.JobPositionStandardDetail, error) {
	query := `
		SELECT jps.job_description_id, jps.job_standard_id, js.standard_code, js.name AS standard_name, js.allowance_value
		FROM job_position_standards jps
		JOIN job_standards js ON jps.job_standard_id = js.id
		WHERE jps.job_description_id = $1
		ORDER BY js.standard_code ASC
	`
	var list []entity.JobPositionStandardDetail
	err := r.DB.SelectContext(ctx, &list, query, positionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.JobPositionStandardDetail{}, nil
		}
		logger.Error("PayrollRepository:GetStandardsByPosition:Error %v", err)
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) RemoveStandardFromPosition(ctx context.Context, positionID uuid.UUID, standardID uuid.UUID) error {
	query := `DELETE FROM job_position_standards WHERE job_description_id = $1 AND job_standard_id = $2`
	_, err := r.DB.SQLx().ExecContext(ctx, query, positionID, standardID)
	if err != nil {
		logger.Error("PayrollRepository:RemoveStandardFromPosition:Error %v", err)
		return err
	}
	return nil
}

// ==========================================
// Employee Competencies Repository
// ==========================================

func (r *PayrollRepository) AssignEmployeeCompetencies(ctx context.Context, userProfileID uuid.UUID, competencyIDs []uuid.UUID) error {
	for _, cid := range competencyIDs {
		query := `
			INSERT INTO employee_competencies (user_profile_id, competency_id, created_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (user_profile_id, competency_id) DO NOTHING
		`
		_, err := r.DB.SQLx().ExecContext(ctx, query, userProfileID, cid)
		if err != nil {
			logger.Error("PayrollRepository:AssignEmployeeCompetencies:Error %v", err)
			return err
		}
	}
	return nil
}

func (r *PayrollRepository) GetCompetenciesByEmployee(ctx context.Context, userProfileID uuid.UUID) ([]entity.EmployeeCompetencyDetail, error) {
	query := `
		SELECT ec.user_profile_id, ec.competency_id, cd.name AS competency_name, cd.code AS competency_code, cd.point_value
		FROM employee_competencies ec
		JOIN competency_dictionaries cd ON ec.competency_id = cd.id
		WHERE ec.user_profile_id = $1
		ORDER BY cd.code ASC
	`
	var list []entity.EmployeeCompetencyDetail
	err := r.DB.SelectContext(ctx, &list, query, userProfileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.EmployeeCompetencyDetail{}, nil
		}
		logger.Error("PayrollRepository:GetCompetenciesByEmployee:Error %v", err)
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) RemoveEmployeeCompetency(ctx context.Context, userProfileID uuid.UUID, competencyID uuid.UUID) error {
	query := `DELETE FROM employee_competencies WHERE user_profile_id = $1 AND competency_id = $2`
	_, err := r.DB.SQLx().ExecContext(ctx, query, userProfileID, competencyID)
	if err != nil {
		logger.Error("PayrollRepository:RemoveEmployeeCompetency:Error %v", err)
		return err
	}
	return nil
}

// ==========================================
// System Settings Repository
// ==========================================

func (r *PayrollRepository) GetSystemSetting(ctx context.Context, key string) (*entity.SystemSetting, error) {
	query := `SELECT key, value, description, updated_at FROM system_settings WHERE key = $1`
	var setting entity.SystemSetting
	err := r.DB.GetContext(ctx, &setting, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetSystemSetting:Error %v", err)
		return nil, err
	}
	return &setting, nil
}

func (r *PayrollRepository) SetSystemSetting(ctx context.Context, setting *entity.SystemSetting) error {
	query := `
		INSERT INTO system_settings (key, value, description, updated_at)
		VALUES (:key, :value, :description, NOW())
		ON CONFLICT (key) 
		DO UPDATE SET value = EXCLUDED.value, description = EXCLUDED.description, updated_at = NOW()
	`
	_, err := r.DB.NamedExecContext(ctx, query, setting)
	if err != nil {
		logger.Error("PayrollRepository:SetSystemSetting:Error %v", err)
		return err
	}
	return nil
}

func (r *PayrollRepository) GetAllSystemSettings(ctx context.Context) ([]entity.SystemSetting, error) {
	query := `SELECT key, value, description, updated_at FROM system_settings ORDER BY key ASC`
	var settings []entity.SystemSetting
	err := r.DB.SelectContext(ctx, &settings, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.SystemSetting{}, nil
		}
		logger.Error("PayrollRepository:GetAllSystemSettings:Error %v", err)
		return nil, err
	}
	return settings, nil
}
