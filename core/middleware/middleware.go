package middleware

import (
	"cal-salary/core/constants"
	"cal-salary/core/controller"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/utils"
	activityEntity "cal-salary/modules/activity_log/entity"
	activitySvc "cal-salary/modules/activity_log/service"
	"cal-salary/modules/auth/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	controller.BaseController
	AuthService service.AuthServiceInterface
}

func NewMiddleware(authService service.AuthServiceInterface) *Middleware {
	return &Middleware{
		BaseController: controller.NewBaseController(),
		AuthService:    authService,
	}
}
func (m *Middleware) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header or query param
			authHeader := c.Request().Header.Get("Authorization")
			var tokenStr string
			if authHeader != "" {
				parts := strings.Fields(authHeader)
				if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
					return m.Unauthorized(errors.ErrInvalidTokenFormat, "invalid token format")
				}
				tokenStr = strings.TrimSpace(parts[1])
			} else if queryToken := c.QueryParam("token"); queryToken != "" {
				tokenStr = queryToken
			} else {
				return m.Unauthorized(errors.ErrMissingAuthorizationHeader, "missing authorization header")
			}

			// Validate token
			claims, err := utils.ValidateAndParseToken(tokenStr)
			if err != nil {
				logger.Error("AuthMiddleware:ValidateAndParseToken:Error:", err)
				return m.Unauthorized(errors.ErrInvalidTokenFormat, "invalid token: "+err.Error())
			}

			// Set user claims in context
			c.Set(constants.ContextTokenData, claims)
			return next(c)
		}
	}
}

func (m *Middleware) PermissionMiddleware(requiredPermissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Lấy thông tin user từ context (đã được set bởi AuthMiddleware)
			userData := c.Get(constants.ContextTokenData)
			if userData == nil {
				return m.Unauthorized(errors.ErrMissingAuthorizationHeader, "missing authorization header")
			}

			// Parse user claims từ token
			claims, ok := userData.(*utils.TokenClaims)
			if !ok {
				return m.Unauthorized(errors.ErrInvalidTokenFormat, "invalid token format")
			}

			// Nếu không có permission nào được yêu cầu, cho phép truy cập
			if len(requiredPermissions) == 0 {
				return next(c)
			}

			// Lấy danh sách permissions trực tiếp của user (bao gồm từ roles và permissions riêng)
			ctx := c.Request().Context()

			// // Thử lấy từ cache trước
			// userPermissions, err := m.AuthService.PrivateGetPermissionsByUserIDFromCache(ctx, claims.UserID)
			// if err != nil {
			// 	// Nếu có lỗi thực sự (không phải cache miss), log và lấy từ database
			// 	logger.Error("PermissionMiddleware:Cache error, fetching from database:", err)
			// 	userPermissions, err = m.AuthService.PrivateGetPermissionsByUserID(ctx, claims.UserID)
			// 	if err != nil {
			// 		logger.Error("PermissionMiddleware:PrivateGetPermissions:Error:", err)
			// 		return m.Unauthorized(errors.ErrUnauthorized, "user has no permissions")
			// 	}
			// } else if userPermissions == nil {
			// 	// Cache miss (không có trong cache), lấy từ database
			// 	logger.Info("PermissionMiddleware:Cache miss, fetching from database")
			// 	userPermissions, err = m.AuthService.PrivateGetPermissionsByUserID(ctx, claims.UserID)
			// 	if err != nil {
			// 		logger.Error("PermissionMiddleware:PrivateGetPermissions:Error:", err)
			// 		return m.Unauthorized(errors.ErrUnauthorized, "user has no permissions")
			// 	}
			// }
			userPermissions, err := m.AuthService.PrivateGetPermissionsByUserID(ctx, claims.UserID)
			if err != nil {
				logger.Error("PermissionMiddleware:PrivateGetPermissions:Error:", err)
				return m.Unauthorized(errors.ErrUnauthorized, "user has no permissions")
			}

			if userPermissions == nil {
				logger.Error("PermissionMiddleware:PrivateGetPermissions:Error: userPermissions is nil")
				return m.Unauthorized(errors.ErrUnauthorized, "user has no permissions")
			}

			// Kiểm tra xem user có ít nhất một trong các permissions được yêu cầu không
			for _, requiredPerm := range requiredPermissions {
				for _, userPerm := range *userPermissions {
					// Hỗ trợ cả slug (ví dụ: "user::create", "user:read") và format "resource:action"
					permissionString := userPerm.Resource + ":" + string(userPerm.Action)
					if requiredPerm == userPerm.Slug || permissionString == requiredPerm {
						return next(c)
					}
				}
			}

			// userPermissions, err = m.AuthService.PrivateGetPermissionsByUserID(ctx, claims.UserID)
			// if err == nil && userPermissions != nil {
			// 	for _, requiredPerm := range requiredPermissions {
			// 		for _, userPerm := range *userPermissions {
			// 			permissionString := userPerm.Resource + ":" + string(userPerm.Action)
			// 			if requiredPerm == userPerm.Slug || permissionString == requiredPerm {
			// 				return next(c)
			// 			}
			// 		}
			// 	}
			// }

			// Nếu không có permission nào khớp, trả về lỗi Forbidden
			return m.Unauthorized(errors.ErrForbidden, "user does not have permission to access this resource")
		}
	}
}

