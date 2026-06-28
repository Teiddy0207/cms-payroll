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
// Job Positions Echo Handlers
// ==========================================

func (ctrl *PayrollController) CreateJobPosition(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CreateJobPositionRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	resp, appErr := ctrl.Service.CreateJobPosition(ctx, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Create job position success")
}

func (ctrl *PayrollController) GetJobPositions(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := ctrl.Service.GetJobPositions(ctx, *qp)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch job positions success")
}

func (ctrl *PayrollController) GetJobPositionByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := ctrl.Service.GetJobPositionByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch job position success")
}

func (ctrl *PayrollController) UpdateJobPosition(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.UpdateJobPositionRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateJobPosition(ctx, id, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update job position success")
}

func (ctrl *PayrollController) DeleteJobPosition(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := ctrl.Service.DeleteJobPosition(ctx, id)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Delete job position success")
}
