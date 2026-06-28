package auth

import (
	"cal-salary/core/cache"
	"cal-salary/core/database"
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	"cal-salary/modules/auth/controller"
	"cal-salary/modules/auth/repository"
	"cal-salary/modules/auth/router"
	"cal-salary/modules/auth/service"

	"github.com/labstack/echo/v4"
)

type AuthModule struct {
	Service service.AuthServiceInterface
	repo    repository.AuthRepositoryInterface
}

// Init initializes the auth module's repository and service layers
func Init(db database.Database, rCache cache.Cache) *AuthModule {
	repo := repository.NewAuthRepository(db)
	svc := service.NewAuthService(repo, rCache)
	return &AuthModule{
		Service: svc,
		repo:    repo,
	}
}

// SetupRouter initializes controller and sets up auth routes with middleware
func (m *AuthModule) SetupRouter(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	ctrl := controller.NewAuthController(m.Service, m.repo)
	router.NewAuthRouter(*ctrl).Setup(e, middlewareInstance, activityLogSvc)
}