func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Record start time
			start := time.Now()

			req := c.Request()
			res := c.Response()

			if err := next(c); err != nil {
				c.Error(err)
			}

			// Calculate latency
			latency := time.Since(start)

			// Log request details
			logger.Info("Request",
				"method", req.Method,
				"uri", req.RequestURI,
				"status", res.Status,
				"remote_ip", c.RealIP(),
				"user_agent", req.UserAgent(),
				"latency", latency.String(),
			)

			return nil
		}
	}
}

func (m *Middleware) ActivityLogMiddleware(svc activitySvc.ActivityLogServiceInterface, module string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Execute the handler first
			err := next(c)

			method := c.Request().Method
			if method == http.MethodGet || method == http.MethodOptions {
				// Don't log read requests
				return err
			}

			var userID *uuid.UUID
			if userData := c.Get(constants.ContextTokenData); userData != nil {
				if claims, ok := userData.(*utils.TokenClaims); ok && claims.UserID != uuid.Nil {
					userID = &claims.UserID
				}
			}

			status := "success"
			var errMsg *string
			if err != nil {
				status = "failed"
				msg := err.Error()
				errMsg = &msg
			}

			action := "unknown"
			switch method {
			case http.MethodPost:
				action = "create"
			case http.MethodPut, http.MethodPatch:
				action = "update"
			case http.MethodDelete:
				action = "delete"
			}
			fmt.Println("ActivityLogMiddleware:action:", action)
			endpoint := c.Request().RequestURI

			// Auto-detect entity type from query param "type" or URL path (e.g. /api/v1/private/auth/roles -> roles)
			entityType := c.QueryParam("type")
			if entityType == "" {
				entityType = "unknown"
				pathParts := strings.Split(strings.TrimPrefix(c.Request().URL.Path, "/"), "/")
				if len(pathParts) >= 5 {
					entityType = pathParts[4]
				}
			}

			ip := c.RealIP()
			ua := c.Request().UserAgent()

			// Read rich data set by the service/controller layer (if any)
			logCtx := getActivityLogContext(c)

			logEntry := &activityEntity.ActivityLog{
				UserID:       userID,
				Module:       module,
				Action:       action,
				EntityType:   entityType,
				HTTPMethod:   &method,
				Endpoint:     &endpoint,
				IPAddress:    &ip,
				UserAgent:    &ua,
				Status:       status,
				ErrorMessage: errMsg,
			}

			if logCtx != nil {
				logEntry.EntityID = logCtx.EntityID
				logEntry.UserProfileID = logCtx.UserProfileID
				logEntry.Description = logCtx.Description
				logEntry.FieldsChanged = logCtx.FieldsChanged
				logEntry.OldData = logCtx.OldData
				logEntry.NewData = logCtx.NewData
			}

			// If UserProfileID is missing, try fetching it from table user_profiles via AuthService
			// if logEntry.UserProfileID == nil && userID != nil && m.AuthService != nil {
			// 	profile, err := m.AuthService.PrivateGetUserProfileByUserID(c.Request().Context(), *userID)
			// 	if err == nil && profile != nil {
			// 		logEntry.UserProfileID = &profile.ID
			// 	}
			// }

			svc.Log(c.Request().Context(), logEntry)

			return err
		}
	}
}
