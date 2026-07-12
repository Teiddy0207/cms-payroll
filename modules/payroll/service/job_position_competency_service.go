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
	"sync"
	"sync/atomic"
	"time"

	"github.com/expr-lang/expr"
	"github.com/google/uuid"
)



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

	// 2. Tính P1 từ Hợp đồng lao động active và Job Score của vị trí
	var p1Total float64
	var p1Score float64
	var p1StdsBreakdown []dto.JobStandardBreakdown

	activeContract, err := s.repo.GetActiveContractByEmployeeID(ctx, employeeID)
	if err != nil {
		logger.Error("PreviewSalary:GetActiveContractByEmployeeID:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch active contract", err)
	}
	if activeContract != nil {
		p1Total = activeContract.PositionBaseRate
	}

	if profile.PositionID != nil {
		pos, err := s.repo.GetJobPositionByID(ctx, *profile.PositionID)
		if err != nil {
			logger.Error("PreviewSalary:GetJobPositionByID:Error %v", err)
			return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch job position", err)
		}
		if pos != nil {
			p1Score = pos.JobScore
			p1StdsBreakdown = append(p1StdsBreakdown, dto.JobStandardBreakdown{
				StandardID:     pos.ID,
				StandardCode:   pos.Code,
				StandardName:   "Lương dải P1 (" + pos.Name + ")",
				AllowanceValue: p1Total,
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
	p2Total := p2Score * systemRate
	subtotal := p1Total + p2Total

	note := fmt.Sprintf("P1 từ tiêu chuẩn vị trí, P2 từ năng lực cá nhân nhân viên %s.", profile.FullName)

	breakdown := map[string]any{
		"Đơn giá điểm (system rate)":       systemRate,
		"Tổng điểm P1 (tiêu chuẩn vị trí)": p1Score,
		"Thành tiền P1 (cơ bản)":           p1Total,
		"Tổng điểm P2 (năng lực cá nhân)":  p2Score,
		"Thành tiền P2 (năng lực)":         p2Total,
	}

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
		P1:             p1Total,
		P2:             p2Total,
		Total:          subtotal,
		SalaryP1:       p1Total,
		SalaryP2:       p2Total,
		TotalSalary:    subtotal,
		Breakdown:      breakdown,
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
		preview, calcErr := s.PreviewSalary(ctx, profile.ID, req.Period)
		if calcErr != nil {
			item := dto.SalaryCalculationItem{
				EmployeeID:   profile.ID,
				FullName:     profile.FullName,
				DepartmentID: profile.DepartmentID,
				PositionID:   profile.PositionID,
				Status:       "FAILED",
				Error:        calcErr.Message,
			}
			items = append(items, item)
			failedCount++
		} else {
			item := mapper.PreviewToSalaryCalculationItem(preview, &profile)
			items = append(items, *item)
			successCount++
		}
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

func parsePeriodDates(period string) (time.Time, time.Time, int, int, error) {
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, err
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end, int(t.Month()), t.Year(), nil
}

func convertToFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	}
	return 0, false
}

func (s *PayrollService) RunSalaryCalculationAsync(ctx context.Context, req *dto.SalaryCalculationRequest) (string, *errors.AppError) {
	if req == nil || req.Period == "" {
		return "", errors.NewAppError(errors.ErrInvalidInput, "period is required", nil)
	}

	start, end, month, year, err := parsePeriodDates(req.Period)
	if err != nil {
		return "", errors.NewAppError(errors.ErrInvalidInput, "invalid period format", err)
	}

	lockKey := "RunSalaryCalculationAsync:lock:payroll:run:" + req.Period
	ok, err := s.cache.GetClient().SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil || !ok {
		return "", errors.NewAppError(errors.ErrResourceLocked, "a calculation job is already running for this period", nil)
	}

	periodRecord, err := s.repo.GetOrCreatePeriod(ctx, month, year)
	if err != nil {
		s.cache.Del(ctx, lockKey)
		return "", errors.NewAppError(errors.ErrInternalServer, "failed to get or create payroll period", err)
	}

	jobID := uuid.New().String()
	jobKey := "RunSalaryCalculationAsync:job:payroll:" + jobID

	jobData := map[string]interface{}{
		"job_id":          jobID,
		"period":          req.Period,
		"status":          "RUNNING",
		"total_employees": 0,
		"success_count":   0,
		"failed_count":    0,
	}
	err = s.cache.GetClient().HSet(ctx, jobKey, jobData).Err()
	if err != nil {
		s.cache.Del(ctx, lockKey)
		return "", errors.NewAppError(errors.ErrInternalServer, "failed to initialize job status", err)
	}
	s.cache.GetClient().Expire(ctx, jobKey, 24*time.Hour)

	go func() {
		bgCtx := context.Background()
		profiles, _, err := s.repo.GetUserProfiles(bgCtx, params.QueryParams{PageNumber: 1, PageSize: 1000})
		if err != nil {
			s.cache.GetClient().HSet(bgCtx, jobKey, "status", "FAILED", "error", err.Error())
			s.cache.Del(bgCtx, lockKey)
			return
		}

		total := len(profiles)
		s.cache.GetClient().HSet(bgCtx, jobKey, "total_employees", total)

		formulas, _ := s.repo.GetPayrollFormulasByPeriod(bgCtx, start, end)
		sortedFormulas := sortFormulas(formulas)

		numWorkers := 5
		if total < numWorkers {
			numWorkers = total
		}
		if numWorkers == 0 {
			numWorkers = 1
		}

		profileChan := make(chan entity.UserProfile, total)
		for _, p := range profiles {
			profileChan <- p
		}
		close(profileChan)

		var wg sync.WaitGroup
		var successCounter int64
		var failedCounter int64

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for p := range profileChan {
					empLockKey := fmt.Sprintf("RunSalaryCalculationAsync:lock:payroll:process:%s:%s", req.Period, p.ID.String())
					ok, err := s.cache.GetClient().SetNX(bgCtx, empLockKey, "1", 30*time.Second).Result()
					if err != nil || !ok {
						atomic.AddInt64(&failedCounter, 1)
						s.cache.GetClient().HSet(bgCtx, jobKey, "failed_count", atomic.LoadInt64(&failedCounter))
						continue
					}

					calcItemErr := s.calculateAndSaveEmployee(bgCtx, periodRecord.ID, p, req.Period, sortedFormulas)
					s.cache.Del(bgCtx, empLockKey)

					if calcItemErr != nil {
						atomic.AddInt64(&failedCounter, 1)
						s.cache.GetClient().HSet(bgCtx, jobKey, "failed_count", atomic.LoadInt64(&failedCounter))
					} else {
						atomic.AddInt64(&successCounter, 1)
						s.cache.GetClient().HSet(bgCtx, jobKey, "success_count", atomic.LoadInt64(&successCounter))
					}
				}
			}()
		}

		wg.Wait()
		s.cache.GetClient().HSet(bgCtx, jobKey, "status", "SUCCESS")
		s.cache.Del(bgCtx, lockKey)
	}()

	return jobID, nil
}

