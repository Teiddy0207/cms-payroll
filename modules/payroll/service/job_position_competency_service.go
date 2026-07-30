package service

import (
	"cal-salary/core/database"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"cal-salary/modules/payroll/mapper"
	"cal-salary/modules/payroll/repository"
	"context"
	"fmt"
	"math"
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

	var pos *entity.JobPosition
	if profile.PositionID != nil {
		var err error
		pos, err = s.repo.GetJobPositionByID(ctx, *profile.PositionID)
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
	p2Total := p2Score
	isCapped := false
	if pos != nil && pos.P2Cap > 0 && p2Total > pos.P2Cap {
		p2Total = pos.P2Cap
		isCapped = true
	}
	subtotal := p1Total + p2Total

	note := fmt.Sprintf("P1 từ dải lương vị trí và hợp đồng, P2 từ phụ cấp năng lực cá nhân nhân viên %s.", profile.FullName)
	if isCapped {
		note += fmt.Sprintf(" (P2 đã được áp dụng mức trần %v của vị trí %s)", pos.P2Cap, pos.Name)
	}

	breakdown := map[string]any{
		"Tổng điểm P1 (tiêu chuẩn vị trí)": p1Score,
		"Thành tiền P1 (cơ bản)":           p1Total,
		"Tổng phụ cấp P2 (năng lực)":       p2Score,
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
		SystemRate:     0.0,
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
		sortedFormulas, sortErr := sortFormulas(formulas)
		if sortErr != nil {
			s.cache.GetClient().HSet(bgCtx, jobKey, "status", "FAILED", "error", sortErr.Error())
			s.cache.Del(bgCtx, lockKey)
			return
		}

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

	var startDay, endDay time.Time
	if t, errParse := time.Parse("2006-01", period); errParse == nil {
		startDay = t
		endDay = t.AddDate(0, 1, -1)
	} else {
		startDay = time.Now().AddDate(0, 0, -30)
		endDay = time.Now()
	}

	var db *database.Database
	if repoConcrete, ok := s.repo.(*repository.PayrollRepository); ok {
		db = &repoConcrete.DB
	}

	// 1. Query actual work days and OT hours from daily_attendance_sheets
	var actualWorkDays float64
	var otHours float64
	if db != nil {
		queryAtt := `
			SELECT COALESCE(SUM(actual_work_day), 0) as actual, COALESCE(SUM(ot_hours), 0) as ot
			FROM daily_attendance_sheets
			WHERE employee_id = $1 AND date >= $2 AND date <= $3
		`
		var resultAtt struct {
			Actual float64 `db:"actual"`
			Ot     float64 `db:"ot"`
		}
		errQuery := db.SQLx().GetContext(ctx, &resultAtt, queryAtt, profile.ID, startDay, endDay)
		if errQuery == nil {
			actualWorkDays = resultAtt.Actual
			otHours = resultAtt.Ot
		} else {
			logger.Error("calculateAndSaveEmployee: failed to query attendance sheets: %v", errQuery)
		}
	}

	// 2. Query standard work days from system_settings
	standardWorkDays := 22.0
	if db != nil {
		var settingVal string
		errSetting := db.SQLx().GetContext(ctx, &settingVal, `SELECT value FROM system_settings WHERE key = 'standard_work_days'`)
		if errSetting == nil && settingVal != "" {
			var parsed float64
			if _, errScan := fmt.Sscanf(settingVal, "%f", &parsed); errScan == nil && parsed > 0 {
				standardWorkDays = parsed
			}
		}
	}

	// 3. Calculate work ratio
	workRatio := 1.0
	if standardWorkDays > 0 {
		workRatio = actualWorkDays / standardWorkDays
		if workRatio > 1.0 {
			workRatio = 1.0
		}
	}

	// Pro-rate position base P1 and competency P2 by work ratio
	p1Base := preview.P1Total
	p2Base := preview.P2Total

	p1 := p1Base * workRatio
	p2 := p2Base * workRatio
	p3 := 0.0

	gross := p1 + p2 + p3
	tax := gross * 0.1
	net := gross - tax

	env := map[string]interface{}{
		"P1":                 p1,
		"p1":                 p1,
		"P2":                 p2,
		"p2":                 p2,
		"P3":                 p3,
		"p3":                 p3,
		"GROSS_SALARY":       gross,
		"gross_salary":       gross,
		"TAX":                tax,
		"tax":                tax,
		"NET_SALARY":         net,
		"net_salary":         net,
		"ACTUAL_WORK_DAYS":   actualWorkDays,
		"actual_work_days":   actualWorkDays,
		"STANDARD_WORK_DAYS": standardWorkDays,
		"standard_work_days": standardWorkDays,
		"OT_HOURS":           otHours,
		"ot_hours":           otHours,
		"WORK_RATIO":         workRatio,
		"work_ratio":         workRatio,
	}

	hasGrossFormula := false
	hasTaxFormula := false
	hasNetFormula := false
	for _, f := range sortedFormulas {
		if f.VariableName == "GROSS_SALARY" || f.VariableName == "gross_salary" {
			hasGrossFormula = true
		} else if f.VariableName == "TAX" || f.VariableName == "tax" {
			hasTaxFormula = true
		} else if f.VariableName == "NET_SALARY" || f.VariableName == "net_salary" {
			hasNetFormula = true
		}
	}

	details := []entity.PayrollRecordDetail{
		{Component: "P1_BASE", Description: fmt.Sprintf("Lương vị trí gốc (HĐ): %v/tháng", p1Base), Source: "POSITION", Amount: p1Base},
		{Component: "P1", Description: fmt.Sprintf("Lương vị trí thực nhận (%v/%v ngày công)", actualWorkDays, standardWorkDays), Source: "POSITION", Amount: p1},
		{Component: "P2_BASE", Description: fmt.Sprintf("Lương năng lực gốc (Đạt): %v/tháng", p2Base), Source: "PERSONAL", Amount: p2Base},
		{Component: "P2", Description: fmt.Sprintf("Lương năng lực thực nhận (Tỷ lệ công: %v%%)", math.Round(workRatio*10000)/100), Source: "PERSONAL", Amount: p2},
		{Component: "P3", Description: "Lương hiệu quả P3", Source: "FORMULA", Amount: p3},
		{Component: "ACTUAL_WORK_DAYS", Description: "Số ngày công thực tế đi làm", Source: "ATTENDANCE", Amount: actualWorkDays},
		{Component: "STANDARD_WORK_DAYS", Description: "Số ngày công chuẩn của tháng", Source: "ATTENDANCE", Amount: standardWorkDays},
		{Component: "OT_HOURS", Description: "Số giờ làm thêm ngoài giờ (OT)", Source: "ATTENDANCE", Amount: otHours},
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
			} else if f.VariableName == "P3" || f.VariableName == "p3" {
				p3 = val
				env["P3"] = val
				env["p3"] = val
				if !hasGrossFormula {
					gross = p1 + p2 + p3
					env["GROSS_SALARY"] = gross
					env["gross_salary"] = gross
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

	for i := range details {
		if details[i].Component == "P3" {
			details[i].Amount = p3
			break
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

	// 4. Validate P1 salary range and insert anomaly if out of bounds
	if db != nil && profile.PositionID != nil {
		pos, errPos := s.repo.GetJobPositionByID(ctx, *profile.PositionID)
		if errPos == nil && pos != nil && pos.MinSalary > 0 && pos.MaxSalary > 0 {
			if p1Base < pos.MinSalary || p1Base > pos.MaxSalary {
				// Base rate is out of bounds
				anomalyID := uuid.New()
				var deviation float64
				var explanation string
				if p1Base < pos.MinSalary {
					deviation = pos.MinSalary - p1Base
					explanation = fmt.Sprintf("Lương vị trí gốc của hợp đồng (%v) thấp hơn mức tối thiểu dải lương (%v) cho vị trí %s (%s)", p1Base, pos.MinSalary, pos.Name, pos.Code)
				} else {
					deviation = p1Base - pos.MaxSalary
					explanation = fmt.Sprintf("Lương vị trí gốc của hợp đồng (%v) cao hơn mức tối đa dải lương (%v) cho vị trí %s (%s)", p1Base, pos.MaxSalary, pos.Name, pos.Code)
				}

				// Insert anomaly
				queryAnom := `
					INSERT INTO ai_anomalies (id, record_id, metric_flagged, deviation_value, ai_explanation, status, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, 'PENDING', NOW(), NOW())
					ON CONFLICT DO NOTHING
				`
				_, _ = db.SQLx().ExecContext(ctx, queryAnom, anomalyID, record.ID, "P1_RANGE", deviation, explanation)

				// Prepend warning label to P1 detail description
				for i := range details {
					if details[i].Component == "P1" {
						details[i].Description = "[CẢNH BÁO VƯỢT KHUNG] " + details[i].Description
						break
					}
				}
			}
		}
	}

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

func (s *PayrollService) GetSavedPayrollRecords(ctx context.Context, period string, param params.QueryParams) ([]dto.SalaryCalculationItem, *errors.AppError) {
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

	records, err := s.repo.GetPayrollRecords(ctx, periodRecord.ID, param)
	if err != nil {
		logger.Error("GetSavedPayrollRecords:Error %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch payroll records", err)
	}

	items := make([]dto.SalaryCalculationItem, 0, len(records))
	for _, r := range records {
		profile, err := s.repo.GetUserProfileByID(ctx, r.EmployeeID)
		if err != nil {
			continue
		}

		details, err := s.repo.GetPayrollRecordDetails(ctx, r.ID, param)
		if err != nil {
			continue
		}

		item := mapper.ToSalaryCalculationItem(&r, profile, details, period)
		items = append(items, *item)
	}

	return items, nil
}
