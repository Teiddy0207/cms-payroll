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
}
