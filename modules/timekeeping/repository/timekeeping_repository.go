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
		INSERT INTO attendance_logs (id, event_id, employee_code, timestamp, location_gps, device_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (event_id) DO NOTHING
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, log.ID, log.EventID, log.EmployeeCode, log.Timestamp, log.LocationGPS, log.DeviceId)
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
		INSERT INTO employee_face_templates (id, employee_code, face_data, face_embedding, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (employee_code) DO UPDATE
		SET face_data = EXCLUDED.face_data,
		    face_embedding = EXCLUDED.face_embedding
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, template.ID, template.EmployeeCode, template.FaceData, template.FaceEmbedding)
	return err
}

func (r *TimekeepingRepositoryImpl) GetFaceTemplates(ctx context.Context) ([]entity.EmployeeFaceTemplate, error) {
	var list []entity.EmployeeFaceTemplate
	query := `
		SELECT id, employee_code, face_data, face_embedding, created_at
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
		SELECT id, employee_code, face_data, face_embedding, created_at
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

func (r *TimekeepingRepositoryImpl) GetAttendanceLogsList(ctx context.Context, employeeCodeFilter *string, employeeIDFilter *uuid.UUID, departmentIDFilter *uuid.UUID, dateFilter *time.Time, qp params.QueryParams) ([]entity.AttendanceLog, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	if offset < 0 {
		offset = 0
	}
	pageSize := qp.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	baseQuery := `
		FROM attendance_logs l
		JOIN user_profiles u ON l.employee_code = u.code
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if employeeCodeFilter != nil && *employeeCodeFilter != "" {
		conditions = append(conditions, fmt.Sprintf("l.employee_code = $%d", argIndex))
		args = append(args, *employeeCodeFilter)
		argIndex++
	}

	if employeeIDFilter != nil {
		conditions = append(conditions, fmt.Sprintf("u.id = $%d", argIndex))
		args = append(args, *employeeIDFilter)
		argIndex++
	}

	if departmentIDFilter != nil {
		conditions = append(conditions, fmt.Sprintf("u.department_id = $%d", argIndex))
		args = append(args, *departmentIDFilter)
		argIndex++
	}

	if dateFilter != nil {
		conditions = append(conditions, fmt.Sprintf("l.timestamp::date = $%d::date", argIndex))
		args = append(args, *dateFilter)
		argIndex++
	}

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR u.code ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+qp.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(l.id) %s %s", baseQuery, whereClause)
	err := r.DB.SQLx().GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT l.id, l.employee_code, l.timestamp, l.location_gps, l.device_id, l.created_at
		%s
		%s
		ORDER BY l.timestamp DESC
		LIMIT %d OFFSET %d
	`, baseQuery, whereClause, pageSize, offset)

	var list []entity.AttendanceLog
	err = r.DB.SQLx().SelectContext(ctx, &list, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return list, totalItems, nil
}

// ==========================================
// Leave Request Repository Methods
// ==========================================

func (r *TimekeepingRepositoryImpl) CreateLeaveRequest(ctx context.Context, req *entity.LeaveRequest) error {
	query := `
		INSERT INTO leave_requests (id, employee_id, start_date, end_date, days_requested, reason, approved_by, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query,
		req.ID, req.EmployeeID, req.StartDate, req.EndDate,
		req.DaysRequested, req.Reason, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetLeaveRequestByID(ctx context.Context, id uuid.UUID) (*entity.LeaveRequest, error) {
	var req entity.LeaveRequest
	query := `
		SELECT id, employee_id, start_date, end_date, days_requested, reason, approved_by, status, created_at, updated_at
		FROM leave_requests
		WHERE id = $1
	`
	err := r.DB.SQLx().GetContext(ctx, &req, query, id)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *TimekeepingRepositoryImpl) UpdateLeaveRequest(ctx context.Context, req *entity.LeaveRequest) error {
	query := `
		UPDATE leave_requests
		SET approved_by = $2, status = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query, req.ID, req.ApprovedBy, req.Status)
	return err
}

func (r *TimekeepingRepositoryImpl) GetLeaveRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID, qp params.QueryParams) ([]entity.LeaveRequest, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	if offset < 0 {
		offset = 0
	}
	pageSize := qp.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	baseQuery := `
		FROM leave_requests lr
		JOIN user_profiles u ON lr.employee_id = u.id
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if employeeID != nil {
		conditions = append(conditions, fmt.Sprintf("lr.employee_id = $%d", argIndex))
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

	if deptStr, ok := qp.Filters["department_id"]; ok && deptStr != "" {
		if deptID, err := uuid.Parse(deptStr); err == nil {
			conditions = append(conditions, fmt.Sprintf("u.department_id = $%d", argIndex))
			args = append(args, deptID)
			argIndex++
		}
	}

	if empStr, ok := qp.Filters["employee_id"]; ok && empStr != "" {
		if empID, err := uuid.Parse(empStr); err == nil {
			conditions = append(conditions, fmt.Sprintf("lr.employee_id = $%d", argIndex))
			args = append(args, empID)
			argIndex++
		}
	}

	if statusStr, ok := qp.Filters["status"]; ok && statusStr != "" && statusStr != "ALL" {
		conditions = append(conditions, fmt.Sprintf("lr.status = $%d", argIndex))
		args = append(args, statusStr)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(lr.id) " + baseQuery + whereClause
	var totalItems int
	err := r.DB.SQLx().GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT lr.id, lr.employee_id, lr.start_date, lr.end_date, lr.days_requested, lr.reason, lr.approved_by, lr.status, lr.created_at, lr.updated_at
	` + baseQuery + whereClause + ` ORDER BY lr.start_date DESC, lr.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var list []entity.LeaveRequest
	err = r.DB.SQLx().SelectContext(ctx, &list, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.LeaveRequest{}, 0, nil
		}
		return nil, 0, err
	}

	return list, totalItems, nil
}

// GetApprovedLeavesForPeriod lấy tất cả leave_requests đã APPROVED có ngày nghỉ giao với [start, end].
// Dùng trong CalculateTimesheets để map từng ngày nghỉ có phép.
func (r *TimekeepingRepositoryImpl) GetApprovedLeavesForPeriod(ctx context.Context, start, end time.Time) ([]entity.LeaveRequest, error) {
	var list []entity.LeaveRequest
	query := `
		SELECT id, employee_id, start_date, end_date, days_requested, reason, approved_by, status, created_at, updated_at
		FROM leave_requests
		WHERE status = 'APPROVED'
		  AND start_date <= $2
		  AND end_date >= $1
		ORDER BY employee_id, start_date ASC
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, start, end)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.LeaveRequest{}, nil
		}
		return nil, err
	}
	return list, nil
}

// ==========================================
// Leave Balance Repository Methods
// ==========================================

func (r *TimekeepingRepositoryImpl) GetLeaveBalance(ctx context.Context, employeeID uuid.UUID, year int) (*entity.LeaveBalance, error) {
	var balance entity.LeaveBalance
	query := `
		SELECT id, employee_id, year, accrued_days, used_days, balance, last_accrual_month, created_at, updated_at
		FROM leave_balances
		WHERE employee_id = $1 AND year = $2
	`
	err := r.DB.SQLx().GetContext(ctx, &balance, query, employeeID, year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &balance, nil
}

// UpsertLeaveBalance tạo mới hoặc cập nhật bản ghi leave_balances.
func (r *TimekeepingRepositoryImpl) UpsertLeaveBalance(ctx context.Context, b *entity.LeaveBalance) error {
	query := `
		INSERT INTO leave_balances (id, employee_id, year, accrued_days, used_days, balance, last_accrual_month, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (employee_id, year) DO UPDATE SET
			accrued_days       = EXCLUDED.accrued_days,
			used_days          = EXCLUDED.used_days,
			balance            = EXCLUDED.balance,
			last_accrual_month = EXCLUDED.last_accrual_month,
			updated_at         = NOW()
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query,
		b.ID, b.EmployeeID, b.Year, b.AccruedDays, b.UsedDays, b.Balance, b.LastAccrualMonth)
	return err
}

func (r *TimekeepingRepositoryImpl) GetAllLeaveBalancesForYear(ctx context.Context, year int) ([]entity.LeaveBalance, error) {
	var list []entity.LeaveBalance
	query := `
		SELECT id, employee_id, year, accrued_days, used_days, balance, last_accrual_month, created_at, updated_at
		FROM leave_balances
		WHERE year = $1
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, year)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.LeaveBalance{}, nil
		}
		return nil, err
	}
	return list, nil
}

// CountLeavePaidDaysInYear đếm tổng actual_work_day của những ngày có status = LEAVE_PAID trong năm.
// Hỗ trợ nửa ngày (actual_work_day = 0.5).
func (r *TimekeepingRepositoryImpl) CountLeavePaidDaysInYear(ctx context.Context, employeeID uuid.UUID, year int) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(actual_work_day), 0)
		FROM daily_attendance_sheets
		WHERE employee_id = $1
		  AND EXTRACT(YEAR FROM date) = $2
		  AND status = 'LEAVE_PAID'
	`
	err := r.DB.SQLx().GetContext(ctx, &total, query, employeeID, year)
	if err != nil {
		return 0, err
	}
	return total, nil
}
