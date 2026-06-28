package router

import (
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	"cal-salary/modules/auth/controller"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	AuthController controller.AuthController
}

func NewAuthRouter(authController controller.AuthController) *AuthRouter {
	return &AuthRouter{
		AuthController: authController,
	}
}

func (r *AuthRouter) Setup(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	v1 := e.Group("/api/v1")

	privateRoutes := v1.Group("/private")
	publicRoutes := v1.Group("/public")

	authRoutes := privateRoutes.Group("/auth")

	// Apply activity log middleware to all auth routes
	authRoutes.Use(middlewareInstance.ActivityLogMiddleware(activityLogSvc, "auth"))

	authPublicRoutes := publicRoutes.Group("/auth")

	authPublicRoutes.POST("/register", r.AuthController.Register)
	authPublicRoutes.POST("/login", r.AuthController.Login)
	authPublicRoutes.POST("/logout", r.AuthController.Logout)

	authPublicRoutes.POST("/forgot-password", r.AuthController.ForgotPassword)
	authPublicRoutes.POST("/verify-otp", r.AuthController.VerifyOTP)
	authPublicRoutes.POST("/reset-password", r.AuthController.ResetPassword)
	authPublicRoutes.POST("/send-otp-change-password", r.AuthController.SendOTPChangePassword)
	authPublicRoutes.POST("/change-password", r.AuthController.ChangePassword)

	authPublicRoutes.POST("/update-password", r.AuthController.ChangePassword, middlewareInstance.AuthMiddleware())
}
