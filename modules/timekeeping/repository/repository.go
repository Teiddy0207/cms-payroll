package repository

import (
	"cal-salary/core/params"
	"cal-salary/modules/timekeeping/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type TimekeepingRepository interface {
	CreateAttendanceLog(ctx context.Context, log *entity.AttendanceLog) error
	GetAttendanceLogs(ctx context.Context, employeeCode string, start, end time.Time) ([]entity.AttendanceLog, error)
	GetAttendanceLogsForCalculation(ctx context.Context, start, end time.Time) ([]entity.AttendanceLog, error)
	GetAttendanceLogsList(ctx context.Context, employeeCodeFilter *string, employeeIDFilter *uuid.UUID, departmentIDFilter *uuid.UUID, dateFilter *time.Time, qp params.QueryParams) ([]entity.AttendanceLog, int, error)

	UpsertDailyAttendanceSheet(ctx context.Context, sheet *entity.DailyAttendanceSheet) error
	GetDailyAttendanceSheets(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID, start, end time.Time, qp params.QueryParams) ([]entity.DailyAttendanceSheet, int, error)

	CreateExplanationRequest(ctx context.Context, req *entity.ExplanationRequest) error
	GetExplanationRequestByID(ctx context.Context, id uuid.UUID) (*entity.ExplanationRequest, error)
	UpdateExplanationRequest(ctx context.Context, req *entity.ExplanationRequest) error
	GetExplanationRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID) ([]entity.ExplanationRequest, error)

	CreateOTRequest(ctx context.Context, req *entity.OTRequest) error
	GetOTRequestByID(ctx context.Context, id uuid.UUID) (*entity.OTRequest, error)
	UpdateOTRequest(ctx context.Context, req *entity.OTRequest) error
	GetOTRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID) ([]entity.OTRequest, error)

	CreateFaceTemplate(ctx context.Context, template *entity.EmployeeFaceTemplate) error
	GetFaceTemplates(ctx context.Context) ([]entity.EmployeeFaceTemplate, error)
	GetFaceTemplateByCode(ctx context.Context, code string) (*entity.EmployeeFaceTemplate, error)
	DeleteFaceTemplate(ctx context.Context, code string) error

	// Leave Requests
	CreateLeaveRequest(ctx context.Context, req *entity.LeaveRequest) error
	GetLeaveRequestByID(ctx context.Context, id uuid.UUID) (*entity.LeaveRequest, error)
	UpdateLeaveRequest(ctx context.Context, req *entity.LeaveRequest) error
	GetLeaveRequests(ctx context.Context, employeeID *uuid.UUID, departmentID *uuid.UUID) ([]entity.LeaveRequest, error)
	// GetApprovedLeavesForPeriod lấy tất cả leave đã APPROVED trong khoảng thời gian
	// (dùng trong CalculateTimesheets để xác định ngày được phép nghỉ có lương).
	GetApprovedLeavesForPeriod(ctx context.Context, start, end time.Time) ([]entity.LeaveRequest, error)

	// Leave Balances
	GetLeaveBalance(ctx context.Context, employeeID uuid.UUID, year int) (*entity.LeaveBalance, error)
	UpsertLeaveBalance(ctx context.Context, balance *entity.LeaveBalance) error
	GetAllLeaveBalancesForYear(ctx context.Context, year int) ([]entity.LeaveBalance, error)
	// CountLeavePaidDaysInYear đếm số ngày LEAVE_PAID trong daily_attendance_sheets của nhân viên trong năm.
	// Dùng để tính lại used_days sau CalculateTimesheets.
	CountLeavePaidDaysInYear(ctx context.Context, employeeID uuid.UUID, year int) (float64, error)
}
