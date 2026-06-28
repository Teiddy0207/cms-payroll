package middleware

import (
	"cal-salary/core/constants"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ActivityLogContext holds rich data that can only be populated by the service/controller layer.
// Services set these fields on echo.Context so ActivityLogMiddleware can pick them up after
// next(c) returns.
type ActivityLogContext struct {
	EntityID      *uuid.UUID
	UserProfileID *uuid.UUID
	Description   *string
	FieldsChanged *string
	OldData       *[]byte // raw JSON (use json.Marshal on your entity before updating)
	NewData       *[]byte // raw JSON (use json.Marshal on your entity after updating)
}

// SetActivityLogContext stores an ActivityLogContext on the echo.Context.
// Call this from your service or controller before (or after) performing the mutation.
//
// Example:
//
//	middleware.SetActivityLogContext(c, &middleware.ActivityLogContext{
//	    EntityID:      &entity.ID,
//	    Description:   utils.StringPtr("updated user profile"),
//	    FieldsChanged: utils.StringPtr("name,email"),
//	    OldData:       oldJSON,
//	    NewData:       newJSON,
//	})
func SetActivityLogContext(c echo.Context, ctx *ActivityLogContext) {
	c.Set(constants.ContextActivityLog, ctx)
}

// getActivityLogContext reads the ActivityLogContext from echo.Context.
// Returns nil if not set.
func getActivityLogContext(c echo.Context) *ActivityLogContext {
	v := c.Get(constants.ContextActivityLog)
	if v == nil {
		return nil
	}
	ctx, _ := v.(*ActivityLogContext)
	return ctx
}
