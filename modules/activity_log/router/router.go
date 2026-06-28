package router

import (
	"cal-salary/core/middleware"
	"cal-salary/modules/activity_log/controller"

	"github.com/labstack/echo/v4"
)

type ActivityLogRouter struct {
	ActivityLogController *controller.ActivityLogController
}

func NewActivityLogRouter(activityLogController *controller.ActivityLogController) *ActivityLogRouter {
	return &ActivityLogRouter{
		ActivityLogController: activityLogController,
	}
}

func (r *ActivityLogRouter) Setup(e *echo.Echo, middlewareInstance *middleware.Middleware) {
	v1 := e.Group("/api/v1")
	privateRoutes := v1.Group("/private")

	activityLogRoutes := privateRoutes.Group("/activity-logs")

	// Apply authentication and permission middleware
	// We assume a generic 'activity-log::read' permission, or whatever suits the system
	activityLogRoutes.GET("", r.ActivityLogController.GetActivityLogs, middlewareInstance.AuthMiddleware(), middlewareInstance.PermissionMiddleware()) //"activity-log::read"
}
