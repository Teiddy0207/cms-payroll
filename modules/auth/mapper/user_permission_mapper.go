package mapper

import (
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
)

func ToUserPermissionEntity(req *dto.UserPermissionRequest) *entity.UserPermission {

	if req == nil {
		return nil
	}

	return &entity.UserPermission{
		UserID:       req.UserID,
		PermissionID: req.PermissionID,
		GrantedBy:    req.GrantedBy,
	}
}
