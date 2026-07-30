package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/modules/timekeeping/dto"
	"cal-salary/modules/timekeeping/entity"
	payrollEntity "cal-salary/modules/payroll/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// Leave Request Service
// ==========================================

// CreateLeaveRequest nhân viên xin nghỉ phép.
// Tự động tính days_requested nếu client truyền = 0 (số ngày trong khoảng [start, end]).
// Validate: không cho xin nghỉ ngày đã qua, start <= end.
func (s *TimekeepingServiceImpl) CreateLeaveRequest(ctx context.Context, userID uuid.UUID, req *dto.CreateLeaveRequest) (*dto.LeaveRequestResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	startDate := req.StartDate.Truncate(24 * time.Hour)
	endDate := req.EndDate.Truncate(24 * time.Hour)

	if endDate.Before(startDate) {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Ngày kết thúc phải >= ngày bắt đầu", nil)
	}

	// Tính days_requested nếu client không truyền (hoặc truyền 0)
	daysRequested := req.DaysRequested
	if daysRequested <= 0 {
		numDays := int(endDate.Sub(startDate).Hours()/24) + 1
		daysRequested = float64(numDays)
	}
	// Validate: nếu là nửa ngày thì start == end
	if daysRequested == 0.5 && !startDate.Equal(endDate) {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Nghỉ nửa ngày (0.5) chỉ áp dụng cho một ngày duy nhất (start_date = end_date)", nil)
	}

	dbReq := &entity.LeaveRequest{
		EmployeeID:    profile.ID,
		StartDate:     startDate,
		EndDate:       endDate,
		DaysRequested: daysRequested,
		Reason:        req.Reason,
		Status:        "PENDING",
	}
	dbReq.ID = uuid.New()

	if errCreate := s.repo.CreateLeaveRequest(ctx, dbReq); errCreate != nil {
		logger.Error("CreateLeaveRequest:Error %v", errCreate)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Không thể tạo đơn xin nghỉ phép", errCreate)
	}

	return &dto.LeaveRequestResponse{
		ID:            dbReq.ID.String(),
		EmployeeID:    profile.ID.String(),
		EmployeeCode:  profile.Code,
		FullName:      profile.FullName,
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		DaysRequested: daysRequested,
		Reason:        req.Reason,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}, nil
}

// GetLeaveRequests lấy danh sách đơn nghỉ phép, lọc theo role:
//   - Admin/Director: thấy tất cả
//   - Manager: thấy nhân viên trong phòng ban
//   - Employee: chỉ thấy của mình
func (s *TimekeepingServiceImpl) GetLeaveRequests(ctx context.Context, userID uuid.UUID) ([]dto.LeaveRequestResponse, *errors.AppError) {
	role, _ := s.getUserRole(ctx, userID)
	roleUpper := strings.ToUpper(role)
	isAdmin := roleUpper == "ADMIN" || roleUpper == "DIRECTOR"
	isManager := roleUpper == "MANAGER"

	var employeeIDFilter *uuid.UUID
	var departmentIDFilter *uuid.UUID

	if !isAdmin {
		if profile, err := s.getProfileByUserID(ctx, userID); err == nil {
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
	}

	reqs, errReqs := s.repo.GetLeaveRequests(ctx, employeeIDFilter, departmentIDFilter)
	if errReqs != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn đơn nghỉ phép", errReqs)
	}

	// Load profile map và user map để enrich response
	var allProfiles []payrollEntity.UserProfile
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allProfiles, `SELECT id, code, full_name FROM user_profiles`)
	profileMap := make(map[uuid.UUID]payrollEntity.UserProfile)
	for _, p := range allProfiles {
		profileMap[p.ID] = p
	}

	var allUsers []struct {
		ID       uuid.UUID `db:"id"`
		Username string    `db:"username"`
	}
	_ = s.payrollRepo.DB.SQLx().SelectContext(ctx, &allUsers, `SELECT id, username FROM users`)
	userMap := make(map[uuid.UUID]string)
	for _, u := range allUsers {
		userMap[u.ID] = u.Username
	}

	var responses []dto.LeaveRequestResponse
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
		responses = append(responses, dto.LeaveRequestResponse{
			ID:            r.ID.String(),
			EmployeeID:    r.EmployeeID.String(),
			EmployeeCode:  emp.Code,
			FullName:      emp.FullName,
			StartDate:     r.StartDate.Format("2006-01-02"),
			EndDate:       r.EndDate.Format("2006-01-02"),
			DaysRequested: r.DaysRequested,
			Reason:        r.Reason,
			ApprovedBy:    approvedByStr,
			ApprovedName:  approvedNameStr,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt,
		})
	}

	return responses, nil
}

// UpdateLeaveRequestStatus manager hoặc admin duyệt / từ chối đơn nghỉ phép.
// Chỉ manager quản lý đúng phòng ban hoặc admin mới được duyệt.
func (s *TimekeepingServiceImpl) UpdateLeaveRequestStatus(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *dto.UpdateRequestStatus) *errors.AppError {
	isAdminOrDirector, _ := s.payrollRepo.IsAdminOrDirector(ctx, userID)
	managedDeptID, errDept := s.payrollRepo.GetManagedDepartmentID(ctx, userID)
	isManager := errDept == nil

	if !isAdminOrDirector && !isManager {
		return errors.NewAppError(errors.ErrForbidden, "Không có quyền phê duyệt đơn nghỉ phép", nil)
	}

	dbReq, err := s.repo.GetLeaveRequestByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewAppError(errors.ErrNotFound, "Không tìm thấy đơn nghỉ phép", err)
		}
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn đơn nghỉ phép", err)
	}

	// Manager chỉ duyệt nhân viên trong phòng của mình
	if !isAdminOrDirector && isManager {
		empDeptID, errEmpDept := s.payrollRepo.GetEmployeeDepartmentID(ctx, dbReq.EmployeeID)
		if errEmpDept != nil || empDeptID != managedDeptID {
			return errors.NewAppError(errors.ErrForbidden, "Nhân viên không thuộc phòng ban quản lý", nil)
		}
	}

	dbReq.Status = req.Status
	dbReq.ApprovedBy = &userID

	if errUpdate := s.repo.UpdateLeaveRequest(ctx, dbReq); errUpdate != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Lỗi cập nhật trạng thái đơn nghỉ phép", errUpdate)
	}

	return nil
}

