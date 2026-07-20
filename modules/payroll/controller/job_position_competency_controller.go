package controller

import (
	"cal-salary/core/errors"
	"cal-salary/modules/payroll/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)



func (ctrl *PayrollController) AssignEmployeeCompetencies(c echo.Context) error {
	ctx := c.Request().Context()
	userProfileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid User Profile ID format", nil)
	}

	req := new(dto.AssignEmployeeCompetenciesRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.AssignEmployeeCompetencies(ctx, userProfileID, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Assign competencies to employee success")
}

func (ctrl *PayrollController) GetCompetenciesByEmployee(c echo.Context) error {
	ctx := c.Request().Context()
	userProfileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid User Profile ID format", nil)
	}

	resp, appErr := ctrl.Service.GetCompetenciesByEmployee(ctx, userProfileID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch competencies by employee success")
}

func (ctrl *PayrollController) RemoveEmployeeCompetency(c echo.Context) error {
	ctx := c.Request().Context()
	userProfileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid User Profile ID format", nil)
	}

	competencyID, err := uuid.Parse(c.Param("competency_id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid Competency ID format", nil)
	}

	appErr := ctrl.Service.RemoveEmployeeCompetency(ctx, userProfileID, competencyID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Remove competency from employee success")
}