func (s *PayrollService) calculateAndSaveEmployee(ctx context.Context, periodID uuid.UUID, profile entity.UserProfile, period string, sortedFormulas []entity.PayrollFormula) error {
	preview, err := s.PreviewSalary(ctx, profile.ID, period)
	if err != nil {
		return err
	}

	p1 := preview.P1Total
	p2 := preview.P2Total
	p3 := 0.0

	gross := p1 + p2 + p3
	tax := gross * 0.1
	net := gross - tax

	env := map[string]interface{}{
		"P1":           p1,
		"p1":           p1,
		"P2":           p2,
		"p2":           p2,
		"P3":           p3,
		"p3":           p3,
		"GROSS_SALARY": gross,
		"gross_salary": gross,
		"TAX":          tax,
		"tax":          tax,
		"NET_SALARY":   net,
		"net_salary":   net,
	}

	hasTaxFormula := false
	hasNetFormula := false
	for _, f := range sortedFormulas {
		if f.VariableName == "TAX" || f.VariableName == "tax" {
			hasTaxFormula = true
		} else if f.VariableName == "NET_SALARY" || f.VariableName == "net_salary" {
			hasNetFormula = true
		}
	}

	details := []entity.PayrollRecordDetail{
		{Component: "P1", Description: "Lương vị trí P1", Source: "POSITION", Amount: p1},
		{Component: "P2", Description: "Lương năng lực P2", Source: "PERSONAL", Amount: p2},
		{Component: "P3", Description: "Lương hiệu quả P3", Source: "FORMULA", Amount: p3},
	}

	for _, f := range sortedFormulas {
		program, err := expr.Compile(f.Expression, expr.Env(env))
		if err != nil {
			continue
		}
		output, err := expr.Run(program, env)
		if err != nil {
			continue
		}
		if val, ok := convertToFloat64(output); ok {
			env[f.VariableName] = val
			if f.VariableName == "GROSS_SALARY" || f.VariableName == "gross_salary" {
				gross = val
				env["GROSS_SALARY"] = val
				env["gross_salary"] = val
				if !hasTaxFormula {
					tax = gross * 0.1
					env["TAX"] = tax
					env["tax"] = tax
				}
				if !hasNetFormula {
					net = gross - tax
					env["NET_SALARY"] = net
					env["net_salary"] = net
				}
			} else if f.VariableName == "TAX" || f.VariableName == "tax" {
				tax = val
				env["TAX"] = val
				env["tax"] = val
				if !hasNetFormula {
					net = gross - tax
					env["NET_SALARY"] = net
					env["net_salary"] = net
				}
			} else if f.VariableName == "NET_SALARY" || f.VariableName == "net_salary" {
				net = val
				env["NET_SALARY"] = val
				env["net_salary"] = val
			} else {
				desc := f.Description
				if desc == "" {
					desc = "Hệ số " + f.VariableName
				}
				details = append(details, entity.PayrollRecordDetail{
					Component:   f.VariableName,
					Description: desc,
					Source:      "FORMULA",
					Amount:      val,
				})
			}
		}
	}

	record := &entity.PayrollRecord{
		PeriodID:    periodID,
		EmployeeID:  profile.ID,
		P1Value:     p1,
		P2Value:     p2,
		P3Value:     p3,
		GrossSalary: gross,
		Tax:         tax,
		NetSalary:   net,
		Status:      "DRAFT",
	}

	details = append(details,
		entity.PayrollRecordDetail{Component: "GROSS", Description: "Tổng thu nhập chịu thuế", Source: "FORMULA", Amount: gross},
		entity.PayrollRecordDetail{Component: "TAX", Description: "Thuế thu nhập cá nhân", Source: "FORMULA", Amount: tax},
		entity.PayrollRecordDetail{Component: "NET", Description: "Thực nhận", Source: "FORMULA", Amount: net},
	)

	return s.repo.UpsertPayrollRecord(ctx, record, details)
}

