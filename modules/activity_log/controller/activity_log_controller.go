package controller

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"

	"github.com/labstack/echo/v4"
)

// GetActivityLogs handles the request to retrieve a paginated list of activity logs.
func (c *ActivityLogController) GetActivityLogs(ctx echo.Context) error {
	p := params.NewQueryParams(ctx)

	logs, err := c.ActivityLogService.GetActivityLogs(ctx.Request().Context(), *p)
	if err != nil {
		logger.Error("ActivityLogController:GetActivityLogs:Error", err)
		return c.InternalServerError(errors.ErrInternalServer, "Failed to get activity logs", err)
	}

	return c.SuccessResponse(ctx, logs, "Get activity logs successfully")
}
