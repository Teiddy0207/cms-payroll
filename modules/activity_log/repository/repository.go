package repository

import (
	"cal-salary/core/database"
	"cal-salary/core/params"
	"cal-salary/modules/activity_log/entity"
	"context"
)

// ActivityLogRepositoryInterface defines the contract for the activity-log repository.
type ActivityLogRepositoryInterface interface {
	Insert(ctx context.Context, log *entity.ActivityLog) error
	GetActivityLogs(ctx context.Context, p params.QueryParams) (*entity.PaginatedActivityLogEntity, error)
}

// ActivityLogRepository is the concrete implementation backed by a PostgreSQL database.
type ActivityLogRepository struct {
	DB database.Database
}

// NewActivityLogRepository creates a new ActivityLogRepository.
func NewActivityLogRepository(db database.Database) ActivityLogRepositoryInterface {
	return &ActivityLogRepository{DB: db}
}
