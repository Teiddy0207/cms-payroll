package controller

import (
	"cal-salary/core/errors"
	"cal-salary/modules/payroll/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)


func (ctrl *PayrollController) AssignStandardToPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	req := new(dto.AssignStandardRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.AssignStandardToPosition(ctx, positionID, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Assign standard to position success")
}

func (ctrl *PayrollController) GetStandardsByPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	resp, appErr := ctrl.Service.GetStandardsByPosition(ctx, positionID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch standards by position success")
}

func (ctrl *PayrollController) RemoveStandardFromPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	standardID, err := uuid.Parse(c.Param("standard_id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Standard ID format", nil)
	}

	appErr := ctrl.Service.RemoveStandardFromPosition(ctx, positionID, standardID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Remove standard from position success")
}

// ==========================================
// Job Position Competencies Handlers
// ==========================================

func (ctrl *PayrollController) AssignCompetencyToPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	req := new(dto.AssignCompetencyRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.AssignCompetencyToPosition(ctx, positionID, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Assign competency to position success")
}

func (ctrl *PayrollController) GetCompetenciesByPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	resp, appErr := ctrl.Service.GetCompetenciesByPosition(ctx, positionID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch competencies by position success")
}

func (ctrl *PayrollController) RemoveCompetencyFromPosition(c echo.Context) error {
	ctx := c.Request().Context()
	positionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Position ID format", nil)
	}

	competencyID, err := uuid.Parse(c.Param("competency_id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Competency ID format", nil)
	}

	appErr := ctrl.Service.RemoveCompetencyFromPosition(ctx, positionID, competencyID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Remove competency from position success")
}

// ==========================================
// System Settings Handlers
// ==========================================

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

// ==========================================
// Calculator Handlers
// ==========================================

func (ctrl *PayrollController) PreviewSalary(c echo.Context) error {
	ctx := c.Request().Context()
	employeeID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Employee ID format", nil)
	}

	period := c.QueryParam("period")
	if period == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "period query parameter is required", nil)
	}

	resp, appErr := ctrl.Service.PreviewSalary(ctx, employeeID, period)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch salary preview success")
}
