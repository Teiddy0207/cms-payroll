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

	timekeepingRoutes.POST("/checkin", r.TimekeepingController.ProcessCheckIn)
	timekeepingRoutes.GET("/sheets", r.TimekeepingController.GetDailyAttendanceSheets)
	timekeepingRoutes.POST("/sheets/calculate", r.TimekeepingController.CalculateTimesheets)

	timekeepingRoutes.POST("/explanations", r.TimekeepingController.CreateExplanationRequest)
	timekeepingRoutes.GET("/explanations", r.TimekeepingController.GetExplanationRequests)
	timekeepingRoutes.PUT("/explanations/:id/status", r.TimekeepingController.UpdateExplanationRequestStatus)

	timekeepingRoutes.POST("/ot-requests", r.TimekeepingController.CreateOTRequest)
	timekeepingRoutes.GET("/ot-requests", r.TimekeepingController.GetOTRequests)
	timekeepingRoutes.PUT("/ot-requests/:id/status", r.TimekeepingController.UpdateOTRequestStatus)

	timekeepingRoutes.POST("/faces", r.TimekeepingController.RegisterFaceTemplate)
	timekeepingRoutes.GET("/faces", r.TimekeepingController.GetFaceTemplates)
	timekeepingRoutes.DELETE("/faces/:code", r.TimekeepingController.DeleteFaceTemplate)
}