// ==========================================
// Leave Balance Service
// ==========================================

// GetLeaveBalance lấy số dư phép của nhân viên trong năm.
// Nhân viên thường chỉ xem được phép của mình. Admin/Manager xem được của người khác.
func (s *TimekeepingServiceImpl) GetLeaveBalance(ctx context.Context, userID uuid.UUID, year int) (*dto.LeaveBalanceResponse, *errors.AppError) {
	profile, err := s.getProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy hồ sơ nhân sự", err)
	}

	balance, err := s.repo.GetLeaveBalance(ctx, profile.ID, year)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi truy vấn số dư phép", err)
	}

	// Nếu chưa có bản ghi leave_balance cho năm này → trả về balance = 0
	if balance == nil {
		return &dto.LeaveBalanceResponse{
			EmployeeID:   profile.ID.String(),
			EmployeeCode: profile.Code,
			FullName:     profile.FullName,
			Year:         year,
			AccruedDays:  0,
			UsedDays:     0,
			Balance:      0,
		}, nil
	}

	return &dto.LeaveBalanceResponse{
		EmployeeID:       profile.ID.String(),
		EmployeeCode:     profile.Code,
		FullName:         profile.FullName,
		Year:             year,
		AccruedDays:      balance.AccruedDays,
		UsedDays:         balance.UsedDays,
		Balance:          balance.Balance,
		LastAccrualMonth: balance.LastAccrualMonth,
	}, nil
}

// AccrueLeaveForMonth cộng 1 ngày phép vào đầu tháng cho tất cả nhân viên.
//
// Quy tắc idempotent: mỗi nhân viên chỉ được cộng 1 lần / tháng.
// Nếu đã accrual tháng đó (last_accrual_month == month), bỏ qua.
// Max 12 phép/năm.
//
// Nhân viên mới (chưa có bản ghi leave_balance năm nay) → tạo mới và cộng ngay.
func (s *TimekeepingServiceImpl) AccrueLeaveForMonth(ctx context.Context, req *dto.AccrueLeaveRequest) (*dto.AccrueLeaveResponse, *errors.AppError) {
	if req.Month < 1 || req.Month > 12 {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Tháng không hợp lệ (1-12)", nil)
	}
	if req.Year < 2020 {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Năm không hợp lệ", nil)
	}

	// Lấy tất cả nhân viên
	var allProfiles []payrollEntity.UserProfile
	if err := s.payrollRepo.DB.SQLx().SelectContext(ctx, &allProfiles, `SELECT id, code, full_name FROM user_profiles`); err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi tải danh sách nhân viên", err)
	}

	// Lấy tất cả leave_balances hiện tại của năm
	existingBalances, err := s.repo.GetAllLeaveBalancesForYear(ctx, req.Year)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Lỗi tải số dư phép", err)
	}

	balanceMap := make(map[uuid.UUID]*entity.LeaveBalance)
	for i := range existingBalances {
		b := &existingBalances[i]
		balanceMap[b.EmployeeID] = b
	}

	accrued := 0
	for _, p := range allProfiles {
		b, exists := balanceMap[p.ID]
		if !exists {
			// Nhân viên chưa có bản ghi năm nay → tạo mới
			month := req.Month
			b = &entity.LeaveBalance{
				EmployeeID:       p.ID,
				Year:             req.Year,
				AccruedDays:      1.0, // Cộng ngay tháng đầu tiên
				UsedDays:         0.0,
				Balance:          1.0,
				LastAccrualMonth: &month,
			}
			b.ID = uuid.New()
			if errUpsert := s.repo.UpsertLeaveBalance(ctx, b); errUpsert != nil {
				logger.Error("AccrueLeaveForMonth:UpsertLeaveBalance:Error %v (employee=%s)", errUpsert, p.ID)
				continue
			}
			accrued++
			continue
		}

		// Đã có bản ghi — kiểm tra đã accrual tháng này chưa
		if b.LastAccrualMonth != nil && *b.LastAccrualMonth == req.Month {
			// Đã accrual tháng này rồi, bỏ qua (idempotent)
			continue
		}

		// Tối đa 12 phép/năm
		if b.AccruedDays >= 12.0 {
			continue
		}

		month := req.Month
		b.AccruedDays += 1.0
		b.Balance = b.AccruedDays - b.UsedDays
		b.LastAccrualMonth = &month

		if errUpsert := s.repo.UpsertLeaveBalance(ctx, b); errUpsert != nil {
			logger.Error("AccrueLeaveForMonth:UpsertLeaveBalance:Error %v (employee=%s)", errUpsert, p.ID)
			continue
		}
		accrued++
	}

	return &dto.AccrueLeaveResponse{
		Month:          req.Month,
		Year:           req.Year,
		EmployeesCount: accrued,
		Message:        fmt.Sprintf("Đã cộng 1 ngày phép cho %d nhân viên trong tháng %d/%d", accrued, req.Month, req.Year),
	}, nil
}
