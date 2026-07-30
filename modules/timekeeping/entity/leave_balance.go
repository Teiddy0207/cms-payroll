package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

// LeaveBalance theo dõi số dư phép của từng nhân viên từng năm.
//
// Quy tắc nghiệp vụ:
//   - Đầu mỗi tháng (ngày 1), AccruedDays tăng thêm 1.0 (tối đa 12/năm).
//   - UsedDays được cập nhật sau mỗi lần CalculateTimesheets, bằng tổng số ngày
//     có status = LEAVE_PAID trong daily_attendance_sheets cho năm đó.
//   - Balance = AccruedDays - UsedDays (luôn >= 0).
//   - Cuối năm (31/12 → 1/1 năm mới): bản ghi năm mới được tạo với accrued=0.
type LeaveBalance struct {
	EmployeeID       uuid.UUID `db:"employee_id"        json:"employee_id"`
	Year             int       `db:"year"               json:"year"`
	AccruedDays      float64   `db:"accrued_days"       json:"accrued_days"`
	UsedDays         float64   `db:"used_days"          json:"used_days"`
	Balance          float64   `db:"balance"            json:"balance"`
	LastAccrualMonth *int      `db:"last_accrual_month" json:"last_accrual_month"` // 1-12
	entity.BaseEntity
}
