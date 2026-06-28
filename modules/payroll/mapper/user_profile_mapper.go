package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToUserProfileEntity(req *dto.CreateUserProfileRequest) *entity.UserProfile {
	return &entity.UserProfile{
		UserID:       req.UserID,
		Code:         req.Code,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Avatar:       req.Avatar,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		PositionID:   req.PositionID,
		DepartmentID: req.DepartmentID,
	}
}

func ToUserProfileDTO(profile *entity.UserProfile) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		ID:           profile.ID,
		UserID:       profile.UserID,
		Code:         profile.Code,
		FullName:     profile.FullName,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		DateOfBirth:  profile.DateOfBirth,
		Gender:       profile.Gender,
		PositionID:   profile.PositionID,
		DepartmentID: profile.DepartmentID,
		CreatedAt:    profile.CreatedAt,
		UpdatedAt:    profile.UpdatedAt,
	}
}

func ToPaginatedUserProfileDTO(items []entity.UserProfile, totalItems, pageNumber, pageSize int) *dto.PaginatedUserProfileDTO {
	responses := make([]dto.UserProfileResponse, len(items))
	for i, u := range items {
		responses[i] = *ToUserProfileDTO(&u)
	}

	return &dto.PaginatedUserProfileDTO{
		Items:      responses,
		TotalItems: totalItems,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
