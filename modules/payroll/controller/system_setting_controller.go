package controller

import (
	"cal-salary/core/errors"
	"cal-salary/modules/payroll/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)


func (ctrl *PayrollController) GetSystemSetting(c echo.Context) error {
	ctx := c.Request().Context()
	key := c.Param("key")
	if key == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Key is required", nil)
	}

	resp, appErr := ctrl.Service.GetSystemSetting(ctx, key)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch system setting success")
}

func (ctrl *PayrollController) UpdateSystemSetting(c echo.Context) error {
	ctx := c.Request().Context()
	key := c.Param("key")
	if key == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Key is required", nil)
	}

	req := new(dto.UpdateSystemSettingRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateSystemSetting(ctx, key, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update system setting success")
}

func (ctrl *PayrollController) GetAllSystemSettings(c echo.Context) error {
	ctx := c.Request().Context()
	resp, appErr := ctrl.Service.GetAllSystemSettings(ctx)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch system settings success")
}