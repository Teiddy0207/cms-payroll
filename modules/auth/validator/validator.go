package validator

import (
	"cal-salary/core/utils"
	"cal-salary/core/validation"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/repository"
	"context"
	"strings"

	"github.com/google/uuid"
)

func ValidateCreateUser(req *dto.CreateUserRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Username) {
		result.AddError("username", "Username is required")
	}

	err := utils.ValidateStrongPassword(req.Password)
	if err != nil {
		result.AddError("password", err.Error())
	}

	if req.Password != req.ConfirmedPassword {
		result.AddError("confirm_password", "Confirm password must be the same as password")
	}

	return result
}

func ValidateVerifyOTPRequest(req *dto.VerifyOTPRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(utils.ToString(req.UserID)) {
		result.AddError("user_id", "User ID is required")
	}

	if utils.IsEmpty(req.OTP) {
		result.AddError("otp", "OTP is required")
	}

	return result
}

func ValidateAssignPermissionToUserRequest(req *dto.UserPermissionRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(utils.ToString(req.UserID)) {
		result.AddError("user_id", "User ID is required")
	}

	if utils.IsEmpty(utils.ToString(req.PermissionID)) {
		result.AddError("permission_id", "Permission ID is required")
	}

	return result
}

func ValidateAssignPermissionToRoleRequest(req *dto.RolePermissionRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(utils.ToString(req.RoleID)) {
		result.AddError("role_id", "Role ID is required")
	}

	if len(req.PermissionID) == 0 {
		result.AddError("permission_id", "Permission ID is required")
	}

	return result
}

func ValidateForgotPasswordRequest(req *dto.ForgotPasswordRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := ValidateIdentifier(req.Identifier)

	return result
}

func ValidateResetPasswordRequest(req *dto.ResetPasswordRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(req.NewPassword) {
		result.AddError("password", "Password is required")
	}

	if utils.IsEmpty(req.NewPassword) {
		result.AddError("new_password", "New password is required")
	}

	if req.NewPassword != req.ConfirmPassword {
		result.AddError("confirm_password", "Confirm password not match")
	}

	return result
}

func ValidateCreateUserRequest(ctx context.Context, req *dto.CreateUserRequest, repo repository.AuthRepositoryInterface) *validation.ValidationResult {
	if req == nil {
		return nil
	}
	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Username) {
		result.AddError("username", "Username is required")
	} else {
		existingUser, err := repo.GetUserByIdentifier(ctx, req.Username)
		if err != nil {
			result.AddError("username", "Không thể kiểm tra username")
		} else if existingUser != nil {
			result.AddError("username", "Username đã tồn tại")
		}
	}

	err := utils.ValidateStrongPassword(req.Password)
	if err != nil {
		result.AddError("password", err.Error())
	}

	if req.Password != req.ConfirmedPassword {
		result.AddError("confirmed_password", "Confirm password must be the same as password")
	}

	return result
}

func ValidateChangePasswordRequest(req *dto.ChangePasswordRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Password) {
		result.AddError("password", "Password is required")
	}

	if utils.IsEmpty(req.NewPassword) {
		result.AddError("new_password", "New password is required")
	}

	if req.Password == req.NewPassword {
		result.AddError("new_password", "New password required not the same as old password")
	}

	if req.NewPassword != req.ConfirmPassword {
		result.AddError("confirm_password", "Confirm password not match")
	}

	return result

}

func ValidateLoginRequest(req *dto.LoginRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := ValidateIdentifier(req.Identifier)

	return result

}

func ValidateRegisterRequest(req *dto.RegisterRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}

	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Phone) {
		result.AddError("phone", "Phone is required")
	}

	if !utils.IsValidPhone(req.Phone) {
		result.AddError("phone", "Phone is invalid")
	}

	if err := utils.ValidateStrongPassword(req.Password); err != nil {
		result.AddError("password", err.Error())
	}

	return result
}

func ValidateRoleRequest(ctx context.Context, repo repository.AuthRepositoryInterface, req *dto.RoleRequest, excludeID *uuid.UUID) *validation.ValidationResult {
	if req == nil {
		return nil
	}
	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Name) {
		result.AddError("name", "Name is required")
	}

	if utils.IsEmpty(req.Slug) {
		result.AddError("slug", "Slug is required")
	} else {
		existingRole, err := repo.GetRoleByNameOrSlug(ctx, req.Name, req.Slug)
		if err != nil {
			result.AddError("name", "Không thể kiểm tra role")
			result.AddError("slug", "Không thể kiểm tra role")
		} else if existingRole != nil {
			if excludeID != nil && existingRole.ID == *excludeID {
			} else {
				if existingRole.Name == req.Name {
					result.AddError("name", "Role với tên này đã tồn tại")
				}
				if existingRole.Slug == req.Slug {
					result.AddError("slug", "Role với slug này đã tồn tại")
				}
			}
		}
	}

	return result
}

func ValidateRoleRequestUpdate(ctx context.Context, repo repository.AuthRepositoryInterface, req *dto.RoleRequest, excludeID *uuid.UUID) *validation.ValidationResult {
	if req == nil {
		return nil
	}
	result := validation.NewValidationResult()

	existingRole, err := repo.GetRoleByNameOrSlug(ctx, req.Name, req.Slug)
	if err != nil {
		result.AddError("name", "Không thể kiểm tra role")
		result.AddError("slug", "Không thể kiểm tra role")
	} else if existingRole != nil {
		if excludeID != nil && existingRole.ID == *excludeID {
		} else {
			if existingRole.Name == req.Name {
				result.AddError("name", "Role với tên này đã tồn tại")
			}
			if existingRole.Slug == req.Slug {
				result.AddError("slug", "Role với slug này đã tồn tại")
			}
		}
	}

	return result
}

func ValidatePermissionRequest(req *dto.PermissionRequest) *validation.ValidationResult {
	if req == nil {
		return nil
	}
	result := validation.NewValidationResult()

	if utils.IsEmpty(req.Name) {
		result.AddError("name", "Name is required")
	}

	if utils.IsEmpty(req.Resource) {
		result.AddError("resource", "Resource is required")
	}

	if utils.IsEmpty(req.Action) {
		result.AddError("action", "Action is required")
	}

	return result
}

// ValidateIdentifier validates a string as either phone number or email
// Returns validation result with appropriate error messages
func ValidateIdentifier(identifier string) *validation.ValidationResult {
	result := validation.NewValidationResult()

	if utils.IsEmpty(identifier) {
		result.AddError("identifier", "Identifier is required")
		return result
	}

	// Check if identifier contains @ symbol to determine if it's an email
	if strings.Contains(identifier, "@") {
		// Validate as email
		if !utils.IsValidEmail(identifier) && !utils.IsValidEmailDomain(identifier) {
			result.AddError("identifier", "Invalid email format")
		}
	}

	return result
}
