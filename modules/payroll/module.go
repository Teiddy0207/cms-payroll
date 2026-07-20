package payroll

import (
	"cal-salary/core/cache"
	"cal-salary/core/database"
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	"cal-salary/modules/payroll/controller"
	"cal-salary/modules/payroll/repository"
	"cal-salary/modules/payroll/router"
	"cal-salary/modules/payroll/service"

	"github.com/labstack/echo/v4"
)

type PayrollModule struct {
	Service service.PayrollServiceInterface
	repo    repository.PayrollRepositoryInterface
}

func Init(db database.Database, redisCache *cache.Cache) *PayrollModule {
	repo := repository.NewPayrollRepository(db)
	svc := service.NewPayrollService(repo, redisCache)
	return &PayrollModule{
		Service: svc,
		repo:    repo,
	}
}

// SetupRouter sets up HTTP routing for all payroll endpoints
func (m *PayrollModule) SetupRouter(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	ctrl := controller.NewPayrollController(m.Service)
	router.NewPayrollRouter(ctrl).Setup(e, middlewareInstance, activityLogSvc)
}
