package dto

import (
	"cal-salary/core/dto"
	authDto "cal-salary/modules/auth/dto"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActivityLogDTO struct {
	ID            uuid.UUID                 `json:"id"`
	UserID        *uuid.UUID                `json:"user_id"`
	UserProfileID *uuid.UUID                `json:"user_profile_id"`
	UserProfile   *authDto.LoginUserProfile `json:"user_profile,omitempty"`
	Position      *authDto.PositionInfo     `json:"position,omitempty"`
	Module        string                    `json:"module"`
	Action        string                    `json:"action"`
	EntityType    string                    `json:"entity_type"`
	EntityID      *uuid.UUID                `json:"entity_id"`
	Description   *string                   `json:"description"`
	FieldsChanged *string                   `json:"fields_changed"`
	OldData       *json.RawMessage          `json:"old_data"`
	NewData       *json.RawMessage          `json:"new_data"`
	HTTPMethod    *string                   `json:"http_method"`
	Endpoint      *string                   `json:"endpoint"`
	IPAddress     *string                   `json:"ip_address"`
	UserAgent     *string                   `json:"user_agent"`
	Status        string                    `json:"status"`
	ErrorMessage  *string                   `json:"error_message"`
	CreatedAt     time.Time                 `json:"created_at"`
}

type PaginatedActivityLogDTO = dto.Pagination[ActivityLogDTO]
