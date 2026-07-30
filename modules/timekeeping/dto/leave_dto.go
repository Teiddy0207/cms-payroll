package dto

import "time"

// ==========================================
// Leave Request DTOs
// ==========================================

// CreateLeaveRequest — payload nhân viên gửi khi xin nghỉ phép.
// DaysRequested: 0.5 cho nửa ngày, 1.0+ cho nguyên ngày.
// Nếu DaysRequested == 0 , backend tự tính = số ngày trong [StartDate, EndDate].
type CreateLeaveRequest struct {
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	DaysRequested float64   `json:"days_requested"` // 0.5 hoặc số nguyên
	Reason        string    `json:"reason"`
}

// LeaveRequestResponse — response trả về sau khi tạo/lấy đơn nghỉ phép.
type LeaveRequestResponse struct {
	ID            string    `json:"id"`
	EmployeeID    string    `json:"employee_id"`
	EmployeeCode  string    `json:"employee_code"`
	FullName      string    `json:"full_name"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	DaysRequested float64   `json:"days_requested"`
	Reason        string    `json:"reason"`
	ApprovedBy    *string   `json:"approved_by,omitempty"`
	ApprovedName  *string   `json:"approved_name,omitempty"`
	Status        string    `json:"status"` // PENDING, APPROVED, REJECTED
	CreatedAt     time.Time `json:"created_at"`
}

// ==========================================
// Leave Balance DTOs
// ==========================================

// LeaveBalanceResponse — response trả về số dư phép của nhân viên.
type LeaveBalanceResponse struct {
	EmployeeID       string  `json:"employee_id"`
	EmployeeCode     string  `json:"employee_code"`
	FullName         string  `json:"full_name"`
	Year             int     `json:"year"`
	AccruedDays      float64 `json:"accrued_days"`  // Tổng phép đã tích lũy
	UsedDays         float64 `json:"used_days"`     // Đã dùng
	Balance          float64 `json:"balance"`       // Còn lại
	LastAccrualMonth *int    `json:"last_accrual_month,omitempty"`
}

// AccrueLeaveRequest — payload cho admin gọi accrual thủ công.
type AccrueLeaveRequest struct {
	Month int `json:"month"` // 1-12
	Year  int `json:"year"`
}

// AccrueLeaveResponse — kết quả sau khi chạy accrual.
type AccrueLeaveResponse struct {
	Month          int    `json:"month"`
	Year           int    `json:"year"`
	EmployeesCount int    `json:"employees_count"` // Số nhân viên được cộng phép
	Message        string `json:"message"`
}
