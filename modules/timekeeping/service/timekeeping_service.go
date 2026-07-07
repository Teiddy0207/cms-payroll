package service

import (
	"cal-salary/core/cache"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	payrollEntity "cal-salary/modules/payroll/entity"
	payrollRepo "cal-salary/modules/payroll/repository"
	"cal-salary/modules/timekeeping/dto"
	"cal-salary/modules/timekeeping/entity"
	"cal-salary/modules/timekeeping/repository"
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type TimekeepingServiceImpl struct {
	repo        repository.TimekeepingRepository
	payrollRepo *payrollRepo.PayrollRepository
	cache       *cache.Cache
}

func NewTimekeepingService(repo repository.TimekeepingRepository, payrollRepo *payrollRepo.PayrollRepository, cache *cache.Cache) TimekeepingService {
	return &TimekeepingServiceImpl{
		repo:        repo,
		payrollRepo: payrollRepo,
		cache:       cache,
	}
}

func (s *TimekeepingServiceImpl) getProfileByCode(ctx context.Context, code string) (*payrollEntity.UserProfile, error) {
	var profile payrollEntity.UserProfile
	query := `SELECT id, user_id, code, full_name, phone, avatar, date_of_birth, gender, position_id, department_id, created_at, updated_at FROM user_profiles WHERE code = $1`
	err := s.payrollRepo.DB.SQLx().GetContext(ctx, &profile, query, code)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *TimekeepingServiceImpl) getProfileByUserID(ctx context.Context, userID uuid.UUID) (*payrollEntity.UserProfile, error) {
	var profile payrollEntity.UserProfile
	query := `SELECT id, user_id, code, full_name, phone, avatar, date_of_birth, gender, position_id, department_id, created_at, updated_at FROM user_profiles WHERE user_id = $1`
	err := s.payrollRepo.DB.SQLx().GetContext(ctx, &profile, query, userID)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *TimekeepingServiceImpl) getUserRole(ctx context.Context, userID uuid.UUID) (string, error) {
	var slug string
	query := `
		SELECT r.slug
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	err := s.payrollRepo.DB.SQLx().GetContext(ctx, &slug, query, userID)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (s *TimekeepingServiceImpl) ProcessCheckIn(ctx context.Context, req *dto.CheckinRequest) (*dto.CheckinResponse, *errors.AppError) {
	_, err := s.getProfileByCode(ctx, req.EmployeeCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy mã nhân viên", err)
		}
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi kiểm tra nhân viên", err)
	}

	redisKey := fmt.Sprintf("bluenet_3ps_backend:checkin:last:%s", req.EmployeeCode)
	exists, errExists := s.cache.GetClient().Exists(ctx, redisKey).Result()

	if errExists == nil && exists > 0 {
		return &dto.CheckinResponse{
			Message:   "Check-in thành công (Lọc trùng lặp 5 phút)",
			Timestamp: req.Timestamp,
		}, nil
	}

	log := &entity.AttendanceLog{
		ID:           uuid.New(),
		EmployeeCode: req.EmployeeCode,
		Timestamp:    req.Timestamp,
		LocationGPS:  req.LocationGPS,
		DeviceId:     req.DeviceId,
	}

	if errCreate := s.repo.CreateAttendanceLog(ctx, log); errCreate != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Không thể lưu thông tin check-in", errCreate)
	}

	_ = s.cache.GetClient().Set(ctx, redisKey, "1", 5*time.Minute).Err()

	return &dto.CheckinResponse{
		Message:   "Check-in khuôn mặt thành công!",
		Timestamp: req.Timestamp,
	}, nil
}

func (s *TimekeepingServiceImpl) GetDailyAttendanceSheets(ctx context.Context, userID uuid.UUID, period string) ([]dto.DailyAttendanceResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân viên", err)
	}

	role, errRole := s.getUserRole(ctx, userID)
	if errRole != nil {
		role = "EMPLOYEE"
	}

	parsedTime, errParse := time.Parse("2006-01", period)
	if errParse != nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Định dạng chu kỳ không hợp lệ (YYYY-MM)", errParse)
	}

	start := time.Date(parsedTime.Year(), parsedTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	var employeeIDFilter *uuid.UUID
	var departmentIDFilter *uuid.UUID

	isAdmin := role == "ADMIN" || role == "DIRECTOR"
	isManager := role == "MANAGER"

	if !isAdmin {
		if isManager {
			deptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
			if errDept == nil {
				departmentIDFilter = &deptID
			} else {
				employeeIDFilter = &profile.ID
			}
		} else {
			employeeIDFilter = &profile.ID
		}
	}

	sheets, errSheets := s.repo.GetDailyAttendanceSheets(ctx, employeeIDFilter, start, end)
	if errSheets != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn bảng công", errSheets)
	}

	var allProfiles []payrollEntity.UserProfile
	query := `
		SELECT u.id, u.user_id, u.code, u.full_name, u.phone, u.avatar, u.date_of_birth, u.gender, u.position_id, u.department_id, u.created_at, u.updated_at
		FROM user_profiles u
	`
	errAll := s.payrollRepo.DB.SQLx().SelectContext(ctx, &allProfiles, query)
	if errAll != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi tải thông tin nhân sự", errAll)
	}

	profileMap := make(map[uuid.UUID]payrollEntity.UserProfile)
	for _, p := range allProfiles {
		profileMap[p.ID] = p
	}

	var allDepts []payrollEntity.Department
	queryDept := `SELECT id, name FROM departments`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allDepts, queryDept)
	deptMap := make(map[uuid.UUID]string)
	for _, d := range allDepts {
		deptMap[d.ID] = d.Name
	}

	var responses []dto.DailyAttendanceResponse
	for _, sh := range sheets {
		empProfile, ok := profileMap[sh.EmployeeID]
		if !ok {
			continue
		}

		if departmentIDFilter != nil && (empProfile.DepartmentID == nil || *empProfile.DepartmentID != *departmentIDFilter) {
			continue
		}

		deptName := "Không có"
		if empProfile.DepartmentID != nil {
			if name, ok := deptMap[*empProfile.DepartmentID]; ok {
				deptName = name
			}
		}

		responses = append(responses, dto.DailyAttendanceResponse{
			ID:             sh.ID.String(),
			EmployeeID:     sh.EmployeeID.String(),
			EmployeeCode:   empProfile.Code,
			FullName:       empProfile.FullName,
			DepartmentName: deptName,
			Date:           sh.Date.Format("2006-02-01"),
			CheckIn:        sh.CheckIn,
			CheckOut:       sh.CheckOut,
			ActualWorkDay:  sh.ActualWorkDay,
			OTHours:        sh.OTHours,
			Status:         sh.Status,
		})
	}

	return responses, nil
}

func (s *TimekeepingServiceImpl) CalculateTimesheets(ctx context.Context, req *dto.CalculateTimesheetRequest) *errors.AppError {
	startDay, errStart := time.Parse("2006-01-02", req.StartDate)
	endDay, errEnd := time.Parse("2006-01-02", req.EndDate)
	if errStart != nil || errEnd != nil {
		return errors.NewAppError(errors.ErrInvalidInput, "Định dạng ngày bắt đầu hoặc ngày kết thúc không hợp lệ", nil)
	}

	start := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(endDay.Year(), endDay.Month(), endDay.Day(), 23, 59, 59, 0, time.UTC)

	var profiles []payrollEntity.UserProfile
	query := `SELECT id, code FROM user_profiles`
	errProfiles := s.payrollRepo.DB.SQLx().SelectContext(ctx, &profiles, query)
	if errProfiles != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn hồ sơ nhân viên", errProfiles)
	}

	logs, errLogs := s.repo.GetAttendanceLogsForCalculation(ctx, start, end)
	if errLogs != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn logs chấm công", errLogs)
	}

	var explanations []entity.ExplanationRequest
	queryExp := `SELECT id, employee_id, date, status FROM explanation_requests WHERE date >= $1 AND date <= $2 AND status = 'APPROVED'`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &explanations, queryExp, startDay, endDay)

	var ots []entity.OTRequest
	queryOT := `SELECT id, employee_id, date, hours_requested, is_night_ot, is_holiday_ot, status FROM ot_requests WHERE date >= $1 AND date <= $2 AND status = 'APPROVED'`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &ots, queryOT, startDay, endDay)

	expMap := make(map[string]bool)
	for _, e := range explanations {
		key := fmt.Sprintf("%s:%s", e.EmployeeID.String(), e.Date.Format("2006-01-02"))
		expMap[key] = true
	}

	otMap := make(map[string]entity.OTRequest)
	for _, o := range ots {
		key := fmt.Sprintf("%s:%s", o.EmployeeID.String(), o.Date.Format("2006-01-02"))
		otMap[key] = o
	}

	logsMap := make(map[string][]time.Time)
	for _, l := range logs {
		dateStr := l.Timestamp.Format("2006-01-02")
		key := fmt.Sprintf("%s:%s", l.EmployeeCode, dateStr)
		logsMap[key] = append(logsMap[key], l.Timestamp)
	}

	for _, p := range profiles {
		for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			logKey := fmt.Sprintf("%s:%s", p.Code, dateStr)
			userLogs := logsMap[logKey]

			sheet := &entity.DailyAttendanceSheet{
				EmployeeID: p.ID,
				Date:       d,
				Status:     "ABSENT",
			}
			sheet.ID = uuid.New()

			expKey := fmt.Sprintf("%s:%s", p.ID.String(), dateStr)
			hasApprovedExplanation := expMap[expKey]

			if len(userLogs) == 0 {
				sheet.ActualWorkDay = 0.0
				sheet.OTHours = 0.0
				if hasApprovedExplanation {
					sheet.ActualWorkDay = 1.0
					sheet.Status = "NORMAL"
				}
			} else {
				var checkIn, checkOut time.Time
				checkIn = userLogs[0]
				sheet.CheckIn = &checkIn

				if len(userLogs) >= 2 {
					checkOut = userLogs[len(userLogs)-1]
					sheet.CheckOut = &checkOut
				}

				checkInHour := checkIn.Hour()
				checkInMin := checkIn.Minute()

				inTimeVal := checkInHour*60 + checkInMin
				limitInVal := 9*60 // 09:00 Grace Limit

				isLate := inTimeVal > limitInVal

				if len(userLogs) < 2 {
					sheet.ActualWorkDay = 0.0
					sheet.Status = "LATE"
					if hasApprovedExplanation {
						sheet.ActualWorkDay = 1.0
						sheet.Status = "NORMAL"
					}
				} else {
					checkOutHour := checkOut.Hour()
					checkOutMin := checkOut.Minute()
					outTimeVal := checkOutHour*60 + checkOutMin

					isEarly := outTimeVal < (17*60 + 30) // Before 17:30

					isHalfDay := outTimeVal >= (12*60) && outTimeVal <= (14*60)

					if !isLate && !isEarly {
						sheet.ActualWorkDay = 1.0
						sheet.Status = "NORMAL"
					} else if !isLate && isHalfDay {
						sheet.ActualWorkDay = 0.5
						sheet.Status = "NORMAL"
					} else {
						sheet.ActualWorkDay = 0.0
						if isLate && isEarly {
							sheet.Status = "LATE_EARLY"
						} else if isLate {
							sheet.Status = "LATE"
						} else {
							sheet.Status = "EARLY"
						}
					}

					if hasApprovedExplanation {
						sheet.ActualWorkDay = 1.0
						sheet.Status = "NORMAL"
					}

					if otReq, hasOT := otMap[expKey]; hasOT && outTimeVal > (17*60+30) {
						diffSecs := checkOut.Sub(time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 17, 30, 0, 0, checkOut.Location())).Seconds()
						diffHours := math.Max(0, diffSecs/3600.0)
						sheet.OTHours = math.Min(diffHours, otReq.HoursRequested)
					}
				}
			}

			errUpsert := s.repo.UpsertDailyAttendanceSheet(ctx, sheet)
			if errUpsert != nil {
				logger.Error("UpsertDailyAttendanceSheet failed for %s on %s: %v", p.Code, dateStr, errUpsert)
			}
		}
	}

	return nil
}

func (s *TimekeepingServiceImpl) CreateExplanationRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateExplanationRequest) (*dto.ExplanationRequestResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	dbReq := &entity.ExplanationRequest{
		EmployeeID: profile.ID,
		Date:       req.Date,
		Reason:     req.Reason,
		Status:     "PENDING",
	}
	dbReq.ID = uuid.New()

	if errCreate := s.repo.CreateExplanationRequest(ctx, dbReq); errCreate != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Không thể tạo đơn giải trình", errCreate)
	}

	return &dto.ExplanationRequestResponse{
		ID:           dbReq.ID.String(),
		EmployeeID:   dbReq.EmployeeID.String(),
		EmployeeCode: profile.Code,
		FullName:     profile.FullName,
		Date:         dbReq.Date.Format("2006-01-02"),
		Reason:       dbReq.Reason,
		Status:       dbReq.Status,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *TimekeepingServiceImpl) GetExplanationRequests(ctx context.Context, userID uuid.UUID) ([]dto.ExplanationRequestResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	role, errRole := s.getUserRole(ctx, userID)
	if errRole != nil {
		role = "EMPLOYEE"
	}

	var employeeIDFilter *uuid.UUID
	var departmentIDFilter *uuid.UUID

	isAdmin := role == "ADMIN" || role == "DIRECTOR"
	isManager := role == "MANAGER"

	if !isAdmin {
		if isManager {
			deptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
			if errDept == nil {
				departmentIDFilter = &deptID
			} else {
				employeeIDFilter = &profile.ID
			}
		} else {
			employeeIDFilter = &profile.ID
		}
	}

	reqs, errReqs := s.repo.GetExplanationRequests(ctx, employeeIDFilter, departmentIDFilter)
	if errReqs != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn đơn giải trình", errReqs)
	}

	var allProfiles []payrollEntity.UserProfile
	query := `SELECT id, code, full_name FROM user_profiles`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allProfiles, query)
	profileMap := make(map[uuid.UUID]payrollEntity.UserProfile)
	for _, p := range allProfiles {
		profileMap[p.ID] = p
	}

	var allUsers []struct {
		ID       uuid.UUID `db:"id"`
		Username string    `db:"username"`
	}
	queryUsers := `SELECT id, username FROM users`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allUsers, queryUsers)
	userMap := make(map[uuid.UUID]string)
	for _, u := range allUsers {
		userMap[u.ID] = u.Username
	}

	var responses []dto.ExplanationRequestResponse
	for _, r := range reqs {
		emp, ok := profileMap[r.EmployeeID]
		if !ok {
			continue
		}

		var approvedByStr, approvedNameStr *string
		if r.ApprovedBy != nil {
			str := r.ApprovedBy.String()
			approvedByStr = &str
			if name, ok := userMap[*r.ApprovedBy]; ok {
				approvedNameStr = &name
			}
		}

		responses = append(responses, dto.ExplanationRequestResponse{
			ID:           r.ID.String(),
			EmployeeID:   r.EmployeeID.String(),
			EmployeeCode: emp.Code,
			FullName:     emp.FullName,
			Date:         r.Date.Format("2006-01-02"),
			Reason:       r.Reason,
			ApprovedBy:   approvedByStr,
			ApprovedName: approvedNameStr,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		})
	}

	return responses, nil
}

func (s *TimekeepingServiceImpl) UpdateExplanationRequestStatus(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *dto.UpdateRequestStatus) *errors.AppError {
	isAdminOrDirector, _ := s.payrollRepo.IsAdminOrDirector(ctx, userID)
	managedDeptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
	isManager := errDept == nil

	if !isAdminOrDirector && !isManager {
		return errors.NewAppError(errors.ErrForbidden, "Không có quyền phê duyệt", nil)
	}

	dbReq, err := s.repo.GetExplanationRequestByID(ctx, id)
	if err != nil {
		return errors.NewAppError(errors.ErrNotFound, "Không tìm thấy đơn giải trình", err)
	}

	if !isAdminOrDirector && isManager {
		empDeptID, errEmpDept := s.payrollRepo.GetEmployeeDepartmentID(ctx, dbReq.EmployeeID)
		if errEmpDept != nil || empDeptID != managedDeptID {
			return errors.NewAppError(errors.ErrForbidden, "Nhân viên không thuộc phòng ban quản lý", nil)
		}
	}

	dbReq.Status = req.Status
	dbReq.ApprovedBy = &userID

	if errUpdate := s.repo.UpdateExplanationRequest(ctx, dbReq); errUpdate != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi cập nhật đơn giải trình", errUpdate)
	}

	return nil
}

func (s *TimekeepingServiceImpl) CreateOTRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateOTRequest) (*dto.OTRequestResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	dbReq := &entity.OTRequest{
		EmployeeID:     profile.ID,
		Date:           req.Date,
		HoursRequested: req.HoursRequested,
		IsNightOT:      req.IsNightOT,
		IsHolidayOT:    req.IsHolidayOT,
		Status:         "PENDING",
	}
	dbReq.ID = uuid.New()

	if errCreate := s.repo.CreateOTRequest(ctx, dbReq); errCreate != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Không thể tạo đơn đăng ký OT", errCreate)
	}

	return &dto.OTRequestResponse{
		ID:             dbReq.ID.String(),
		EmployeeID:     dbReq.EmployeeID.String(),
		EmployeeCode:   profile.Code,
		FullName:       profile.FullName,
		Date:           dbReq.Date.Format("2006-01-02"),
		HoursRequested: dbReq.HoursRequested,
		IsNightOT:      dbReq.IsNightOT,
		IsHolidayOT:    dbReq.IsHolidayOT,
		Status:         dbReq.Status,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *TimekeepingServiceImpl) GetOTRequests(ctx context.Context, userID uuid.UUID) ([]dto.OTRequestResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	role, errRole := s.getUserRole(ctx, userID)
	if errRole != nil {
		role = "EMPLOYEE"
	}

	var employeeIDFilter *uuid.UUID
	var departmentIDFilter *uuid.UUID

	isAdmin := role == "ADMIN" || role == "DIRECTOR"
	isManager := role == "MANAGER"

	if !isAdmin {
		if isManager {
			deptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
			if errDept == nil {
				departmentIDFilter = &deptID
			} else {
				employeeIDFilter = &profile.ID
			}
		} else {
			employeeIDFilter = &profile.ID
		}
	}

	reqs, errReqs := s.repo.GetOTRequests(ctx, employeeIDFilter, departmentIDFilter)
	if errReqs != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn đơn đăng ký OT", errReqs)
	}

	var allProfiles []payrollEntity.UserProfile
	query := `SELECT id, code, full_name FROM user_profiles`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allProfiles, query)
	profileMap := make(map[uuid.UUID]payrollEntity.UserProfile)
	for _, p := range allProfiles {
		profileMap[p.ID] = p
	}

	var allUsers []struct {
		ID       uuid.UUID `db:"id"`
		Username string    `db:"username"`
	}
	queryUsers := `SELECT id, username FROM users`
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allUsers, queryUsers)
	userMap := make(map[uuid.UUID]string)
	for _, u := range allUsers {
		userMap[u.ID] = u.Username
	}

	var responses []dto.OTRequestResponse
	for _, r := range reqs {
		emp, ok := profileMap[r.EmployeeID]
		if !ok {
			continue
		}

		var approvedByStr, approvedNameStr *string
		if r.ApprovedBy != nil {
			str := r.ApprovedBy.String()
			approvedByStr = &str
			if name, ok := userMap[*r.ApprovedBy]; ok {
				approvedNameStr = &name
			}
		}

		responses = append(responses, dto.OTRequestResponse{
			ID:             r.ID.String(),
			EmployeeID:     r.EmployeeID.String(),
			EmployeeCode:   emp.Code,
			FullName:       emp.FullName,
			Date:           r.Date.Format("2006-01-02"),
			HoursRequested: r.HoursRequested,
			IsNightOT:      r.IsNightOT,
			IsHolidayOT:    r.IsHolidayOT,
			ApprovedBy:     approvedByStr,
			ApprovedName:   approvedNameStr,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt,
		})
	}

	return responses, nil
}

func (s *TimekeepingServiceImpl) UpdateOTRequestStatus(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *dto.UpdateRequestStatus) *errors.AppError {
	isAdminOrDirector, _ := s.payrollRepo.IsAdminOrDirector(ctx, userID)
	managedDeptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
	isManager := errDept == nil

	if !isAdminOrDirector && !isManager {
		return errors.NewAppError(errors.ErrForbidden, "Không có quyền phê duyệt", nil)
	}

	dbReq, err := s.repo.GetOTRequestByID(ctx, id)
	if err != nil {
		return errors.NewAppError(errors.ErrNotFound, "Không tìm thấy đơn đăng ký OT", err)
	}

	if !isAdminOrDirector && isManager {
		empDeptID, errEmpDept := s.payrollRepo.GetEmployeeDepartmentID(ctx, dbReq.EmployeeID)
		if errEmpDept != nil || empDeptID != managedDeptID {
			return errors.NewAppError(errors.ErrForbidden, "Nhân viên không thuộc phòng ban quản lý", nil)
		}
	}

	dbReq.Status = req.Status
	dbReq.ApprovedBy = &userID

	if errUpdate := s.repo.UpdateOTRequest(ctx, dbReq); errUpdate != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi cập nhật đơn đăng ký OT", errUpdate)
	}

	return nil
}
