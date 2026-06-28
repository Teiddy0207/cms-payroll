package mapper

import (
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
	"encoding/json"
)

func ToUserEntity(user *dto.CreateUserRequest) *entity.User {
	if user == nil {
		return nil
	}
	return &entity.User{
		Username: &user.Username,
		Password: user.Password,
	}
}

func ToUserDTO(user *entity.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Password:        user.Password,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func ToUserPaginationDTO(entity *entity.PaginatedUserEntity) *dto.PaginatedUserDTO {

	if entity == nil {
		return nil
	}

	userResponses := make([]dto.UserResponse, len(entity.Items))
	for i, user := range entity.Items {
		userResponses[i] = *ToUserDTO(&user)
	}

	// Tính total pages
	totalPages := 0
	if entity.PageSize > 0 {
		totalPages = (entity.TotalItems + entity.PageSize - 1) / entity.PageSize
	}

	return &dto.PaginatedUserDTO{
		Items:      userResponses,
		TotalItems: entity.TotalItems,
		TotalPages: totalPages,
		PageNumber: entity.PageNumber,
		PageSize:   entity.PageSize,
	}
}

func ToUserWithProfileDTO(userWithProfile *entity.UserWithProfile) *dto.UserResponse {
	if userWithProfile == nil {
		return nil
	}

	userResponse := ToUserDTO(&userWithProfile.User)

	// Thêm user_profile nếu có
	if userWithProfile.ProfileID != nil {
		userResponse.UserProfile = &dto.UserProfileInfo{
			ID:       userWithProfile.ProfileID,
			Code:     userWithProfile.ProfileCode,
			FullName: userWithProfile.ProfileFullName,
			Phone:    userWithProfile.ProfilePhone,
			Gender:   userWithProfile.ProfileGender,
			Avatar:   userWithProfile.ProfileAvatar,
		}
	}

	// Thêm official_id và official nếu có
	if userWithProfile.OfficialID != nil {
		userResponse.OfficialID = userWithProfile.OfficialID
		if userWithProfile.OfficeCode != nil || userWithProfile.OfficeName != nil {
			userResponse.Official = &dto.OfficeInfo{
				Code: userWithProfile.OfficeCode,
				Name: userWithProfile.OfficeName,
			}
		}
	}

	// Thêm role nếu có
	if userWithProfile.RoleID != nil && userWithProfile.RoleCode != nil && userWithProfile.RoleName != nil {
		userResponse.Roles = &dto.RoleInfo{
			ID:          *userWithProfile.RoleID,
			Code:        *userWithProfile.RoleCode,
			Name:        *userWithProfile.RoleName,
			Description: userWithProfile.RoleDesc,
		}
	}

	return userResponse
}

func ToUserWithProfilePaginationDTO(entity *entity.PaginatedUserWithProfileEntity) *dto.PaginatedUserDTO {
	if entity == nil {
		return nil
	}

	userResponses := make([]dto.UserResponse, len(entity.Items))
	for i, userWithProfile := range entity.Items {
		userResponses[i] = *ToUserWithProfileDTO(&userWithProfile)
	}

	totalPages := 0
	if entity.PageSize > 0 {
		totalPages = (entity.TotalItems + entity.PageSize - 1) / entity.PageSize
	}

	return &dto.PaginatedUserDTO{
		Items:      userResponses,
		TotalItems: entity.TotalItems,
		TotalPages: totalPages,
		PageNumber: entity.PageNumber,
		PageSize:   entity.PageSize,
	}
}

func ToUserDetailDTO(user *entity.UserDetail) *dto.UserDetailDTO {
	if user == nil {
		return nil
	}

	userDTO := &dto.UserDetailDTO{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		DisplayName: user.DisplayName,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
	}

	// Thêm role nếu có
	if user.RoleID != nil && user.RoleCode != nil && user.RoleName != nil {
		userDTO.Roles = &dto.RoleInfo{
			ID:          *user.RoleID,
			Code:        *user.RoleCode,
			Name:        *user.RoleName,
			Description: user.RoleDesc,
		}
	}

	// Parse permissions từ JSON string
	if user.Permissions != nil && *user.Permissions != "" && *user.Permissions != "[]" {
		var permissions []dto.PermissionResponse
		if err := json.Unmarshal([]byte(*user.Permissions), &permissions); err == nil {
			userDTO.Permissions = permissions
		}
	}

	return userDTO
}

func ToUserUpdateEntity(req *dto.UserUpdateRequest, existing *entity.User) *entity.User {
	if existing == nil {
		return nil
	}

	updated := &entity.User{
		Email:           existing.Email,
		Username:        existing.Username,
		Password:        existing.Password,
		EmailVerifiedAt: existing.EmailVerifiedAt,
		PositionID:      existing.PositionID,
		IsActive:        existing.IsActive,
		BaseEntity:      existing.BaseEntity,
	}

	if req == nil {
		return updated
	}

	if req.Email != nil {
		updated.Email = req.Email
	}
	if req.Username != nil {
		updated.Username = req.Username
	}
	if req.Password != nil && *req.Password != "" {
		updated.Password = *req.Password
	}
	if req.PositionID != nil {
		updated.PositionID = req.PositionID
	}
	if req.IsActive != nil {
		updated.IsActive = *req.IsActive
	}

	return updated
}

func ToUsersWithProfileDTO(entity *entity.PaginatedUserWithProfileEntity) *dto.PaginatedUserDTO {
	if entity == nil {
		return nil
	}

	userResponses := make([]dto.UserResponse, len(entity.Items))
	for i, userWithProfile := range entity.Items {
		userResponses[i] = *ToUserWithProfileDTO(&userWithProfile)
	}
	return &dto.PaginatedUserDTO{
		Items: userResponses,
	}
}

func ToUserWithProfileDTOs(users []*entity.UserWithProfile) []dto.UserResponse {
	if users == nil {
		return []dto.UserResponse{}
	}

	result := make([]dto.UserResponse, 0, len(users))

	for _, u := range users {
		if u == nil {
			continue
		}

		dtoPtr := ToUserWithProfileDTO(u)
		if dtoPtr != nil {
			result = append(result, *dtoPtr)
		}
	}

	return result
}
