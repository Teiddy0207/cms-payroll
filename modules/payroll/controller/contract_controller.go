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
// Contracts Echo Handlers
// ==========================================

func (ctrl *PayrollController) CreateContract(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CreateContractRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	resp, appErr := ctrl.Service.CreateContract(ctx, req)
	if appErr != nil {
		if appErr.Code == errors.ErrAlreadyExists {
			return ctrl.BadRequest(appErr.Code, appErr.Message, nil)
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Create contract success")
}

func (ctrl *PayrollController) GetContracts(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := ctrl.Service.GetContracts(ctx, *qp)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch contracts success")
}

func (ctrl *PayrollController) GetContractByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := ctrl.Service.GetContractByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch contract success")
}

func (ctrl *PayrollController) UpdateContract(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.UpdateContractRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateContract(ctx, id, req)
	if appErr != nil {
		if appErr.Code == errors.ErrAlreadyExists || appErr.Code == errors.ErrNotFound {
			return ctrl.BadRequest(appErr.Code, appErr.Message, nil)
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update contract success")
}

func (ctrl *PayrollController) DeleteContract(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := ctrl.Service.DeleteContract(ctx, id)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Delete contract success")
}
