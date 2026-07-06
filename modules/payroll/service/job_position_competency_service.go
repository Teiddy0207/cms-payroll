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
	"strconv"

	"github.com/google/uuid"
)

// ==========================================
// Job Position Standards Service
// ==========================================

func (s *PayrollService) AssignStandardToPosition(ctx context.Context, positionID uuid.UUID, req *dto.AssignStandardRequest) *errors.AppError {
	jps := &entity.JobPositionStandard{
		JobDescriptionID: positionID,
		JobStandardID:    req.JobStandardID,
	}
	if err := s.repo.AssignStandardToPosition(ctx, jps); err != nil {
		logger.Error("PayrollService:AssignStandardToPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to assign standard to position", err)
	}
	return nil
}

func (s *PayrollService) GetStandardsByPosition(ctx context.Context, positionID uuid.UUID) ([]dto.JobPositionStandardResponse, *errors.AppError) {
	list, err := s.repo.GetStandardsByPosition(ctx, positionID)
	if err != nil {
		logger.Error("PayrollService:GetStandardsByPosition:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch standards by position", err)
	}
	return mapper.ToJobPositionStandardDTOList(list), nil
}

func (s *PayrollService) RemoveStandardFromPosition(ctx context.Context, positionID uuid.UUID, standardID uuid.UUID) *errors.AppError {
	if err := s.repo.RemoveStandardFromPosition(ctx, positionID, standardID); err != nil {
		logger.Error("PayrollService:RemoveStandardFromPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to remove standard from position", err)
	}
	return nil
}

// ==========================================
// Employee Competencies Service
// ==========================================

func (s *PayrollService) AssignEmployeeCompetencies(ctx context.Context, userProfileID uuid.UUID, req *dto.AssignEmployeeCompetenciesRequest) *errors.AppError {
	if err := s.repo.AssignEmployeeCompetencies(ctx, userProfileID, req.CompetencyIDs); err != nil {
		logger.Error("PayrollService:AssignEmployeeCompetencies:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to assign competencies to employee", err)
	}
	return nil
}

func (s *PayrollService) GetCompetenciesByEmployee(ctx context.Context, userProfileID uuid.UUID) ([]dto.EmployeeCompetencyResponse, *errors.AppError) {
	list, err := s.repo.GetCompetenciesByEmployee(ctx, userProfileID)
	if err != nil {
		logger.Error("PayrollService:GetCompetenciesByEmployee:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch competencies by employee", err)
	}
	return mapper.ToEmployeeCompetencyDTOList(list), nil
}

func (s *PayrollService) RemoveEmployeeCompetency(ctx context.Context, userProfileID uuid.UUID, competencyID uuid.UUID) *errors.AppError {
	if err := s.repo.RemoveEmployeeCompetency(ctx, userProfileID, competencyID); err != nil {
		logger.Error("PayrollService:RemoveEmployeeCompetency:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to remove competency from employee", err)
	}
	return nil
}

// ==========================================
// System Settings Service
// ==========================================

func (s *PayrollService) GetSystemSetting(ctx context.Context, key string) (*dto.SystemSettingResponse, *errors.AppError) {
	setting, err := s.repo.GetSystemSetting(ctx, key)
	if err != nil {
		logger.Error("PayrollService:GetSystemSetting:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch system setting", err)
	}
	if setting == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "system setting not found", nil)
	}
	return mapper.ToSystemSettingDTO(setting), nil
}

func (s *PayrollService) UpdateSystemSetting(ctx context.Context, key string, req *dto.UpdateSystemSettingRequest) *errors.AppError {
	setting := &entity.SystemSetting{
		Key:         key,
		Value:       req.Value,
		Description: req.Description,
	}
	if err := s.repo.SetSystemSetting(ctx, setting); err != nil {
		logger.Error("PayrollService:UpdateSystemSetting:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to update system setting", err)
	}
	return nil
}

func (s *PayrollService) GetAllSystemSettings(ctx context.Context) ([]dto.SystemSettingResponse, *errors.AppError) {
	settings, err := s.repo.GetAllSystemSettings(ctx)
	if err != nil {
		logger.Error("PayrollService:GetAllSystemSettings:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch system settings", err)
	}
	return mapper.ToSystemSettingDTOList(settings), nil
}

