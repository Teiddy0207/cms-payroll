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
	privateRoutes.Use(middlewareInstance.AuthMiddleware())
	publicRoutes := v1.Group("/public")

	authRoutes := privateRoutes.Group("/auth")

	// Apply activity log middleware to all auth routes
	authRoutes.Use(middlewareInstance.ActivityLogMiddleware(activityLogSvc, "auth"))

	// Private Auth - Roles
	authRoutes.POST("/roles", r.AuthController.CreateRole, middlewareInstance.PermissionMiddleware("role::edit"))
	authRoutes.GET("/roles", r.AuthController.GetRoles, middlewareInstance.PermissionMiddleware("role::read"))
	authRoutes.GET("/roles/:id", r.AuthController.GetRoleByID, middlewareInstance.PermissionMiddleware("role::read"))
	authRoutes.PUT("/roles/:id", r.AuthController.UpdateRole, middlewareInstance.PermissionMiddleware("role::edit"))
	authRoutes.DELETE("/roles/:id", r.AuthController.DeleteRole, middlewareInstance.PermissionMiddleware("role::delete"))

	// Private Auth - Permissions
	authRoutes.POST("/permissions", r.AuthController.CreatePermission, middlewareInstance.PermissionMiddleware("permission::edit"))
	authRoutes.GET("/permissions", r.AuthController.GetPermissions, middlewareInstance.PermissionMiddleware("permission::read"))
	authRoutes.GET("/permissions/:id", r.AuthController.GetPermissionByID, middlewareInstance.PermissionMiddleware("permission::read"))
	authRoutes.PUT("/permissions/:id", r.AuthController.UpdatePermission, middlewareInstance.PermissionMiddleware("permission::edit"))
	authRoutes.DELETE("/permissions/:id", r.AuthController.DeletePermission, middlewareInstance.PermissionMiddleware("permission::delete"))

	// Private Auth - Assignments
	authRoutes.POST("/roles/assign-permission", r.AuthController.AssignPermissionToRole, middlewareInstance.PermissionMiddleware("permission::edit"))
	authRoutes.POST("/users/assign-role", r.AuthController.AssignRoleToUser, middlewareInstance.PermissionMiddleware("role::edit"))
	authRoutes.POST("/users/assign-permission", r.AuthController.AssignPermissionToUser, middlewareInstance.PermissionMiddleware("permission::edit"))
	authRoutes.GET("/users/:id/permissions", r.AuthController.GetUserPermissions, middlewareInstance.PermissionMiddleware("permission::read"))
	authRoutes.GET("/me/permissions", r.AuthController.GetMePermissions)

	// Private Auth - Users CRUD
	authRoutes.POST("/users", r.AuthController.CreateUser, middlewareInstance.PermissionMiddleware("user::create"))
	authRoutes.GET("/users", r.AuthController.GetUsers, middlewareInstance.PermissionMiddleware("user::read"))
	authRoutes.GET("/users/:id", r.AuthController.GetUserByID, middlewareInstance.PermissionMiddleware("user::read"))
	authRoutes.PUT("/users/:id", r.AuthController.UpdateUser, middlewareInstance.PermissionMiddleware("user::edit"))
	authRoutes.DELETE("/users/:id", r.AuthController.DeleteUser, middlewareInstance.PermissionMiddleware("user::delete"))

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
