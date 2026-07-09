package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"cal-salary/modules/payroll/mapper"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// Contracts Service Implementation
// ==========================================

func (s *PayrollService) CreateContract(ctx context.Context, req *dto.CreateContractRequest) (*dto.ContractResponse, *errors.AppError) {
	// Rule: Tại một thời điểm chỉ có tối đa 1 hợp đồng ở trạng thái ACTIVE của mỗi nhân viên.
	if req.Status == "ACTIVE" {
		hasActive, err := s.repo.HasActiveContract(ctx, req.EmployeeID, nil)
		if err != nil {
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to check active contracts", err)
		}
		if hasActive {
			return nil, errors.NewAppError(errors.ErrAlreadyExists, "nhân viên đã có một hợp đồng khác đang hoạt động (ACTIVE)", nil)
		}
	}

	// Validate P1 salary range
	profile, err := s.repo.GetUserProfileByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee profile", err)
	}
	if profile != nil && profile.PositionID != nil {
		pos, err := s.repo.GetJobPositionByID(ctx, *profile.PositionID)
		if err != nil {
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job position info", err)
		}
		if pos != nil && pos.MinSalary > 0 && pos.MaxSalary > 0 {
			if req.PositionBaseRate < pos.MinSalary || req.PositionBaseRate > pos.MaxSalary {
				return nil, errors.NewAppError(errors.ErrInvalidInput, fmt.Sprintf("Mức lương đề xuất (%.0fđ) vượt ngoài khung dải lương P1 của vị trí %s (Min: %.0fđ - Max: %.0fđ)", req.PositionBaseRate, pos.Name, pos.MinSalary, pos.MaxSalary), nil)
			}
		}
	}

	contract := mapper.ToContractEntity(req)
	created, err := s.repo.CreateContract(ctx, contract)
	if err != nil {
		logger.Error("PayrollService:CreateContract:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create contract", err)
	}

	return mapper.ToContractDTO(created), nil
}

func (s *PayrollService) GetContracts(ctx context.Context, qp params.QueryParams) (*dto.PaginatedContractDTO, *errors.AppError) {
	contracts, total, err := s.repo.GetContracts(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetContracts:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch contracts", err)
	}

	return mapper.ToPaginatedContractDTO(contracts, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetContractByID(ctx context.Context, id uuid.UUID) (*dto.ContractResponse, *errors.AppError) {
	c, err := s.repo.GetContractByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetContractByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch contract", err)
	}
	if c == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "contract not found", nil)
	}

	return mapper.ToContractDTO(c), nil
}

func (s *PayrollService) UpdateContract(ctx context.Context, id uuid.UUID, req *dto.UpdateContractRequest) *errors.AppError {
	existing, err := s.repo.GetContractByID(ctx, id)
	if err != nil {
		return errors.NewAppError(errors.ErrInternalServer, "failed to fetch existing contract", err)
	}
	if existing == nil {
		return errors.NewAppError(errors.ErrNotFound, "contract not found", nil)
	}

	// Rule: Tại một thời điểm chỉ có tối đa 1 hợp đồng ở trạng thái ACTIVE của mỗi nhân viên.
	if req.Status == "ACTIVE" {
		hasActive, err := s.repo.HasActiveContract(ctx, existing.EmployeeID, &id)
		if err != nil {
			return errors.NewAppError(errors.ErrInternalServer, "failed to check active contracts", err)
		}
		if hasActive {
			return errors.NewAppError(errors.ErrAlreadyExists, "nhân viên đã có một hợp đồng khác đang hoạt động (ACTIVE)", nil)
		}
	}

	// Validate P1 salary range
	profile, err := s.repo.GetUserProfileByID(ctx, existing.EmployeeID)
	if err == nil && profile != nil && profile.PositionID != nil {
		pos, err := s.repo.GetJobPositionByID(ctx, *profile.PositionID)
		if err == nil && pos != nil && pos.MinSalary > 0 && pos.MaxSalary > 0 {
			if req.PositionBaseRate < pos.MinSalary || req.PositionBaseRate > pos.MaxSalary {
				return errors.NewAppError(errors.ErrInvalidInput, fmt.Sprintf("Mức lương đề xuất (%.0fđ) vượt ngoài khung dải lương P1 của vị trí %s (Min: %.0fđ - Max: %.0fđ)", req.PositionBaseRate, pos.Name, pos.MinSalary, pos.MaxSalary), nil)
			}
		}
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &parsed
	}

	contract := &entity.Contract{
		PositionBaseRate: req.PositionBaseRate,
		StartDate:        startDate,
		EndDate:          endDate,
		Status:           req.Status,
	}

	err = s.repo.UpdateContract(ctx, id, contract)
	if err != nil {
		logger.Error("PayrollService:UpdateContract:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update contract", err)
	}
	return nil
}

func (s *PayrollService) DeleteContract(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteContract(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteContract:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete contract", err)
	}
	return nil
}
