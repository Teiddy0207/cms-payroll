package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/modules/timekeeping/dto"
	"context"

	"github.com/google/uuid"
)

type TimekeepingService interface {
	ProcessCheckIn(ctx context.Context, req *dto.CheckinRequest) (*dto.CheckinResponse, *errors.AppError)
	GetDailyAttendanceSheets(ctx context.Context, userID uuid.UUID, qp params.QueryParams) (*dto.PaginatedDailyAttendanceResponse, *errors.AppError)
	GetAttendanceLogsList(ctx context.Context, userID uuid.UUID, qp params.QueryParams) (*dto.PaginatedAttendanceLogResponse, *errors.AppError)
	CalculateTimesheets(ctx context.Context, req *dto.CalculateTimesheetRequest) *errors.AppError

	CreateExplanationRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateExplanationRequest) (*dto.ExplanationRequestResponse, *errors.AppError)
	GetExplanationRequests(ctx context.Context, userID uuid.UUID) ([]dto.ExplanationRequestResponse, *errors.AppError)
	UpdateExplanationRequestStatus(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *dto.UpdateRequestStatus) *errors.AppError

	CreateOTRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateOTRequest) (*dto.OTRequestResponse, *errors.AppError)
	GetOTRequests(ctx context.Context, userID uuid.UUID) ([]dto.OTRequestResponse, *errors.AppError)
	UpdateOTRequestStatus(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *dto.UpdateRequestStatus) *errors.AppError

	RegisterFaceTemplate(ctx context.Context, req *dto.RegisterFaceRequest) *errors.AppError
	GetFaceTemplates(ctx context.Context) ([]dto.FaceTemplateResponse, *errors.AppError)
	DeleteFaceTemplate(ctx context.Context, code string) *errors.AppError

	// EnsureStreamAndConsumer creates/updates the JetStream stream and durable
	// consumer used for the check-in pipeline. Call once at startup.
	EnsureStreamAndConsumer(ctx context.Context) error
	// StartCheckinConsumer starts consuming check-in events and persisting them
	// to attendance_logs. Blocks until the consumer context is set up; message
	// handling itself runs in the background.
	StartCheckinConsumer(ctx context.Context) error
}
