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
// Competency Dictionaries Service Implementation
// ==========================================

func (s *PayrollService) CreateCompetency(ctx context.Context, req *dto.CreateCompetencyRequest) (*dto.CompetencyResponse, *errors.AppError) {
	c := mapper.ToCompetencyEntity(req)
	created, err := s.repo.CreateCompetency(ctx, c)
	if err != nil {
		logger.Error("PayrollService:CreateCompetency:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create competency", err)
	}

	return mapper.ToCompetencyDTO(created), nil
}

func (s *PayrollService) GetCompetencies(ctx context.Context, qp params.QueryParams) (*dto.PaginatedCompetencyDTO, *errors.AppError) {
	list, total, err := s.repo.GetCompetencies(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetCompetencies:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch competencies", err)
	}

	return mapper.ToPaginatedCompetencyDTO(list, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetCompetencyByID(ctx context.Context, id uuid.UUID) (*dto.CompetencyResponse, *errors.AppError) {
	c, err := s.repo.GetCompetencyByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetCompetencyByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch competency", err)
	}
	if c == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "competency not found", nil)
	}

	return mapper.ToCompetencyDTO(c), nil
}

func (s *PayrollService) UpdateCompetency(ctx context.Context, id uuid.UUID, req *dto.UpdateCompetencyRequest) *errors.AppError {
	c := &entity.Competency{
		Name:        req.Name,
		PointValue:  req.PointValue,
		Description: req.Description,
		TypeID:      req.TypeID,
	}

	err := s.repo.UpdateCompetency(ctx, id, c)
	if err != nil {
		logger.Error("PayrollService:UpdateCompetency:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update competency", err)
	}
	return nil
}

func (s *PayrollService) DeleteCompetency(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteCompetency(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteCompetency:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete competency", err)
	}
	return nil
}
