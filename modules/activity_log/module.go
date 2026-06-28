package activity_log

import (
	"cal-salary/core/database"
	"cal-salary/core/middleware"
	"cal-salary/modules/activity_log/controller"
	"cal-salary/modules/activity_log/repository"
	"cal-salary/modules/activity_log/router"
	"cal-salary/modules/activity_log/service"

	"github.com/labstack/echo/v4"
)

// Init initializes the activity log module and returns the service instance
// so it can be injected into other modules.
func Init(e *echo.Echo, db database.Database, middlewareInstance *middleware.Middleware) service.ActivityLogServiceInterface {
	repo := repository.NewActivityLogRepository(db)
	svc := service.NewActivityLogService(repo)
	ctrl := controller.NewActivityLogController(svc)
	router.NewActivityLogRouter(ctrl).Setup(e, middlewareInstance)

	return svc
}
