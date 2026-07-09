package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"cal-salary/modules/payroll/mapper"
	"context"
	"strconv"

	"github.com/google/uuid"
)

// ==========================================
// Job Positions Service Implementation
// ==========================================

func (s *PayrollService) calculateP1SalaryRange(ctx context.Context, pos *entity.JobPosition) {
	kFactor := 4000000.0
	if kSetting, err := s.repo.GetSystemSetting(ctx, "payroll_k_factor"); err == nil && kSetting != nil {
		if parsed, err := strconv.ParseFloat(kSetting.Value, 64); err == nil {
			kFactor = parsed
		}
	}

	pos.JobScore = (pos.EScore * pos.WEWeight) + (pos.CScore * pos.WCWeight) + (pos.RScore * pos.WRWeight)
	pos.Midpoint = pos.JobScore * kFactor

	if pos.SalarySpread <= 0 {
		pos.MinSalary = pos.Midpoint
		pos.MaxSalary = pos.Midpoint
	} else {
		pos.MinSalary = pos.Midpoint / (1.0 + (pos.SalarySpread / 2.0))
		pos.MaxSalary = pos.MinSalary * (1.0 + pos.SalarySpread)
	}
}

func (s *PayrollService) CreateJobPosition(ctx context.Context, req *dto.CreateJobPositionRequest) (*dto.JobPositionResponse, *errors.AppError) {
	pos := mapper.ToJobPositionEntity(req)
	s.calculateP1SalaryRange(ctx, pos)
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
		EScore:       req.EScore,
		CScore:       req.CScore,
		RScore:       req.RScore,
		WEWeight:     req.WEWeight,
		WCWeight:     req.WCWeight,
		WRWeight:     req.WRWeight,
		SalarySpread: req.SalarySpread,
	}
	s.calculateP1SalaryRange(ctx, pos)

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