func (s *PayrollService) GetCalculationJobStatus(ctx context.Context, jobID string) (map[string]any, *errors.AppError) {
	jobKey := "RunSalaryCalculationAsync:job:payroll:" + jobID
	data, err := s.cache.GetClient().HGetAll(ctx, jobKey).Result()
	if err != nil {
		logger.Error("GetCalculationJobStatus:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get job status", err)
	}
	if len(data) == 0 {
		return nil, errors.NewAppError(errors.ErrNotFound, "job not found", nil)
	}

	result := make(map[string]any)
	for k, v := range data {
		result[k] = v
	}
	return result, nil
}

func (s *PayrollService) GetSavedPayrollRecords(ctx context.Context, period string) ([]dto.SalaryCalculationItem, *errors.AppError) {
	_, _, month, year, err := parsePeriodDates(period)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "invalid period format", err)
	}

	periodRecord, err := s.repo.GetPayrollPeriodByMonthYear(ctx, month, year)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch payroll period", err)
	}
	if periodRecord == nil {
		return []dto.SalaryCalculationItem{}, nil
	}

	records, err := s.repo.GetPayrollRecords(ctx, periodRecord.ID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch payroll records", err)
	}

	items := make([]dto.SalaryCalculationItem, 0, len(records))
	for _, r := range records {
		profile, err := s.repo.GetUserProfileByID(ctx, r.EmployeeID)
		if err != nil {
			continue
		}

		details, err := s.repo.GetPayrollRecordDetails(ctx, r.ID)
		if err != nil {
			continue
		}

		item := mapper.ToSalaryCalculationItem(&r, profile, details, period)
		items = append(items, *item)
	}

	return items, nil
}

func sortFormulas(formulas []entity.PayrollFormula) []entity.PayrollFormula {
	var other []entity.PayrollFormula
	var gross *entity.PayrollFormula
	var tax *entity.PayrollFormula
	var net *entity.PayrollFormula

	for i := range formulas {
		f := formulas[i]
		switch f.VariableName {
		case "GROSS_SALARY":
			gross = &formulas[i]
		case "TAX":
			tax = &formulas[i]
		case "NET_SALARY":
			net = &formulas[i]
		default:
			other = append(other, f)
		}
	}

	result := make([]entity.PayrollFormula, 0, len(formulas))
	result = append(result, other...)
	if gross != nil {
		result = append(result, *gross)
	}
	if tax != nil {
		result = append(result, *tax)
	}
	if net != nil {
		result = append(result, *net)
	}
	return result
}
