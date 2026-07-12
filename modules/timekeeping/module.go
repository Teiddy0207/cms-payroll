package timekeeping

import (
	"cal-salary/core/cache"
	"cal-salary/core/config"
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"cal-salary/core/messaging"
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	payrollRepo "cal-salary/modules/payroll/repository"
	"cal-salary/modules/timekeeping/controller"
	"cal-salary/modules/timekeeping/repository"
	"cal-salary/modules/timekeeping/router"
	"cal-salary/modules/timekeeping/service"
	"context"

	"github.com/labstack/echo/v4"
)

type TimekeepingModule struct {
	Service service.TimekeepingService
	repo    repository.TimekeepingRepository
}

func Init(db database.Database, redisCache *cache.Cache, natsClient *messaging.NatsClient, natsCfg config.NatsConfig) *TimekeepingModule {
	repo := repository.NewTimekeepingRepository(db)
	pRepo := payrollRepo.NewPayrollRepository(db)
	svc := service.NewTimekeepingService(repo, pRepo, redisCache, natsClient, natsCfg.StreamName, natsCfg.CheckinSubject, natsCfg.DurableConsumer)

	ctx := context.Background()
	if err := svc.EnsureStreamAndConsumer(ctx); err != nil {
		logger.Error("timekeeping: failed to set up JetStream checkin stream/consumer", "error", err)
	} else if err := svc.StartCheckinConsumer(ctx); err != nil {
		logger.Error("timekeeping: failed to start checkin consumer", "error", err)
	}

	return &TimekeepingModule{
		Service: svc,
		repo:    repo,
	}
}

func (m *TimekeepingModule) SetupRouter(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	ctrl := controller.NewTimekeepingController(m.Service)
	router.NewTimekeepingRouter(ctrl).Setup(e, middlewareInstance, activityLogSvc)
}
