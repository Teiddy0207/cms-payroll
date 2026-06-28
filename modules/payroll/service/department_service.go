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
// Departments Service Implementation
// ==========================================

func (s *PayrollService) CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, *errors.AppError) {
	dept := mapper.ToDepartmentEntity(req)
	created, err := s.repo.CreateDepartment(ctx, dept)
	if err != nil {
		logger.Error("PayrollService:CreateDepartment:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create department", err)
	}

	return mapper.ToDepartmentDTO(created), nil
}

func (s *PayrollService) GetDepartments(ctx context.Context, qp params.QueryParams) (*dto.PaginatedDepartmentDTO, *errors.AppError) {
	depts, total, err := s.repo.GetDepartments(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetDepartments:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch departments", err)
	}

	return mapper.ToPaginatedDepartmentDTO(depts, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*dto.DepartmentResponse, *errors.AppError) {
	d, err := s.repo.GetDepartmentByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetDepartmentByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch department", err)
	}
	if d == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "department not found", nil)
	}

	return mapper.ToDepartmentDTO(d), nil
}

func (s *PayrollService) UpdateDepartment(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartmentRequest) *errors.AppError {
	dept := &entity.Department{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		ManagerID:   req.ManagerID,
	}

	err := s.repo.UpdateDepartment(ctx, id, dept)
	if err != nil {
		logger.Error("PayrollService:UpdateDepartment:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update department", err)
	}
	return nil
}

func (s *PayrollService) DeleteDepartment(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteDepartment(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteDepartment:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete department", err)
	}
	return nil
}
