package repository

import (
	"cal-salary/core/database"
	"cal-salary/core/params"
	"cal-salary/modules/timekeeping/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TimekeepingRepositoryImpl struct {
	DB database.Database
}

func NewTimekeepingRepository(db database.Database) TimekeepingRepository {
	return &TimekeepingRepositoryImpl{DB: db}
}

func (r *TimekeepingRepositoryImpl) CreateAttendanceLog(ctx context.Context, log *entity.AttendanceLog) error {
	query := `
		INSERT INTO attendance_logs (id, employee_code, timestamp, location_gps, device_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, log.ID, log.EmployeeCode, log.Timestamp, log.LocationGPS, log.DeviceId)
	return err
}

func (r *TimekeepingRepositoryImpl) GetAttendanceLogs(ctx context.Context, employeeCode string, start, end time.Time) ([]entity.AttendanceLog, error) {
	var list []entity.AttendanceLog
	query := `
		SELECT id, employee_code, timestamp, location_gps, device_id, created_at
		FROM attendance_logs
		WHERE employee_code = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp ASC
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, employeeCode, start, end)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.AttendanceLog{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *TimekeepingRepositoryImpl) GetAttendanceLogsForCalculation(ctx context.Context, start, end time.Time) ([]entity.AttendanceLog, error) {
	var list []entity.AttendanceLog
	query := `
		SELECT id, employee_code, timestamp, location_gps, device_id, created_at
		FROM attendance_logs
		WHERE timestamp >= $1 AND timestamp <= $2
		ORDER BY employee_code, timestamp ASC
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, start, end)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.AttendanceLog{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *TimekeepingRepositoryImpl) UpsertDailyAttendanceSheet(ctx context.Context, sheet *entity.DailyAttendanceSheet) error {
	query := `
		INSERT INTO daily_attendance_sheets (id, employee_id, date, check_in, check_out, actual_work_day, ot_hours, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (employee_id, date) DO UPDATE SET
			check_in = EXCLUDED.check_in,
			check_out = EXCLUDED.check_out,
			actual_work_day = EXCLUDED.actual_work_day,
			ot_hours = EXCLUDED.ot_hours,
			status = EXCLUDED.status,
			updated_at = NOW()
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, sheet.ID, sheet.EmployeeID, sheet.Date, sheet.CheckIn, sheet.CheckOut, sheet.ActualWorkDay, sheet.OTHours, sheet.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetDailyAttendanceSheets(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID, start, end time.Time, qp params.QueryParams) ([]entity.DailyAttendanceSheet, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	if offset < 0 {
		offset = 0
	}
	pageSize := qp.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	baseQuery := `
		FROM daily_attendance_sheets sh
		JOIN user_profiles u ON sh.employee_id = u.id
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("sh.date >= $%d AND sh.date <= $%d", argIndex, argIndex+1))
	args = append(args, start, end)
	argIndex += 2

	if employeeID != nil {
		conditions = append(conditions, fmt.Sprintf("sh.employee_id = $%d", argIndex))
		args = append(args, *employeeID)
		argIndex++
	}

	if departmentID != nil {
		conditions = append(conditions, fmt.Sprintf("u.department_id = $%d", argIndex))
		args = append(args, *departmentID)
		argIndex++
	}

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR u.code ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+qp.Search+"%")
		argIndex++
	}

	if status, ok := qp.Filters["status"]; ok && status != "" && status != "ALL" {
		conditions = append(conditions, fmt.Sprintf("sh.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if dateStr, ok := qp.Filters["date"]; ok && dateStr != "" {
		conditions = append(conditions, fmt.Sprintf("sh.date = $%d", argIndex))
		args = append(args, dateStr)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalItems int
	err := r.DB.SQLx().GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT sh.id, sh.employee_id, sh.date, sh.check_in, sh.check_out, sh.actual_work_day, sh.ot_hours, sh.status, sh.created_at, sh.updated_at
	` + baseQuery + whereClause + ` ORDER BY sh.date ASC, u.full_name ASC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var list []entity.DailyAttendanceSheet
	err = r.DB.SQLx().SelectContext(ctx, &list, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.DailyAttendanceSheet{}, 0, nil
		}
		return nil, 0, err
	}

	return list, totalItems, nil
}

func (r *TimekeepingRepositoryImpl) CreateExplanationRequest(ctx context.Context, req *entity.ExplanationRequest) error {
	query := `
		INSERT INTO explanation_requests (id, employee_id, date, reason, approved_by, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, req.ID, req.EmployeeID, req.Date, req.Reason, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetExplanationRequestByID(ctx context.Context, id uuid.UUID) (*entity.ExplanationRequest, error) {
	var req entity.ExplanationRequest
	query := `
		SELECT id, employee_id, date, reason, approved_by, status, created_at, updated_at
		FROM explanation_requests
		WHERE id = $1
	`
	err := r.DB.SQLx().GetContext(ctx, &req, query, id)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *TimekeepingRepositoryImpl) UpdateExplanationRequest(ctx context.Context, req *entity.ExplanationRequest) error {
	query := `
		UPDATE explanation_requests
		SET approved_by = $2, status = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, req.ID, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetExplanationRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID) ([]entity.ExplanationRequest, error) {
	var list []entity.ExplanationRequest
	var err error

	if employeeID != nil {
		query := `
			SELECT id, employee_id, date, reason, approved_by, status, created_at, updated_at
			FROM explanation_requests
			WHERE employee_id = $1
			ORDER BY created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, *employeeID)
	} else if departmentID != nil {
		query := `
			SELECT r.id, r.employee_id, r.date, r.reason, r.approved_by, r.status, r.created_at, r.updated_at
			FROM explanation_requests r
			JOIN user_profiles u ON r.employee_id = u.id
			WHERE u.department_id = $1
			ORDER BY r.created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, *departmentID)
	} else {
		query := `
			SELECT id, employee_id, date, reason, approved_by, status, created_at, updated_at
			FROM explanation_requests
			ORDER BY created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.ExplanationRequest{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *TimekeepingRepositoryImpl) CreateOTRequest(ctx context.Context, req *entity.OTRequest) error {
	query := `
		INSERT INTO ot_requests (id, employee_id, date, hours_requested, is_night_ot, is_holiday_ot, approved_by, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, req.ID, req.EmployeeID, req.Date, req.HoursRequested, req.IsNightOT, req.IsHolidayOT, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetOTRequestByID(ctx context.Context, id uuid.UUID) (*entity.OTRequest, error) {
	var req entity.OTRequest
	query := `
		SELECT id, employee_id, date, hours_requested, is_night_ot, is_holiday_ot, approved_by, status, created_at, updated_at
		FROM ot_requests
		WHERE id = $1
	`
	err := r.DB.SQLx().GetContext(ctx, &req, query, id)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *TimekeepingRepositoryImpl) UpdateOTRequest(ctx context.Context, req *entity.OTRequest) error {
	query := `
		UPDATE ot_requests
		SET approved_by = $2, status = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, req.ID, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetOTRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID) ([]entity.OTRequest, error) {
	var list []entity.OTRequest
	var err error

	if employeeID != nil {
		query := `
			SELECT id, employee_id, date, hours_requested, is_night_ot, is_holiday_ot, approved_by, status, created_at, updated_at
			FROM ot_requests
			WHERE employee_id = $1
			ORDER BY created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, *employeeID)
	} else if departmentID != nil {
		query := `
			SELECT r.id, r.employee_id, r.date, r.hours_requested, r.is_night_ot, r.is_holiday_ot, r.approved_by, r.status, r.created_at, r.updated_at
			FROM ot_requests r
			JOIN user_profiles u ON r.employee_id = u.id
			WHERE u.department_id = $1
			ORDER BY r.created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, *departmentID)
	} else {
		query := `
			SELECT id, employee_id, date, hours_requested, is_night_ot, is_holiday_ot, approved_by, status, created_at, updated_at
			FROM ot_requests
			ORDER BY created_at DESC
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.OTRequest{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *TimekeepingRepositoryImpl) CreateFaceTemplate(ctx context.Context, template *entity.EmployeeFaceTemplate) error {
	query := `
		INSERT INTO employee_face_templates (id, employee_code, face_data, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (employee_code) DO UPDATE
		SET face_data = EXCLUDED.face_data
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, template.ID, template.EmployeeCode, template.FaceData)
	return err
}

func (r *TimekeepingRepositoryImpl) GetFaceTemplates(ctx context.Context) ([]entity.EmployeeFaceTemplate, error) {
	var list []entity.EmployeeFaceTemplate
	query := `
		SELECT id, employee_code, face_data, created_at
		FROM employee_face_templates
		ORDER BY created_at DESC
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.EmployeeFaceTemplate{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *TimekeepingRepositoryImpl) GetFaceTemplateByCode(ctx context.Context, code string) (*entity.EmployeeFaceTemplate, error) {
	var template entity.EmployeeFaceTemplate
	query := `
		SELECT id, employee_code, face_data, created_at
		FROM employee_face_templates
		WHERE employee_code = $1
	`
	err := r.DB.SQLx().GetContext(ctx, &template, query, code)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TimekeepingRepositoryImpl) DeleteFaceTemplate(ctx context.Context, code string) error {
	query := `
		DELETE FROM employee_face_templates
		WHERE employee_code = $1
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, code)
	return err
}
