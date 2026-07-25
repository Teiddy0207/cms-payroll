package router

import (
	"cal-salary/modules/meeting/controller"

	"github.com/labstack/echo/v4"
)

func RegisterMeetingRoutes(g *echo.Group, ctrl *controller.MeetingController, authMiddleware echo.MiddlewareFunc) {
	meetingGroup := g.Group("/meetings")
	if authMiddleware != nil {
		meetingGroup.Use(authMiddleware)
	}

	meetingGroup.POST("", ctrl.CreateMeeting)
	meetingGroup.GET("", ctrl.GetMeetings)
	meetingGroup.POST("/:id/rsvp", ctrl.UpdateRSVP)
}
