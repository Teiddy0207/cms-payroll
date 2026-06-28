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
// Job Positions Service Implementation
// ==========================================

func (s *PayrollService) CreateJobPosition(ctx context.Context, req *dto.CreateJobPositionRequest) (*dto.JobPositionResponse, *errors.AppError) {
	pos := mapper.ToJobPositionEntity(req)
	created, err := s.repo.CreateJobPosition(ctx, pos)
	if err != nil {
		logger.Error("PayrollService:CreateJobPosition:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create job position", err)
	}

	return mapper.ToJobPositionDTO(created), nil
}

func (s *PayrollService) GetJobPositions(ctx context.Context, qp params.QueryParams) (*dto.PaginatedJobPositionDTO, *errors.AppError) {
	poses, total, err := s.repo.GetJobPositions(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetJobPositions:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job positions", err)
	}

	return mapper.ToPaginatedJobPositionDTO(poses, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetJobPositionByID(ctx context.Context, id uuid.UUID) (*dto.JobPositionResponse, *errors.AppError) {
	p, err := s.repo.GetJobPositionByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetJobPositionByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job position", err)
	}
	if p == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "job position not found", nil)
	}

	return mapper.ToJobPositionDTO(p), nil
}

func (s *PayrollService) UpdateJobPosition(ctx context.Context, id uuid.UUID, req *dto.UpdateJobPositionRequest) *errors.AppError {
	pos := &entity.JobPosition{
		Name:         req.Name,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
	}

	err := s.repo.UpdateJobPosition(ctx, id, pos)
	if err != nil {
		logger.Error("PayrollService:UpdateJobPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update job position", err)
	}
	return nil
}

func (s *PayrollService) DeleteJobPosition(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteJobPosition(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteJobPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete job position", err)
	}
	return nil
}
