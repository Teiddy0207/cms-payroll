package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// PermissionSeedItem defines a permission to seed
type PermissionSeedItem struct {
	Name        string  // Vietnamese display name
	Slug        string  // resource::action (e.g. role::edit)
	Resource    string  // e.g. roles_fe, permissions_fe
	Action      string  // Vietnamese: Xem, Chỉnh sửa, Xoá, Xét duyệt, Hệ thống
	Description *string // optional
	IsSystem    bool
}

// SeedPermissions seeds permissions if they do not exist (checked by slug)
func SeedPermissions(ctx context.Context, db database.Database) error {
	permCreateUser := "Permission to create users"
	permissions := []PermissionSeedItem{
		// Trang quản trị
		{Name: "Trang quản trị", Slug: "dasboard::read", Resource: "dashboard_fe", Action: "Xem", IsSystem: false},
		// Người dùng
		{Name: "Người dùng", Slug: "user::edit", Resource: "users_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Người dùng", Slug: "user::delete", Resource: "users_fe", Action: "Xoá", IsSystem: false},
		{Name: "Người dùng", Slug: "user::read", Resource: "users_fe", Action: "Xem", IsSystem: false},
		{Name: "create User", Slug: "user::create", Resource: "users", Action: "create", Description: &permCreateUser, IsSystem: false},
		// Vai trò
		{Name: "Vai trò", Slug: "role::read", Resource: "roles_fe", Action: "Xem", IsSystem: false},
		{Name: "Vai trò", Slug: "role::edit", Resource: "roles_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Vai trò", Slug: "role::delete", Resource: "roles_fe", Action: "Xoá", IsSystem: false},
		// Phân quyền
		{Name: "Phân quyền", Slug: "permission::read", Resource: "permissions_fe", Action: "Xem", IsSystem: false},
		{Name: "Phân quyền", Slug: "permission::edit", Resource: "permissions_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phân quyền", Slug: "permission::delete", Resource: "permissions_fe", Action: "Xoá", IsSystem: false},
		// Nhật ký
		{Name: "Nhật ký", Slug: "log::read", Resource: "logs_fe", Action: "Xem", IsSystem: false},
		{Name: "Nhật ký", Slug: "log::delete", Resource: "logs_fe", Action: "Xoá", IsSystem: false},
		// Thư viện năng lực
		{Name: "Thư viện năng lực", Slug: "jobCapability::read", Resource: "job_capabilities_fe", Action: "Xem", IsSystem: false},
		{Name: "Thư viện năng lực", Slug: "jobCapability::edit", Resource: "job_capabilities_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Thư viện năng lực", Slug: "jobCapability::delete", Resource: "job_capabilities_fe", Action: "Xoá", IsSystem: false},
		{Name: "Thư viện năng lực", Slug: "jobCapability::review", Resource: "job_capabilities_fe", Action: "Xét duyệt", IsSystem: false},
		// Danh mục năng lực
		{Name: "Danh mục năng lực", Slug: "catJobCapability::read", Resource: "cat_job_capabilities_fe", Action: "Xem", IsSystem: false},
		{Name: "Danh mục năng lực", Slug: "catJobCapability::edit", Resource: "cat_job_capabilities_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Danh mục năng lực", Slug: "catJobCapability::delete", Resource: "cat_job_capabilities_fe", Action: "Xoá", IsSystem: false},
		{Name: "Danh mục năng lực", Slug: "catJobCapability::review", Resource: "cat_job_capabilities_fe", Action: "Xét duyệt", IsSystem: false},
		// Tiêu chuẩn năng lực của chức danh
		{Name: "Tiêu chuẩn năng lực của chức danh", Slug: "jobRequiredCompetency::read", Resource: "job_required_competencies_fe", Action: "Xem", IsSystem: false},
		{Name: "Tiêu chuẩn năng lực của chức danh", Slug: "jobRequiredCompetency::edit", Resource: "job_required_competencies_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Tiêu chuẩn năng lực của chức danh", Slug: "jobRequiredCompetency::delete", Resource: "job_required_competencies_fe", Action: "Xoá", IsSystem: false},
		{Name: "Tiêu chuẩn năng lực của chức danh", Slug: "jobRequiredCompetency::review", Resource: "job_required_competencies_fe", Action: "Xét duyệt", IsSystem: false},
		// Đánh giá năng lực
		{Name: "Đánh giá năng lực", Slug: "userProfileCompetency::read", Resource: "user_profile_competencies_fe", Action: "Xem", IsSystem: false},
		{Name: "Đánh giá năng lực", Slug: "userProfileCompetency::edit", Resource: "user_profile_competencies_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Đánh giá năng lực", Slug: "userProfileCompetency::delete", Resource: "user_profile_competencies_fe", Action: "Xoá", IsSystem: false},
		{Name: "Đánh giá năng lực", Slug: "userProfileCompetency::review", Resource: "user_profile_competencies_fe", Action: "Xét duyệt", IsSystem: false},
		// Thư viện KPI
		{Name: "Thư viện KPI", Slug: "jobKpi::read", Resource: "job_kpis_fe", Action: "Xem", IsSystem: false},
		{Name: "Thư viện KPI", Slug: "jobKpi::edit", Resource: "job_kpis_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Thư viện KPI", Slug: "jobKpi::delete", Resource: "job_kpis_fe", Action: "Xoá", IsSystem: false},
		{Name: "Thư viện KPI", Slug: "jobKpi::review", Resource: "job_kpis_fe", Action: "Xét duyệt", IsSystem: false},
		// Phân bố tỉ trọng KPI
		{Name: "Phân bố tỉ trọng KPI", Slug: "jobAllocateKpi::read", Resource: "job_allocate_kpis_fe", Action: "Xem", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPI", Slug: "jobAllocateKpi::edit", Resource: "job_allocate_kpis_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPI", Slug: "jobAllocateKpi::delete", Resource: "job_allocate_kpis_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPI", Slug: "jobAllocateKpi::review", Resource: "job_allocate_kpis_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPI", Slug: "jobAllocateKpi::system", Resource: "job_allocate_kpis_fe", Action: "Hệ thống", IsSystem: true},
		// Kết quả KPI
		{Name: "Kết quả KPI", Slug: "kpiResult::read", Resource: "kpi_results_fe", Action: "Xem", IsSystem: false},
		{Name: "Kết quả KPI", Slug: "kpiResult::edit", Resource: "kpi_results_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Kết quả KPI", Slug: "kpiResult::delete", Resource: "kpi_results_fe", Action: "Xoá", IsSystem: false},
		{Name: "Kết quả KPI", Slug: "kpiResult::review", Resource: "kpi_results_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Kết quả KPI", Slug: "kpiResult::system", Resource: "kpi_results_fe", Action: "Hệ thống", IsSystem: true},
		// Kết quả thực hiện KPI của công ty
		{Name: "Kết quả thực hiện KPI của công ty", Slug: "kpiResultSummary::read", Resource: "kpi_result_summaries_fe", Action: "Xem", IsSystem: false},
		{Name: "Kết quả thực hiện KPI của công ty", Slug: "kpiResultSummary::review", Resource: "kpi_result_summaries_fe", Action: "Xét duyệt", IsSystem: false},
		// Mục tiêu KPI
		{Name: "Mục tiêu KPI", Slug: "kpiTarget::read", Resource: "kpi_targets_fe", Action: "Xem", IsSystem: false},
		{Name: "Mục tiêu KPI", Slug: "kpiTarget::edit", Resource: "kpi_targets_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Mục tiêu KPI", Slug: "kpiTarget::delete", Resource: "kpi_targets_fe", Action: "Xoá", IsSystem: false},
		{Name: "Mục tiêu KPI", Slug: "kpiTarget::review", Resource: "kpi_targets_fe", Action: "Xét duyệt", IsSystem: false},
		// Chấm điểm KPI
		{Name: "Chấm điểm KPI", Slug: "kpiScoring::read", Resource: "kpi_scorings_fe", Action: "Xem", IsSystem: false},
		{Name: "Chấm điểm KPI", Slug: "kpiScoring::edit", Resource: "kpi_scorings_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Chấm điểm KPI", Slug: "kpiScoring::delete", Resource: "kpi_scorings_fe", Action: "Xoá", IsSystem: false},
		{Name: "Chấm điểm KPI", Slug: "kpiScoring::review", Resource: "kpi_scorings_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Chấm điểm KPI", Slug: "kpiScoring::system", Resource: "kpi_scorings_fe", Action: "Hệ thống", IsSystem: true},
		// Viễn cảnh
		{Name: "Viễn cảnh", Slug: "jobPerspective::read", Resource: "job_perspectives_fe", Action: "Xem", IsSystem: false},
		{Name: "Viễn cảnh", Slug: "jobPerspective::edit", Resource: "job_perspectives_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Viễn cảnh", Slug: "jobPerspective::delete", Resource: "job_perspectives_fe", Action: "Xoá", IsSystem: false},
		{Name: "Viễn cảnh", Slug: "jobPerspective::review", Resource: "job_perspectives_fe", Action: "Xét duyệt", IsSystem: false},
		// Phân bố tỉ trọng viễn cảnh
		{Name: "Phân bố tỉ trọng viễn cảnh", Slug: "jobAllocatePerspective::read", Resource: "job_allocate_perspectives_fe", Action: "Xem", IsSystem: false},
		{Name: "Phân bố tỉ trọng viễn cảnh", Slug: "jobAllocatePerspective::edit", Resource: "job_allocate_perspectives_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phân bố tỉ trọng viễn cảnh", Slug: "jobAllocatePerspective::delete", Resource: "job_allocate_perspectives_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phân bố tỉ trọng viễn cảnh", Slug: "jobAllocatePerspective::review", Resource: "job_allocate_perspectives_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Phân bố tỉ trọng viễn cảnh", Slug: "jobAllocatePerspective::system", Resource: "job_allocate_perspectives_fe", Action: "Hệ thống", IsSystem: true},
		// KPO
		{Name: "KPO", Slug: "jobKpo::read", Resource: "job_kpos_fe", Action: "Xem", IsSystem: false},
		{Name: "KPO", Slug: "jobKpo::edit", Resource: "job_kpos_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "KPO", Slug: "jobKpo::delete", Resource: "job_kpos_fe", Action: "Xoá", IsSystem: false},
		{Name: "KPO", Slug: "jobKpo::review", Resource: "job_kpos_fe", Action: "Xét duyệt", IsSystem: false},
		// Phân bố tỉ trọng KPO
		{Name: "Phân bố tỉ trọng KPO", Slug: "jobAllocateKpo::read", Resource: "job_allocate_kpos_fe", Action: "Xem", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPO", Slug: "jobAllocateKpo::edit", Resource: "job_allocate_kpos_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPO", Slug: "jobAllocateKpo::delete", Resource: "job_allocate_kpos_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPO", Slug: "jobAllocateKpo::review", Resource: "job_allocate_kpos_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Phân bố tỉ trọng KPO", Slug: "jobAllocateKpo::system", Resource: "job_allocate_kpos_fe", Action: "Hệ thống", IsSystem: true},
		// Vị trí chức danh
		{Name: "Vị trí chức danh", Slug: "jobPosition::read", Resource: "job_positions_fe", Action: "Xem", IsSystem: false},
		{Name: "Vị trí chức danh", Slug: "jobPosition::edit", Resource: "job_positions_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Vị trí chức danh", Slug: "jobPosition::delete", Resource: "job_positions_fe", Action: "Xoá", IsSystem: false},
		{Name: "Vị trí chức danh", Slug: "jobPosition::review", Resource: "job_positions_fe", Action: "Xét duyệt", IsSystem: false},
		// Phòng ban
		{Name: "Phòng ban", Slug: "department::read", Resource: "departments_fe", Action: "Xem", IsSystem: false},
		{Name: "Phòng ban", Slug: "department::edit", Resource: "departments_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phòng ban", Slug: "department::delete", Resource: "departments_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phòng ban", Slug: "department::review", Resource: "departments_fe", Action: "Xét duyệt", IsSystem: false},
		// Tiêu chuẩn công việc
		{Name: "Tiêu chuẩn công việc", Slug: "jobStandard::read", Resource: "job_standards_fe", Action: "Xem", IsSystem: false},
		{Name: "Tiêu chuẩn công việc", Slug: "jobStandard::edit", Resource: "job_standards_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Tiêu chuẩn công việc", Slug: "jobStandard::delete", Resource: "job_standards_fe", Action: "Xoá", IsSystem: false},
		{Name: "Tiêu chuẩn công việc", Slug: "jobStandard::review", Resource: "job_standards_fe", Action: "Xét duyệt", IsSystem: false},
		// Hồ sơ nhân sự
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::read", Resource: "user_profiles_fe", Action: "Xem", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::readAll", Resource: "user_profiles_fe", Action: "Xem tất cả", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::edit", Resource: "user_profiles_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::editAll", Resource: "user_profiles_fe", Action: "Chỉnh sửa tất cả", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::delete", Resource: "user_profiles_fe", Action: "Xoá", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::review", Resource: "user_profiles_fe", Action: "Xét duyệt", IsSystem: false},
		// Khen thưởng
		{Name: "Khen thưởng", Slug: "reward::read", Resource: "rewards_fe", Action: "Xem", IsSystem: false},
		{Name: "Khen thưởng", Slug: "reward::edit", Resource: "rewards_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Khen thưởng", Slug: "reward::delete", Resource: "rewards_fe", Action: "Xoá", IsSystem: false},
		{Name: "Khen thưởng", Slug: "reward::review", Resource: "rewards_fe", Action: "Xét duyệt", IsSystem: false},
		// Phạt
		{Name: "Phạt", Slug: "punish::read", Resource: "punishes_fe", Action: "Xem", IsSystem: false},
		{Name: "Phạt", Slug: "punish::edit", Resource: "punishes_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phạt", Slug: "punish::delete", Resource: "punishes_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phạt", Slug: "punish::review", Resource: "punishes_fe", Action: "Xét duyệt", IsSystem: false},
		// Phụ cấp
		{Name: "Phụ cấp", Slug: "allowance::read", Resource: "allowances_fe", Action: "Xem", IsSystem: false},
		{Name: "Phụ cấp", Slug: "allowance::edit", Resource: "allowances_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Phụ cấp", Slug: "allowance::delete", Resource: "allowances_fe", Action: "Xoá", IsSystem: false},
		{Name: "Phụ cấp", Slug: "allowance::review", Resource: "allowances_fe", Action: "Xét duyệt", IsSystem: false},
		// Lịch làm việc
		{Name: "Lịch làm việc", Slug: "workSchedule::read", Resource: "work_schedules_fe", Action: "Xem", IsSystem: false},
		{Name: "Lịch làm việc", Slug: "workSchedule::edit", Resource: "work_schedules_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Lịch làm việc", Slug: "workSchedule::delete", Resource: "work_schedules_fe", Action: "Xoá", IsSystem: false},
		{Name: "Lịch làm việc", Slug: "workSchedule::review", Resource: "work_schedules_fe", Action: "Xét duyệt", IsSystem: false},
		// Ngày lễ
		{Name: "Ngày lễ", Slug: "holiday::read", Resource: "holidays_fe", Action: "Xem", IsSystem: false},
		{Name: "Ngày lễ", Slug: "holiday::edit", Resource: "holidays_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Ngày lễ", Slug: "holiday::delete", Resource: "holidays_fe", Action: "Xoá", IsSystem: false},
		{Name: "Ngày lễ", Slug: "holiday::review", Resource: "holidays_fe", Action: "Xét duyệt", IsSystem: false},
		// Ngày nghỉ phép
		{Name: "Ngày nghỉ phép", Slug: "dayOfLeave::read", Resource: "day_of_leaves_fe", Action: "Xem", IsSystem: false},
		{Name: "Ngày nghỉ phép", Slug: "dayOfLeave::edit", Resource: "day_of_leaves_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Ngày nghỉ phép", Slug: "dayOfLeave::delete", Resource: "day_of_leaves_fe", Action: "Xoá", IsSystem: false},
		{Name: "Ngày nghỉ phép", Slug: "dayOfLeave::review", Resource: "day_of_leaves_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Ngày nghỉ phép", Slug: "dayOfLeave::system", Resource: "day_of_leaves_fe", Action: "Hệ thống", IsSystem: true},
		// Chấm công ngày
		{Name: "Chấm công ngày", Slug: "dailyTimekeeping::read", Resource: "daily_timekeepings_fe", Action: "Xem", IsSystem: false},
		{Name: "Chấm công ngày", Slug: "dailyTimekeeping::edit", Resource: "daily_timekeepings_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Chấm công ngày", Slug: "dailyTimekeeping::delete", Resource: "daily_timekeepings_fe", Action: "Xoá", IsSystem: false},
		{Name: "Chấm công ngày", Slug: "dailyTimekeeping::review", Resource: "daily_timekeepings_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Chấm công ngày", Slug: "dailyTimekeeping::system", Resource: "daily_timekeepings_fe", Action: "Hệ thống", IsSystem: true},
		// Khấu trừ
		{Name: "Khấu trừ", Slug: "deduction::read", Resource: "deductions_fe", Action: "Xem", IsSystem: false},
		{Name: "Khấu trừ", Slug: "deduction::edit", Resource: "deductions_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Khấu trừ", Slug: "deduction::delete", Resource: "deductions_fe", Action: "Xoá", IsSystem: false},
		{Name: "Khấu trừ", Slug: "deduction::review", Resource: "deductions_fe", Action: "Xét duyệt", IsSystem: false},
		// Bảng công
		{Name: "Bảng công", Slug: "timekeepingSheet::read", Resource: "timekeeping_sheets_fe", Action: "Xem", IsSystem: false},
		{Name: "Bảng công", Slug: "timekeepingSheet::edit", Resource: "timekeeping_sheets_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Bảng công", Slug: "timekeepingSheet::delete", Resource: "timekeeping_sheets_fe", Action: "Xoá", IsSystem: false},
		{Name: "Bảng công", Slug: "timekeepingSheet::review", Resource: "timekeeping_sheets_fe", Action: "Xét duyệt", IsSystem: false},
		// Chi nhánh
		{Name: "Chi nhánh", Slug: "branch::read", Resource: "branches_fe", Action: "Xem", IsSystem: false},
		{Name: "Chi nhánh", Slug: "branch::edit", Resource: "branches_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Chi nhánh", Slug: "branch::delete", Resource: "branches_fe", Action: "Xoá", IsSystem: false},
		// Bảng lương
		{Name: "Bảng lương", Slug: "salary::read", Resource: "salaries_fe", Action: "Xem", IsSystem: false},
		{Name: "Bảng lương", Slug: "salary::edit", Resource: "salaries_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Bảng lương", Slug: "salary::delete", Resource: "salaries_fe", Action: "Xoá", IsSystem: false},
		{Name: "Bảng lương", Slug: "salary::review", Resource: "salaries_fe", Action: "Xét duyệt", IsSystem: false},
		// Quỹ lương bổ sung
		{Name: "Quỹ lương bổ sung", Slug: "payrollFundSupplementary::read", Resource: "payroll_fund_supplementaries_fe", Action: "Xem", IsSystem: false},
		{Name: "Quỹ lương bổ sung", Slug: "payrollFundSupplementary::edit", Resource: "payroll_fund_supplementaries_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Quỹ lương bổ sung", Slug: "payrollFundSupplementary::delete", Resource: "payroll_fund_supplementaries_fe", Action: "Xoá", IsSystem: false},
		{Name: "Quỹ lương bổ sung", Slug: "payrollFundSupplementary::review", Resource: "payroll_fund_supplementaries_fe", Action: "Xét duyệt", IsSystem: false},
		// Công thức động
		{Name: "Công thức động", Slug: "formulaDynamic::read", Resource: "formula_dynamics_fe", Action: "Xem", IsSystem: false},
		{Name: "Công thức động", Slug: "formulaDynamic::edit", Resource: "formula_dynamics_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Công thức động", Slug: "formulaDynamic::delete", Resource: "formula_dynamics_fe", Action: "Xoá", IsSystem: false},
		{Name: "Công thức động", Slug: "formulaDynamic::review", Resource: "formula_dynamics_fe", Action: "Xét duyệt", IsSystem: false},
		// Báo cáo
		{Name: "Báo cáo hệ thống", Slug: "reportSystem::read", Resource: "report_systems_fe", Action: "Xem", IsSystem: false},
		{Name: "Báo cáo phòng ban", Slug: "reportDepartment::read", Resource: "report_departments_fe", Action: "Xem", IsSystem: false},
		{Name: "Báo cáo cá nhân", Slug: "reportSelf::read", Resource: "report_selves_fe", Action: "Xem", IsSystem: false},
		{Name: "Báo cáo KPI", Slug: "reportKPI::read", Resource: "report_kpis_fe", Action: "Xem", IsSystem: false},
		{Name: "Báo cáo lương", Slug: "reportPayroll::read", Resource: "report_payrolls_fe", Action: "Xem", IsSystem: false},
		// Cài đặt
		{Name: "Cài đặt", Slug: "setting::read", Resource: "settings_fe", Action: "Xem", IsSystem: false},
		{Name: "Cài đặt", Slug: "setting::edit", Resource: "settings_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Cài đặt", Slug: "setting::delete", Resource: "settings_fe", Action: "Xoá", IsSystem: false},
		{Name: "Cài đặt", Slug: "setting::review", Resource: "settings_fe", Action: "Xét duyệt", IsSystem: false},
		// Tệp tin
		{Name: "Tệp tin", Slug: "storage::edit", Resource: "storages_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Tệp tin", Slug: "storage::delete", Resource: "storages_fe", Action: "Xoá", IsSystem: false},
		// Báo cáo HR
		{Name: "Báo cáo HR", Slug: "reportHR::read", Resource: "report_hrs_fe", Action: "Xem", IsSystem: false},
		// Loại KPI / Đơn vị KPI / Tổng hợp KPI (2026-04-06)
		{Name: "Loại KPI", Slug: "kpiType::read", Resource: "kpi_types_fe", Action: "Xem", IsSystem: false},
		{Name: "Loại KPI", Slug: "kpiType::edit", Resource: "kpi_types_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Loại KPI", Slug: "kpiType::delete", Resource: "kpi_types_fe", Action: "Xoá", IsSystem: false},
		{Name: "Đơn vị KPI", Slug: "kpiUnit::read", Resource: "kpi_units_fe", Action: "Xem", IsSystem: false},
		{Name: "Đơn vị KPI", Slug: "kpiUnit::edit", Resource: "kpi_units_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Đơn vị KPI", Slug: "kpiUnit::delete", Resource: "kpi_units_fe", Action: "Xoá", IsSystem: false},
		{Name: "Tổng hợp KPI", Slug: "kpiAggregate::read", Resource: "kpi_aggregates_fe", Action: "Xem", IsSystem: false},
		{Name: "Tổng hợp KPI", Slug: "kpiAggregate::edit", Resource: "kpi_aggregates_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Tổng hợp KPI", Slug: "kpiAggregate::delete", Resource: "kpi_aggregates_fe", Action: "Xoá", IsSystem: false},
		{Name: "Tổng hợp KPI", Slug: "kpiAggregate::review", Resource: "kpi_aggregates_fe", Action: "Xét duyệt", IsSystem: false},
		{Name: "Hội đồng cấp cao", Slug: "highCouncil::all", Resource: "representatives_fe", Action: "Đại diện", IsSystem: false},
		{Name: "Giám đốc", Slug: "director::all", Resource: "representatives_fe", Action: "Đại diện", IsSystem: false},
		{Name: "Công ty", Slug: "companyRepresentative::all", Resource: "representatives_fe", Action: "Đại diện", IsSystem: false},
	}

	inserted := 0
	skipped := 0

	for _, p := range permissions {
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM permissions WHERE slug = $1 LIMIT 1`, p.Slug)
		if err == nil && existingID != uuid.Nil {
			skipped++
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedPermissions: Error checking existing permission", "slug", p.Slug, "error", err)
			return fmt.Errorf("check permission %s: %w", p.Slug, err)
		}

		id := uuid.New()
		query := `INSERT INTO permissions (id, name, slug, resource, action, description, is_system, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`
		err = db.ExecContext(ctx, query, id, p.Name, p.Slug, p.Resource, p.Action, p.Description, p.IsSystem)
		if err != nil {
			logger.Error("SeedPermissions: Failed to insert permission", "slug", p.Slug, "error", err)
			return fmt.Errorf("insert permission %s: %w", p.Slug, err)
		}

		inserted++
		logger.Info("SeedPermissions: Created permission", "slug", p.Slug, "name", p.Name)
	}

	logger.Info("SeedPermissions: Completed",
		"inserted", inserted,
		"skipped", skipped,
		"total", len(permissions))
	return nil
}
