package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/activity_log/dto"
	"cal-salary/modules/activity_log/entity"
	"cal-salary/modules/activity_log/mapper"
	"context"
	"time"
)

// Log asynchronously records an activity log entry.
func (s *ActivityLogService) Log(ctx context.Context, log *entity.ActivityLog) {
	// Execute in a separate goroutine so we don't block the caller (typically an HTTP request).
	go func() {
		// Use a background context with a timeout so standard request cancellation
		// doesn't abort the DB insertion if the HTTP request has already completed.
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.repo.Insert(bgCtx, log); err != nil {
			logger.Error("ActivityLogService:Log:Error inserting activity log", "error", err)
		}
	}()
}

// GetActivityLogs retrieves paginated activity logs based on query params.
func (s *ActivityLogService) GetActivityLogs(ctx context.Context, p params.QueryParams) (*dto.PaginatedActivityLogDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.repo.GetActivityLogs(ctx, p)
	if err != nil {
		logger.Error("ActivityLogService:GetActivityLogs:Error", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get activity logs", err)
	}

	return mapper.ToActivityLogPaginationDTO(result), nil
}
