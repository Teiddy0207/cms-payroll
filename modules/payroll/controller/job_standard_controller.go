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
// Job Standards Echo Handlers
// ==========================================

func (ctrl *PayrollController) CreateJobStandard(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CreateJobStandardRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	resp, appErr := ctrl.Service.CreateJobStandard(ctx, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Create job standard success")
}

func (ctrl *PayrollController) GetJobStandards(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := ctrl.Service.GetJobStandards(ctx, *qp)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch job standards success")
}

func (ctrl *PayrollController) GetJobStandardByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := ctrl.Service.GetJobStandardByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch job standard success")
}

func (ctrl *PayrollController) UpdateJobStandard(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.UpdateJobStandardRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateJobStandard(ctx, id, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update job standard success")
}

func (ctrl *PayrollController) DeleteJobStandard(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := ctrl.Service.DeleteJobStandard(ctx, id)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Delete job standard success")
}
