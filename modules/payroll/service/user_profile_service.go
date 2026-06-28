package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"cal-salary/modules/payroll/mapper"
	"context"

	"github.com/google/uuid"
)

// ==========================================
// User Profiles Service Implementation
// ==========================================

func (s *PayrollService) CreateUserProfile(ctx context.Context, req *dto.CreateUserProfileRequest) (*dto.UserProfileResponse, *errors.AppError) {
	profile := mapper.ToUserProfileEntity(req)
	created, err := s.repo.CreateUserProfile(ctx, profile)
	if err != nil {
		logger.Error("PayrollService:CreateUserProfile:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create user profile", err)
	}

	return mapper.ToUserProfileDTO(created), nil
}

func (s *PayrollService) GetUserProfiles(ctx context.Context, qp params.QueryParams) (*dto.PaginatedUserProfileDTO, *errors.AppError) {
	profiles, total, err := s.repo.GetUserProfiles(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetUserProfiles:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch user profiles", err)
	}

	return mapper.ToPaginatedUserProfileDTO(profiles, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetUserProfileByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, *errors.AppError) {
	u, err := s.repo.GetUserProfileByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetUserProfileByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch user profile", err)
	}
	if u == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "user profile not found", nil)
	}

	return mapper.ToUserProfileDTO(u), nil
}

func (s *PayrollService) UpdateUserProfile(ctx context.Context, id uuid.UUID, req *dto.UpdateUserProfileRequest) *errors.AppError {
	profile := &entity.UserProfile{
		FullName:     req.FullName,
		Phone:        req.Phone,
		Avatar:       req.Avatar,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		PositionID:   req.PositionID,
		DepartmentID: req.DepartmentID,
	}

	err := s.repo.UpdateUserProfile(ctx, id, profile)
	if err != nil {
		logger.Error("PayrollService:UpdateUserProfile:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update user profile", err)
	}
	return nil
}

func (s *PayrollService) DeleteUserProfile(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteUserProfile(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteUserProfile:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete user profile", err)
	}
	return nil
}
