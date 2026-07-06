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

	resp, appErr := ctrl.Service.RunSalaryCalculation(ctx, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Salary calculation completed")
}
