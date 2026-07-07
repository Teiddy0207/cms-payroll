package timekeeping

import (
	"cal-salary/core/cache"
	"cal-salary/core/database"
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	payrollRepo "cal-salary/modules/payroll/repository"
	"cal-salary/modules/timekeeping/controller"
	"cal-salary/modules/timekeeping/repository"
	"cal-salary/modules/timekeeping/router"
	"cal-salary/modules/timekeeping/service"

	"github.com/labstack/echo/v4"
)

type TimekeepingModule struct {
	Service service.TimekeepingService
	repo    repository.TimekeepingRepository
}

func Init(db database.Database, redisCache *cache.Cache) *TimekeepingModule {
	repo := repository.NewTimekeepingRepository(db)
	pRepo := payrollRepo.NewPayrollRepository(db)
	svc := service.NewTimekeepingService(repo, pRepo, redisCache)
	return &TimekeepingModule{
		Service: svc,
		repo:    repo,
	}
}

func (m *TimekeepingModule) SetupRouter(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	ctrl := controller.NewTimekeepingController(m.Service)
	router.NewTimekeepingRouter(ctrl).Setup(e, middlewareInstance, activityLogSvc)
}
