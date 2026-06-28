package controller

import (
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ==========================================
// User Profiles Echo Handlers
// ==========================================

func (ctrl *PayrollController) CreateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CreateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	resp, appErr := ctrl.Service.CreateUserProfile(ctx, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Create user profile success")
}

func (ctrl *PayrollController) GetUserProfiles(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := ctrl.Service.GetUserProfiles(ctx, *qp)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch user profiles success")
}

func (ctrl *PayrollController) GetUserProfileByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := ctrl.Service.GetUserProfileByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch user profile success")
}

func (ctrl *PayrollController) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.UpdateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateUserProfile(ctx, id, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update user profile success")
}

func (ctrl *PayrollController) DeleteUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := ctrl.Service.DeleteUserProfile(ctx, id)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Delete user profile success")
}
