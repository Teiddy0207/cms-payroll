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

		// Hồ sơ nhân sự
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::read", Resource: "user_profiles_fe", Action: "Xem", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::readAll", Resource: "user_profiles_fe", Action: "Xem tất cả", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::edit", Resource: "user_profiles_fe", Action: "Chỉnh sửa", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::editAll", Resource: "user_profiles_fe", Action: "Chỉnh sửa tất cả", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::delete", Resource: "user_profiles_fe", Action: "Xoá", IsSystem: false},
		{Name: "Hồ sơ nhân sự", Slug: "userProfile::review", Resource: "user_profiles_fe", Action: "Xét duyệt", IsSystem: false},
		// Khen thưởng

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
