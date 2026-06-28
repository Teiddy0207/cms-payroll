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
// Job Standards Service Implementation
// ==========================================

func (s *PayrollService) CreateJobStandard(ctx context.Context, req *dto.CreateJobStandardRequest) (*dto.JobStandardResponse, *errors.AppError) {
	js := mapper.ToJobStandardEntity(req)
	created, err := s.repo.CreateJobStandard(ctx, js)
	if err != nil {
		logger.Error("PayrollService:CreateJobStandard:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create job standard", err)
	}

	return mapper.ToJobStandardDTO(created), nil
}

func (s *PayrollService) GetJobStandards(ctx context.Context, qp params.QueryParams) (*dto.PaginatedJobStandardDTO, *errors.AppError) {
	standards, total, err := s.repo.GetJobStandards(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetJobStandards:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job standards", err)
	}

	return mapper.ToPaginatedJobStandardDTO(standards, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetJobStandardByID(ctx context.Context, id uuid.UUID) (*dto.JobStandardResponse, *errors.AppError) {
	js, err := s.repo.GetJobStandardByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetJobStandardByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job standard", err)
	}
	if js == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "job standard not found", nil)
	}

	return mapper.ToJobStandardDTO(js), nil
}

func (s *PayrollService) UpdateJobStandard(ctx context.Context, id uuid.UUID, req *dto.UpdateJobStandardRequest) *errors.AppError {
	js := &entity.JobStandard{
		Name:           req.Name,
		AllowanceValue: req.AllowanceValue,
		Description:    req.Description,
	}

	err := s.repo.UpdateJobStandard(ctx, id, js)
	if err != nil {
		logger.Error("PayrollService:UpdateJobStandard:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update job standard", err)
	}
	return nil
}

func (s *PayrollService) DeleteJobStandard(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteJobStandard(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteJobStandard:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete job standard", err)
	}
	return nil
}
