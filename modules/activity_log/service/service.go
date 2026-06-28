package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/modules/activity_log/dto"
	"cal-salary/modules/activity_log/entity"
	"cal-salary/modules/activity_log/repository"
	"context"
)

// ActivityLogServiceInterface defines the contract for logging activities.
type ActivityLogServiceInterface interface {
	Log(ctx context.Context, log *entity.ActivityLog)
	GetActivityLogs(ctx context.Context, p params.QueryParams) (*dto.PaginatedActivityLogDTO, *errors.AppError)
}

// ActivityLogService implements the activity logging logic.
type ActivityLogService struct {
	repo repository.ActivityLogRepositoryInterface
}

// NewActivityLogService creates a new ActivityLogService.
func NewActivityLogService(repo repository.ActivityLogRepositoryInterface) ActivityLogServiceInterface {
	return &ActivityLogService{repo: repo}
}
