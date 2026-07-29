package controller

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/core/utils"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/validator"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ==========================================
// Roles Endpoints
// ==========================================

func (controller *AuthController) CreateRole(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.RoleRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidateRoleRequest(ctx, controller.AuthRepo, req, nil)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	resp, appErr := controller.AuthService.PrivateCreateRole(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, resp, "Create role success")
}

func (controller *AuthController) GetRoles(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := controller.AuthService.PrivateGetRoles(ctx, *qp)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get roles success")
}

func (controller *AuthController) GetRoleByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := controller.AuthService.PrivateGetRoleByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get role by ID success")
}

func (controller *AuthController) UpdateRole(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.RoleRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidateRoleRequestUpdate(ctx, controller.AuthRepo, req, &id)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	appErr := controller.AuthService.PrivateUpdateRole(ctx, id, req)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Update role success")
}

func (controller *AuthController) DeleteRole(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := controller.AuthService.PrivateDeleteRole(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Delete role success")
}

// ==========================================
// Permissions Endpoints
// ==========================================

func (controller *AuthController) CreatePermission(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.PermissionRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidatePermissionRequest(req)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	appErr := controller.AuthService.PrivateCreatePermission(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Create permission success")
}

func (controller *AuthController) GetPermissions(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := controller.AuthService.PrivateGetPermissions(ctx, *qp)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get permissions success")
}

func (controller *AuthController) GetPermissionByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := controller.AuthService.PrivateGetPermissionByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get permission by ID success")
}

func (controller *AuthController) UpdatePermission(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.PermissionRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidatePermissionRequest(req)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	appErr := controller.AuthService.PrivateUpdatePermission(ctx, id, req)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Update permission success")
}

func (controller *AuthController) DeletePermission(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := controller.AuthService.PrivateDeletePermission(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Delete permission success")
}

// ==========================================
// Assignment Endpoints
// ==========================================

func (controller *AuthController) AssignRoleToUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.UserRoleRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	appErr := controller.AuthService.PrivateAssignRoleToUser(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Assign role to user success")
}

func (controller *AuthController) AssignPermissionToRole(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.RolePermissionRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidateAssignPermissionToRoleRequest(req)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	appErr := controller.AuthService.PrivateAssignPermissionToRole(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Assign permissions to role success")
}

func (controller *AuthController) AssignPermissionToUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.UserPermissionRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidateAssignPermissionToUserRequest(req)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	appErr := controller.AuthService.PrivateAssignPermissionToUser(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Assign permission to user success")
}

func (controller *AuthController) GetUserPermissions(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := controller.AuthService.PrivateGetPermissionsByUserID(ctx, id)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get user permissions success")
}

// ==========================================
// User Account Endpoints
// ==========================================

func (controller *AuthController) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	valRes := validator.ValidateCreateUserRequest(ctx, req, controller.AuthRepo)
	if valRes.HasError() {
		return controller.BadRequest(errors.ErrInvalidInput, "Validation failed", valRes)
	}

	resp, appErr := controller.AuthService.PrivateCreateUser(ctx, req)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, resp, "Create user success")
}

func (controller *AuthController) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := controller.AuthService.PrivateGetUsers(ctx, *qp)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get users success")
}

func (controller *AuthController) GetUserByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := controller.AuthService.PrivateGetUser(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get user success")
}

func (controller *AuthController) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.UserUpdateRequest)
	if err := c.Bind(req); err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid request body", nil)
	}

	appErr := controller.AuthService.PrivateUpdateUser(ctx, id, req)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Update user success")
}

func (controller *AuthController) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return controller.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := controller.AuthService.PrivateDeleteUser(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return controller.SuccessResponse(c, nil, "Delete user success")
}

func (controller *AuthController) GetMePermissions(c echo.Context) error {
	ctx := c.Request().Context()
	userData := c.Get(constants.ContextTokenData)
	if userData == nil {
		return controller.Unauthorized(errors.ErrMissingAuthorizationHeader, "missing authorization header", nil)
	}

	claims, ok := userData.(*utils.TokenClaims)
	if !ok {
		return controller.Unauthorized(errors.ErrInvalidTokenFormat, "invalid token format", nil)
	}

	resp, appErr := controller.AuthService.PrivateGetPermissionsByUserID(ctx, claims.UserID)
	if appErr != nil {
		return controller.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return controller.SuccessResponse(c, resp, "Get me permissions success")
}

