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



func (s *PayrollService) CreateBatchEvaluations(ctx context.Context, evaluatorID uuid.UUID, req *dto.BatchEvaluationRequest) ([]dto.EvaluationResponse, *errors.AppError) {
	// 1. Kiểm tra xem người chấm có phải Admin hoặc Giám đốc không (bỏ qua RLS check)
	isAdminOrDirector, err := s.repo.IsAdminOrDirector(ctx, evaluatorID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to check user roles", err)
	}

	if !isAdminOrDirector {
		managedDeptID, err := s.repo.GetManagedDepartmentID(ctx, evaluatorID)
		if err != nil {
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to check managed departments", err)
		}
		if managedDeptID == uuid.Nil {
			return nil, errors.NewAppError(errors.ErrForbidden, "Bạn không có quyền chấm điểm (chỉ Trưởng phòng mới được chấm điểm)", nil)
		}

		employeeDeptID, err := s.repo.GetEmployeeDepartmentID(ctx, req.EmployeeID)
		if err != nil {
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to check employee department", err)
		}
		if employeeDeptID == uuid.Nil || employeeDeptID != managedDeptID {
			return nil, errors.NewAppError(errors.ErrForbidden, "Bạn chỉ được phép chấm điểm cho nhân viên thuộc phòng ban do mình quản lý", nil)
		}
	}

	// 3. Xóa các đánh giá cũ trong cùng kỳ của nhân viên này để tránh trùng lặp dữ liệu khi cập nhật
	err = s.repo.DeleteEvaluationsByPeriodAndEmployee(ctx, req.EmployeeID, req.EvaluationPeriod)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to clear old evaluations", err)
	}

	var results []dto.EvaluationResponse
	for _, scoreInput := range req.Scores {
		evalEntity := mapper.ToEvaluationEntity(req, scoreInput, evaluatorID)
		created, err := s.repo.CreateEvaluation(ctx, evalEntity)
		if err != nil {
			logger.Error("PayrollService:CreateBatchEvaluations:Error %v", err)
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to save evaluation score", err)
		}
		results = append(results, *mapper.ToEvaluationDTO(created))
	}

	return results, nil
}

func (s *PayrollService) GetEvaluations(ctx context.Context, qp params.QueryParams) (*dto.PaginatedEvaluationDTO, *errors.AppError) {
	evals, total, err := s.repo.GetEvaluations(ctx, qp)
	if err != nil {
		logger.Error("PayrollService:GetEvaluations:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch evaluations", err)
	}

	return mapper.ToPaginatedEvaluationDTO(evals, total, qp.PageNumber, qp.PageSize), nil
}

func (s *PayrollService) GetEvaluationByID(ctx context.Context, id uuid.UUID) (*dto.EvaluationResponse, *errors.AppError) {
	ce, err := s.repo.GetEvaluationByID(ctx, id)
	if err != nil {
		logger.Error("PayrollService:GetEvaluationByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch evaluation", err)
	}
	if ce == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "evaluation not found", nil)
	}

	return mapper.ToEvaluationDTO(ce), nil
}

func (s *PayrollService) UpdateEvaluation(ctx context.Context, id uuid.UUID, req *dto.CompetencyScoreInput) *errors.AppError {
	ce := &entity.CompetencyEvaluation{
		Score:  req.Score,
		Weight: req.Weight,
	}

	err := s.repo.UpdateEvaluation(ctx, id, ce)
	if err != nil {
		logger.Error("PayrollService:UpdateEvaluation:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update evaluation", err)
	}
	return nil
}

func (s *PayrollService) DeleteEvaluation(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeleteEvaluation(ctx, id)
	if err != nil {
		logger.Error("PayrollService:DeleteEvaluation:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete evaluation", err)
	}
	return nil
}
