package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

// ActivityLog maps to the activity_logs table.
// Pointer fields are optional / nullable columns.
type ActivityLog struct {
	ID                      uuid.UUID  `db:"id"`
	UserID                  *uuid.UUID `db:"user_id"`
	UserProfileID           *uuid.UUID `db:"user_profile_id"`
	UserProfileFullName     *string    `db:"user_profile_full_name"`
	UserProfilePositionCode *string    `db:"user_profile_position_code"`
	UserProfilePositionName *string    `db:"user_profile_position_name"`
	UserProfilePositionDesc *string    `db:"user_profile_position_description"`
	Module                  string     `db:"module"`
	Action                  string     `db:"action"` // create | update | delete
	EntityType              string     `db:"entity_type"`
	EntityID                *uuid.UUID `db:"entity_id"`
	Description             *string    `db:"description"`
	FieldsChanged           *string    `db:"fields_changed"`
	OldData                 *[]byte    `db:"old_data"` // JSONB stored as raw bytes
	NewData                 *[]byte    `db:"new_data"`
	HTTPMethod              *string    `db:"http_method"`
	Endpoint                *string    `db:"endpoint"`
	IPAddress               *string    `db:"ip_address"`
	UserAgent               *string    `db:"user_agent"`
	Status                  string     `db:"status"` // success | failed
	ErrorMessage            *string    `db:"error_message"`
	CreatedAt               time.Time  `db:"created_at"`
}

type PaginatedActivityLogEntity = entity.Pagination[ActivityLog]
