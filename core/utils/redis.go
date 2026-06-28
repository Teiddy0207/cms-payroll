package utils

import (
	"cal-salary/core/constants"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateRolePermissionsKey tạo Redis key cho permissions của role
func GenerateRolePermissionsKey(roleSlug string) string {
	return constants.RedisKeyPrefix + "role:" + roleSlug + ":permissions"
}

// GenerateUserPermissionsKey tạo Redis key cho permissions của user
func GenerateUserPermissionsKey(userID string) string {
	return constants.RedisKeyPrefix + "user:" + userID + ":permissions"
}

// GenerateUserRolesKey tạo Redis key cho roles của user
func GenerateUserRolesKey(userID string) string {
	return constants.RedisKeyPrefix + "user:" + userID + ":roles"
}

// GenerateP2TotalFinalPointKey tạo Redis key cho total_final_point
func GenerateP2TotalFinalPointKey(
	kpiTargetID uuid.UUID,
	quarter int,
	year int,
	positionID *uuid.UUID,
	departmentID *uuid.UUID,
	listType string,
) string {
	var entityID string
	if positionID != nil {
		entityID = positionID.String()
	} else if departmentID != nil {
		entityID = departmentID.String()
	} else {
		entityID = "global"
	}

	listTypeLower := strings.ToLower(strings.TrimSpace(listType))
	return fmt.Sprintf("%sp2:total_final_point:%s:%d:%d:%s:%s",
		constants.RedisKeyPrefix,
		kpiTargetID.String(),
		quarter,
		year,
		entityID,
		listTypeLower,
	)
}
