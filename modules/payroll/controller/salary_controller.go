package controller

import (
	"cal-salary/core/errors"
	"cal-salary/modules/payroll/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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

func (ctrl *PayrollController) RunSalaryCalculation(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(dto.SalaryCalculationRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	if req.Period == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "period is required", nil)
	}

	jobID, appErr := ctrl.Service.RunSalaryCalculationAsync(ctx, req)
	if appErr != nil {
		if appErr.Code == errors.ErrResourceLocked {
			return c.JSON(http.StatusConflict, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return c.JSON(http.StatusAccepted, echo.Map{
		"job_id": jobID,
		"status": "RUNNING",
	})
}

func (ctrl *PayrollController) GetCalculationJobStatus(c echo.Context) error {
	ctx := c.Request().Context()
	jobID := c.Param("job_id")
	if jobID == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "job_id is required", nil)
	}

	status, appErr := ctrl.Service.GetCalculationJobStatus(ctx, jobID)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, status, "Fetch job status success")
}

func (ctrl *PayrollController) GetSavedPayrollRecords(c echo.Context) error {
	ctx := c.Request().Context()
	period := c.QueryParam("period")
	if period == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "period query parameter is required", nil)
	}

	records, appErr := ctrl.Service.GetSavedPayrollRecords(ctx, period)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, records, "Fetch saved payroll records success")
}
