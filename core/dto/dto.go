package dto

import (
	"errors"

	"cal-salary/core/constants"
	"cal-salary/core/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Pagination[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

// BaseRequest provides common fields for all API requests
// It automatically extracts and stores the authorization token
type BaseRequest struct {
	// Token is automatically extracted from Authorization header
	// Format: "Bearer <token>"
	Token string `json:"-" header:"Authorization"`

	// UserID is extracted from the JWT token after validation
	UserID uuid.UUID `json:"-"`
}

// LoadFromContext loads token and user info from echo context (set by middleware)
func (br *BaseRequest) LoadFromContext(c echo.Context) error {
	// Try to get token claims from context first (set by AuthMiddleware)
	tokenData := c.Get(constants.ContextTokenData)
	if tokenData != nil {
		if claims, ok := tokenData.(*utils.TokenClaims); ok {
			// Get clean token from header for storage
			token, err := utils.GetTokenFromHeader(c)
			if err != nil {
				return err
			}
			br.Token = token
			br.UserID = claims.UserID
			return nil
		}
	}

	// Fallback: parse token directly if not in context
	claims := utils.ParseDataFromToken(c)
	if claims == nil {
		return errors.New("invalid or missing token")
	}

	token, err := utils.GetTokenFromHeader(c)
	if err != nil {
		return err
	}

	br.Token = token
	br.UserID = claims.UserID
	return nil
}

// Bind binds JSON request and automatically extracts token and user info
// Usage: var req MyRequest; err := dto.Bind(c, &req)
func Bind(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	if baseReq, ok := req.(interface{ LoadFromContext(echo.Context) error }); ok {
		return baseReq.LoadFromContext(c)
	}

	return nil
}
