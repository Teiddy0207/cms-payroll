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

	err := s.repo.AssignStandardToPosition(ctx, jps)
	if err != nil {
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
	err := s.repo.RemoveStandardFromPosition(ctx, positionID, standardID)
	if err != nil {
		logger.Error("PayrollService:RemoveStandardFromPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to remove standard from position", err)
	}
	return nil
}

// ==========================================
// Job Position Competencies Service
// ==========================================

func (s *PayrollService) AssignCompetencyToPosition(ctx context.Context, positionID uuid.UUID, req *dto.AssignCompetencyRequest) *errors.AppError {
	jpc := &entity.JobPositionCompetency{
		JobDescriptionID: positionID,
		CompetencyID:     req.CompetencyID,
		RequiredLevel:    req.RequiredLevel,
		Weight:           req.Weight,
	}

	err := s.repo.AssignCompetencyToPosition(ctx, jpc)
	if err != nil {
		logger.Error("PayrollService:AssignCompetencyToPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to assign competency to position", err)
	}
	return nil
}

func (s *PayrollService) GetCompetenciesByPosition(ctx context.Context, positionID uuid.UUID) ([]dto.JobPositionCompetencyResponse, *errors.AppError) {
	list, err := s.repo.GetCompetenciesByPosition(ctx, positionID)
	if err != nil {
		logger.Error("PayrollService:GetCompetenciesByPosition:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch competencies by position", err)
	}

	return mapper.ToJobPositionCompetencyDTOList(list), nil
}

func (s *PayrollService) RemoveCompetencyFromPosition(ctx context.Context, positionID uuid.UUID, competencyID uuid.UUID) *errors.AppError {
	err := s.repo.RemoveCompetencyFromPosition(ctx, positionID, competencyID)
	if err != nil {
		logger.Error("PayrollService:RemoveCompetencyFromPosition:Error %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to remove competency from position", err)
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

	err := s.repo.SetSystemSetting(ctx, setting)
	if err != nil {
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

// ==========================================
// Core Calculator Service
// ==========================================

func (s *PayrollService) PreviewSalary(ctx context.Context, employeeID uuid.UUID, period string) (*dto.SalaryPreviewResponse, *errors.AppError) {
	// 1. Get employee profile
	profile, err := s.repo.GetUserProfileByID(ctx, employeeID)
	if err != nil {
		logger.Error("PreviewSalary:GetUserProfileByID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee profile", err)
	}
	if profile == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "employee profile not found", nil)
	}

	if profile.PositionID == nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "employee does not have an assigned job position", nil)
	}

	// 2. Get dynamically calculated P1 Score (SUM of Job Standards for this Position)
	posStds, err := s.repo.GetStandardsByPosition(ctx, *profile.PositionID)
	if err != nil {
		logger.Error("PreviewSalary:GetStandardsByPosition:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job position standards", err)
	}

	var p1Score float64 = 0.0
	var p1StdsBreakdown []dto.JobStandardBreakdown
	for _, item := range posStds {
		p1Score += item.AllowanceValue
		p1StdsBreakdown = append(p1StdsBreakdown, dto.JobStandardBreakdown{
			StandardID:     item.JobStandardID,
			StandardCode:     item.StandardCode,
			StandardName:     item.StandardName,
			AllowanceValue: item.AllowanceValue,
		})
	}

	// 3. Get system_rate (company_point_rate)
	systemRate := 5000.0 // Default rate
	rateSetting, err := s.repo.GetSystemSetting(ctx, "company_point_rate")
	if err == nil && rateSetting != nil {
		parsedRate, errParse := strconv.ParseFloat(rateSetting.Value, 64)
		if errParse == nil {
			systemRate = parsedRate
		}
	}

	// 4. Get required competencies for this job position
	reqComps, err := s.repo.GetCompetenciesByPosition(ctx, *profile.PositionID)
	if err != nil {
		logger.Error("PreviewSalary:GetCompetenciesByPosition:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch position competencies", err)
	}

	// 5. Get actual evaluations for this employee in this period
	qp := params.QueryParams{
		PageNumber: 1,
		PageSize:   1000,
		Filters: map[string]string{
			"employee_id":       employeeID.String(),
			"evaluation_period": period,
		},
	}
	evals, _, err := s.repo.GetEvaluations(ctx, qp)
	if err != nil {
		logger.Error("PreviewSalary:GetEvaluations:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch employee evaluations", err)
	}

	// Map evaluations by competency_id for faster lookup
	evalMap := make(map[uuid.UUID]entity.CompetencyEvaluation)
	for _, ev := range evals {
		evalMap[ev.CompetencyID] = ev
	}

	// 6. Calculate P2 Score by points accumulation (Sum up points of competencies where employee meets/exceeds position requirement)
	var competencyBreakdowns []dto.CompetencyScoreBreakdown
	var p2Score float64 = 0.0

	for _, reqComp := range reqComps {
		achievedLevel := 0
		if ev, exists := evalMap[reqComp.CompetencyID]; exists {
			achievedLevel = ev.Score
		}

		isAchieved := achievedLevel >= reqComp.RequiredLevel
		earnedPoints := 0
		if isAchieved {
			earnedPoints = reqComp.PointValue
			p2Score += float64(reqComp.PointValue)
		}

		competencyBreakdowns = append(competencyBreakdowns, dto.CompetencyScoreBreakdown{
			CompetencyID:   reqComp.CompetencyID,
			CompetencyName: reqComp.CompetencyName,
			CompetencyCode: reqComp.CompetencyCode,
			AchievedLevel:  achievedLevel,
			RequiredLevel:  reqComp.RequiredLevel,
			Weight:         reqComp.Weight,
			PointValue:     reqComp.PointValue,
			IsAchieved:     isAchieved,
			EarnedPoints:   earnedPoints,
		})
	}

	// 7. Compute P1, P2 and Total in currency (VND)
	p1Total := p1Score * systemRate
	p2Total := p2Score * systemRate
	subtotal := p1Total + p2Total

	// 8. Prepare note
	note := fmt.Sprintf("Calculated dynamically based on Job Standards (P1) and Accumulated Competency Points (P2) for employee %s.", profile.FullName)

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