func (s *PayrollService) PreviewSalary(ctx context.Context, employeeID uuid.UUID, period string) (*dto.SalaryPreviewResponse, *errors.AppError) {
	// 1. Lấy hồ sơ nhân sự
	profile, err := s.repo.GetUserProfileByID(ctx, employeeID)
	if err != nil {
		logger.Error("PreviewSalary:GetUserProfileByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee profile", err)
	}
	if profile == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "employee profile not found", nil)
	}

	// 2. Tính P1 động từ các tiêu chuẩn công việc của vị trí (nếu có vị trí)
	var p1Score float64
	var p1StdsBreakdown []dto.JobStandardBreakdown
	if profile.PositionID != nil {
		posStds, err := s.repo.GetStandardsByPosition(ctx, *profile.PositionID)
		if err != nil {
			logger.Error("PreviewSalary:GetStandardsByPosition:Error %v", err)
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job position standards", err)
		}
		for _, item := range posStds {
			p1Score += item.AllowanceValue
			p1StdsBreakdown = append(p1StdsBreakdown, dto.JobStandardBreakdown{
				StandardID:     item.JobStandardID,
				StandardCode:   item.StandardCode,
				StandardName:   item.StandardName,
				AllowanceValue: item.AllowanceValue,
			})
		}
	}

	// 3. Lấy đơn giá điểm từ system_settings
	systemRate := 5000.0
	if rateSetting, err := s.repo.GetSystemSetting(ctx, "company_point_rate"); err == nil && rateSetting != nil {
		if parsed, err := strconv.ParseFloat(rateSetting.Value, 64); err == nil {
			systemRate = parsed
		}
	}

	// 4. Tính P2 từ danh sách năng lực đã gán cho nhân sự (employee_competencies)
	empComps, err := s.repo.GetCompetenciesByEmployee(ctx, employeeID)
	if err != nil {
		logger.Error("PreviewSalary:GetCompetenciesByEmployee:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee competencies", err)
	}

	var p2Score float64
	var competencyBreakdowns []dto.CompetencyScoreBreakdown
	for _, comp := range empComps {
		p2Score += float64(comp.PointValue)
		competencyBreakdowns = append(competencyBreakdowns, dto.CompetencyScoreBreakdown{
			CompetencyID:   comp.CompetencyID,
			CompetencyName: comp.CompetencyName,
			CompetencyCode: comp.CompetencyCode,
			PointValue:     comp.PointValue,
			IsAchieved:     true,
			EarnedPoints:   comp.PointValue,
		})
	}

	// 5. Tính tiền
	p1Total := p1Score * systemRate
	p2Total := p2Score * systemRate
	subtotal := p1Total + p2Total

	note := fmt.Sprintf("P1 từ tiêu chuẩn vị trí, P2 từ năng lực cá nhân nhân viên %s.", profile.FullName)

	return &dto.SalaryPreviewResponse{
		EmployeeID:     employeeID,
		FullName:       profile.FullName,
		Period:         period,
		P1Score:        p1Score,
		P1Total:        p1Total,
		P2Score:        p2Score,
		P2Total:        p2Total,
		SystemRate:     systemRate,
		SubtotalP1P2:   subtotal,
		P1Standards:    p1StdsBreakdown,
		P2Competencies: competencyBreakdowns,
		Note:           note,
	}, nil
}

func (s *PayrollService) RunSalaryCalculation(ctx context.Context, req *dto.SalaryCalculationRequest) (*dto.SalaryCalculationResponse, *errors.AppError) {
	if req == nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "request body is required", nil)
	}
	if req.Period == "" {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "period is required", nil)
	}

	// Lấy toàn bộ nhân sự
	profiles, _, err := s.repo.GetUserProfiles(ctx, params.QueryParams{PageNumber: 1, PageSize: 500})
	if err != nil {
		logger.Error("RunSalaryCalculation:GetUserProfiles:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee profiles", err)
	}

	items := make([]dto.SalaryCalculationItem, 0, len(profiles))
	successCount := 0
	failedCount := 0

	for _, profile := range profiles {
		item := dto.SalaryCalculationItem{
			EmployeeID:   profile.ID,
			FullName:     profile.FullName,
			DepartmentID: profile.DepartmentID,
			PositionID:   profile.PositionID,
			Status:       "SUCCESS",
		}

		preview, calcErr := s.PreviewSalary(ctx, profile.ID, req.Period)
		if calcErr != nil {
			item.Status = "FAILED"
			item.Error = calcErr.Message
			failedCount++
		} else {
			item.Preview = preview
			successCount++
		}

		items = append(items, item)
	}

	return &dto.SalaryCalculationResponse{
		Period:           req.Period,
		TotalEmployees:   len(items),
		SuccessEmployees: successCount,
		FailedEmployees:  failedCount,
		Items:            items,
		Note:             "Tính lương theo công thức: P1 (tiêu chuẩn vị trí) + P2 (năng lực cá nhân) × đơn giá điểm.",
	}, nil
}
