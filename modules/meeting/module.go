package meeting

import (
	"cal-salary/core/database"
	"cal-salary/core/google"
	"cal-salary/core/messaging"
	"cal-salary/core/middleware"
	"cal-salary/modules/meeting/controller"
	"cal-salary/modules/meeting/repository"
	"cal-salary/modules/meeting/router"
	"cal-salary/modules/meeting/service"

	"github.com/labstack/echo/v4"
)

type MeetingModule struct {
	Controller *controller.MeetingController
	Service    service.MeetingServiceInterface
}

func InitMeetingModule(db database.IDatabase, gcal *google.CalendarService, natsClient *messaging.NatsClient) *MeetingModule {
	repo := repository.NewMeetingRepository(db)
	svc := service.NewMeetingService(repo, gcal, natsClient)
	ctrl := controller.NewMeetingController(svc)

	return &MeetingModule{
		Controller: ctrl,
		Service:    svc,
	}
}

func (m *MeetingModule) SetupRouter(e *echo.Echo, middlewareInstance *middleware.Middleware) {
	v1 := e.Group("/api/v1")
	router.RegisterMeetingRoutes(v1, m.Controller, middlewareInstance.AuthMiddleware())
}

func (m *MeetingModule) RegisterRoutes(g *echo.Group, authMiddleware echo.MiddlewareFunc) {
	router.RegisterMeetingRoutes(g, m.Controller, authMiddleware)
}
