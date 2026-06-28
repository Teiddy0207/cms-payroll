package mapper

import (
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
)

func ToUserRoleEntity(req *dto.UserRoleRequest) *entity.UserRole {
	if req == nil {
		return nil
	}

	return &entity.UserRole{
		UserID: req.UserID,
		// RoleID:     req.RoleID,
		// AssignedBy: req.AssignedBy,
		// AssignedAt: req.AssignedAt,
		// ExpiresAt:  req.ExpiresAt,
		// IsActive:   req.IsActive,
	}
}
