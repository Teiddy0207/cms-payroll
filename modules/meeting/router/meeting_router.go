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

	// Static routes first before wildcard :id
	meetingGroup.GET("/notifications/stream", ctrl.StreamNotifications)
	meetingGroup.POST("/signal", ctrl.SendWebRTCSignal)
	meetingGroup.GET("", ctrl.GetMeetings)
	meetingGroup.GET("/", ctrl.GetMeetings)
	meetingGroup.POST("", ctrl.CreateMeeting)
	meetingGroup.POST("/", ctrl.CreateMeeting)
	meetingGroup.POST("/:id/rsvp", ctrl.UpdateRSVP)
	meetingGroup.PUT("/:id", ctrl.UpdateMeeting)
	meetingGroup.POST("/:id/complete", ctrl.CompleteMeeting) // kết thúc họp + AI summarize
	meetingGroup.POST("/:id/summary", ctrl.SaveSummary)
	meetingGroup.GET("/:id/summary", ctrl.GetSummary)
	meetingGroup.GET("/:id", ctrl.GetMeetingByID)
}

