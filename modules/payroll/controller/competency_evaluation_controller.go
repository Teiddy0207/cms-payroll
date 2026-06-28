package controller

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/core/utils"
	"cal-salary/modules/payroll/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ==========================================
// Competency Evaluations Echo Handlers
// ==========================================

func (ctrl *PayrollController) CreateBatchEvaluations(c echo.Context) error {
	ctx := c.Request().Context()

	// Extract evaluator ID from context token claims (set by AuthMiddleware)
	userData := c.Get(constants.ContextTokenData)
	if userData == nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Unauthorized", nil)
	}
	claims, ok := userData.(*utils.TokenClaims)
	if !ok {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Unauthorized", nil)
	}
	evaluatorID := claims.UserID

	req := new(dto.BatchEvaluationRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	resp, appErr := ctrl.Service.CreateBatchEvaluations(ctx, evaluatorID, req)
	if appErr != nil {
		if appErr.Code == errors.ErrForbidden {
			return ctrl.Forbidden(appErr.Code, appErr.Message, nil)
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}

	return ctrl.SuccessResponse(c, resp, "Submit competency evaluations success")
}

func (ctrl *PayrollController) GetEvaluations(c echo.Context) error {
	ctx := c.Request().Context()
	qp := params.NewQueryParams(c)
	resp, appErr := ctrl.Service.GetEvaluations(ctx, *qp)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch evaluations success")
}

func (ctrl *PayrollController) GetEvaluationByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	resp, appErr := ctrl.Service.GetEvaluationByID(ctx, id)
	if appErr != nil {
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Fetch evaluation success")
}

func (ctrl *PayrollController) UpdateEvaluation(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	req := new(dto.CompetencyScoreInput)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid request data", nil)
	}

	appErr := ctrl.Service.UpdateEvaluation(ctx, id, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Update evaluation success")
}

func (ctrl *PayrollController) DeleteEvaluation(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Invalid ID format", nil)
	}

	appErr := ctrl.Service.DeleteEvaluation(ctx, id)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Delete evaluation success")
}
