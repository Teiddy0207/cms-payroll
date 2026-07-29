package router

import (
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	"cal-salary/modules/timekeeping/controller"

	"github.com/labstack/echo/v4"
)

type TimekeepingRouter struct {
	TimekeepingController *controller.TimekeepingController
}

func NewTimekeepingRouter(timekeepingController *controller.TimekeepingController) *TimekeepingRouter {
	return &TimekeepingRouter{
		TimekeepingController: timekeepingController,
	}
}

func (r *TimekeepingRouter) Setup(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	v1 := e.Group("/api/v1")
	privateRoutes := v1.Group("/private")

	privateRoutes.Use(middlewareInstance.AuthMiddleware())

	timekeepingRoutes := privateRoutes.Group("/timekeeping")
	timekeepingRoutes.Use(middlewareInstance.ActivityLogMiddleware(activityLogSvc, "timekeeping"))

	timekeepingRoutes.POST("/checkin", r.TimekeepingController.ProcessCheckIn) // Open to all authenticated users
	timekeepingRoutes.GET("/sheets", r.TimekeepingController.GetDailyAttendanceSheets, middlewareInstance.PermissionMiddleware("timekeepingSheet::read"))
	timekeepingRoutes.GET("/logs", r.TimekeepingController.GetAttendanceLogs, middlewareInstance.PermissionMiddleware("dailyTimekeeping::read"))
	timekeepingRoutes.POST("/sheets/calculate", r.TimekeepingController.CalculateTimesheets, middlewareInstance.PermissionMiddleware("timekeepingSheet::edit"))

	timekeepingRoutes.POST("/explanations", r.TimekeepingController.CreateExplanationRequest) // Open for employees to request
	timekeepingRoutes.GET("/explanations", r.TimekeepingController.GetExplanationRequests, middlewareInstance.PermissionMiddleware("dailyTimekeeping::read"))
	timekeepingRoutes.PUT("/explanations/:id/status", r.TimekeepingController.UpdateExplanationRequestStatus, middlewareInstance.PermissionMiddleware("dailyTimekeeping::review"))

	timekeepingRoutes.POST("/ot-requests", r.TimekeepingController.CreateOTRequest) // Open for employees to request
	timekeepingRoutes.GET("/ot-requests", r.TimekeepingController.GetOTRequests, middlewareInstance.PermissionMiddleware("dailyTimekeeping::read"))
	timekeepingRoutes.PUT("/ot-requests/:id/status", r.TimekeepingController.UpdateOTRequestStatus, middlewareInstance.PermissionMiddleware("dailyTimekeeping::review"))

	timekeepingRoutes.POST("/faces", r.TimekeepingController.RegisterFaceTemplate, middlewareInstance.PermissionMiddleware("dailyTimekeeping::edit"))
	timekeepingRoutes.GET("/faces", r.TimekeepingController.GetFaceTemplates, middlewareInstance.PermissionMiddleware("dailyTimekeeping::read"))
	timekeepingRoutes.DELETE("/faces/:code", r.TimekeepingController.DeleteFaceTemplate, middlewareInstance.PermissionMiddleware("dailyTimekeeping::delete"))
}
