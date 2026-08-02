package database

import (
	"cal-salary/core/constants"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IDatabase interface {
	ExecContext(ctx context.Context, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	NamedQueryContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	SQLx() *sqlx.DB
}

type Database struct {
	db   *sql.DB
	sqlx *sqlx.DB
}

type DatabaseConfig struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	DBName                 string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetime        int    // in minutes
	SSLMode                string // disable, require, verify-ca, verify-full
	ConnectTimeout         int    // in seconds
	StatementTimeout       int    // in seconds
	IdleInTxSessionTimeout int    // in seconds
}

var (
	instance *Database
)

func GetDB() IDatabase {
	return instance
}

func runSchemaInit(db *Database) error {
	logger.Info("Running schema initialization...")
	schemaFile := "db/init_schema.sql"
	content, err := os.ReadFile(schemaFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("Schema file not found, skipping initialization", "path", schemaFile)
			return nil
		}
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.sqlx.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema DDL: %w", err)
	}

	_, _ = db.sqlx.Exec("ALTER TABLE employee_face_templates ADD COLUMN IF NOT EXISTS face_embedding jsonb;")
	_, _ = db.sqlx.Exec("ALTER TABLE competency_dictionaries ADD COLUMN IF NOT EXISTS point_value INTEGER NOT NULL DEFAULT 0;")
	_, _ = db.sqlx.Exec("ALTER TABLE competency_dictionaries ALTER COLUMN point_value TYPE INTEGER USING point_value::integer;")
	_, _ = db.sqlx.Exec(`CREATE TABLE IF NOT EXISTS employee_competencies (
		user_profile_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
		competency_id UUID NOT NULL REFERENCES competency_dictionaries(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		PRIMARY KEY (user_profile_id, competency_id)
	);`)

	_, _ = db.sqlx.Exec(`CREATE TABLE IF NOT EXISTS meeting_summaries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
		summary TEXT,
		key_decisions TEXT,
		action_items TEXT,
		efficiency_score VARCHAR(50),
		sentiment VARCHAR(50),
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		CONSTRAINT meeting_summaries_meeting_id_key UNIQUE (meeting_id)
	);`)
	_, _ = db.sqlx.Exec("ALTER TABLE meeting_summaries ALTER COLUMN key_decisions TYPE TEXT USING key_decisions::text;")
	_, _ = db.sqlx.Exec("ALTER TABLE meeting_summaries ALTER COLUMN action_items TYPE TEXT USING action_items::text;")

	_, _ = db.sqlx.Exec("ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS p2_base_allowance DECIMAL(15, 2) NOT NULL DEFAULT 0.00;")
	_, _ = db.sqlx.Exec("ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS p2_cap DECIMAL(15, 2) NOT NULL DEFAULT 0.00;")
	_, _ = db.sqlx.Exec("UPDATE users SET password = '$2a$10$1G9Y7UIa9SJIwA0SxLgu7.3YkVKaQ/HtCATAAGfxCLzuO5wADcmzq' WHERE username = 'admin';")
	_, _ = db.sqlx.Exec(`
		INSERT INTO role_permissions (role_id, permission_id, granted_by, granted_at)
		SELECT '10000000-0000-0000-0000-000000000001', id, '00000000-0000-0000-0000-000000000000', NOW()
		FROM permissions
		ON CONFLICT (role_id, permission_id) DO NOTHING;
	`)
	_, _ = db.sqlx.Exec("DELETE FROM role_permissions WHERE role_id = '10000000-0000-0000-0000-000000000003';")
	_, _ = db.sqlx.Exec(`
		INSERT INTO role_permissions (role_id, permission_id, granted_by, granted_at)
		SELECT '10000000-0000-0000-0000-000000000003', id, '00000000-0000-0000-0000-000000000000', NOW()
		FROM permissions
		WHERE slug IN (
			'dashboard::read',
			'dailyTimekeeping::read'
		)
		ON CONFLICT (role_id, permission_id) DO NOTHING;
	`)

	descMap := map[string]string{
		"role::read":                        "Quyền này dùng để Xem danh sách vai trò, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"role::edit":                        "Quyền này dùng để Tạo mới và Chỉnh sửa vai trò, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"role::delete":                      "Quyền này dùng để Xoá vai trò, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"permission::read":                  "Quyền này dùng để Xem danh sách quyền hạn, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"permission::edit":                  "Quyền này dùng để Cấu hình gán quyền hạn, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"permission::delete":                "Quyền này dùng để Xoá quyền hạn, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"dashboard::read":                   "Quyền này dùng để Xem báo cáo tổng quan, cho màn Dashboard thuộc menu Thống kê sidebar main",
		"user::read":                        "Quyền này dùng để Xem danh sách tài khoản, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"user::edit":                        "Quyền này dùng để Chỉnh sửa tài khoản, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"user::delete":                      "Quyền này dùng để Khoá/Xoá tài khoản, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"user::create":                      "Quyền này dùng để Tạo mới tài khoản, cho màn Phân quyền thuộc menu Hệ thống sidebar main",
		"log::read":                         "Quyền này dùng để Xem nhật ký hoạt động, cho màn Nhật ký chấm công thuộc menu Thống kê sidebar main",
		"log::delete":                       "Quyền này dùng để Xoá nhật ký hoạt động, cho màn Nhật ký chấm công thuộc menu Thống kê sidebar main",
		"department::read":                  "Quyền này dùng để Xem danh sách phòng ban, cho màn Phòng ban thuộc menu Quản lý nhân sự sidebar main",
		"department::edit":                  "Quyền này dùng để Chỉnh sửa phòng ban, cho màn Phòng ban thuộc menu Quản lý nhân sự sidebar main",
		"department::delete":                "Quyền này dùng để Xoá phòng ban, cho màn Phòng ban thuộc menu Quản lý nhân sự sidebar main",
		"jobPosition::read":                 "Quyền này dùng để Xem danh sách vị trí công việc, cho màn Vị trí công việc thuộc menu Đánh giá sidebar main",
		"jobPosition::edit":                 "Quyền này dùng để Chỉnh sửa vị trí công việc, cho màn Vị trí công việc thuộc menu Đánh giá sidebar main",
		"jobPosition::delete":               "Quyền này dùng để Xoá vị trí công việc, cho màn Vị trí công việc thuộc menu Đánh giá sidebar main",
		"userProfile::read":                 "Quyền này dùng để Xem thông tin nhân sự và hợp đồng, cho màn Nhân viên thuộc menu Quản lý nhân sự sidebar main",
		"userProfile::edit":                 "Quyền này dùng để Chỉnh sửa thông tin nhân sự và hợp đồng, cho màn Nhân viên thuộc menu Quản lý nhân sự sidebar main",
		"userProfile::delete":               "Quyền này dùng để Xoá thông tin nhân sự và hợp đồng, cho màn Nhân viên thuộc menu Quản lý nhân sự sidebar main",
		"jobCapability::read":               "Quyền này dùng để Xem thư viện năng lực, cho màn Năng lực thuộc menu Đánh giá sidebar main",
		"jobCapability::edit":               "Quyền này dùng để Chỉnh sửa thư viện năng lực, cho màn Năng lực thuộc menu Đánh giá sidebar main",
		"jobCapability::delete":             "Quyền này dùng để Xoá thư viện năng lực, cho màn Năng lực thuộc menu Đánh giá sidebar main",
		"userProfileCompetency::read":       "Quyền này dùng để Xem đánh giá năng lực nhân viên, cho màn Đánh giá năng lực thuộc menu Đánh giá sidebar main",
		"userProfileCompetency::edit":       "Quyền này dùng để Thực hiện đánh giá năng lực nhân viên, cho màn Đánh giá năng lực thuộc menu Đánh giá sidebar main",
		"userProfileCompetency::delete":     "Quyền này dùng để Xoá đánh giá năng lực nhân viên, cho màn Đánh giá năng lực thuộc menu Đánh giá sidebar main",
		"salary::read":                      "Quyền này dùng để Xem bảng lương nhân viên, cho màn Tính lương tháng thuộc menu Lương sidebar main",
		"salary::edit":                      "Quyền này dùng để Chạy tính toán lương và cập nhật lương, cho màn Tính lương tháng thuộc menu Lương sidebar main",
		"salary::delete":                    "Quyền này dùng để Xoá bản ghi bảng lương, cho màn Tính lương tháng thuộc menu Lương sidebar main",
		"salary::review":                    "Quyền này dùng để Xét duyệt bảng lương tháng, cho màn Tính lương tháng thuộc menu Lương sidebar main",
		"timekeepingSheet::read":            "Quyền này dùng để Xem bảng tổng hợp công, cho màn Bảng chấm công thuộc menu Chấm công sidebar main",
		"timekeepingSheet::edit":            "Quyền này dùng để Chốt và tính toán bảng công, cho màn Bảng chấm công thuộc menu Chấm công sidebar main",
		"timekeepingSheet::delete":          "Quyền này dùng để Xoá bảng tổng hợp công, cho màn Bảng chấm công thuộc menu Chấm công sidebar main",
		"dailyTimekeeping::read":            "Quyền này dùng để Xem nhật ký chấm công ngày, cho màn Nhật ký chấm công thuộc menu Chấm công sidebar main",
		"dailyTimekeeping::edit":            "Quyền này dùng để Điều chỉnh check-in và đăng ký khuôn mặt, cho màn Nhật ký chấm công thuộc menu Chấm công sidebar main",
		"dailyTimekeeping::delete":          "Quyền này dùng để Xoá lượt chấm công ngày, cho màn Nhật ký chấm công thuộc menu Chấm công sidebar main",
		"formulaDynamic::read":              "Quyền này dùng để Xem công thức tính lương động, cho màn Công thức lương thuộc menu Lương sidebar main",
		"formulaDynamic::edit":              "Quyền này dùng để Thiết lập và chỉnh sửa công thức tính lương động, cho màn Công thức lương thuộc menu Lương sidebar main",
		"formulaDynamic::delete":            "Quyền này dùng để Xoá công thức tính lương động, cho màn Công thức lương thuộc menu Lương sidebar main",
		"setting::read":                     "Quyền này dùng để Xem cấu hình hệ thống, cho màn Cài đặt thuộc menu Hệ thống sidebar main",
		"setting::edit":                     "Quyền này dùng để Cập nhật cấu hình hệ thống, cho màn Cài đặt thuộc menu Hệ thống sidebar main",
		"setting::delete":                   "Quyền này dùng để Reset cấu hình hệ thống, cho màn Cài đặt thuộc menu Hệ thống sidebar main",
		"storage::edit":                     "Quyền này dùng để Thực hiện nén, gộp, đóng dấu PDF, cho màn Công cụ PDF thuộc menu Hỗ trợ sidebar main",
		"storage::delete":                   "Quyền này dùng để Xoá các tệp tin trong hệ thống, cho màn Công cụ PDF thuộc menu Hỗ trợ sidebar main",
	}

	for slug, desc := range descMap {
		_, _ = db.sqlx.Exec("UPDATE permissions SET description = $1 WHERE slug = $2;", desc, slug)
	}

	_, _ = db.sqlx.Exec(`
		DELETE FROM permissions 
		WHERE slug LIKE 'jobKpi%' 
		   OR slug LIKE 'kpi%' 
		   OR slug LIKE 'jobPerspective%' 
		   OR slug LIKE 'jobAllocatePerspective%' 
		   OR slug LIKE 'jobKpo%' 
		   OR slug LIKE 'jobAllocateKpo%' 
		   OR slug LIKE 'reportKPI%'
		   OR slug LIKE 'allowance%'
		   OR slug LIKE 'punish%'
		   OR slug LIKE 'workSchedule%'
		   OR slug LIKE 'holiday%'
		   OR slug LIKE 'reward%';
	`)

	logger.Info("Schema initialized successfully!")
	return nil
}

func InitDB(config DatabaseConfig) (Database, error) {
	logger.Info("Initializing database...")
	var err error

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, constants.DatabaseSSLMode)

	sqlxDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return Database{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB := sqlxDB.DB
	sqlDB.SetMaxOpenConns(constants.DatabaseMaxOpenConns)
	sqlDB.SetMaxIdleConns(constants.DatabaseMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(constants.DatabaseConnMaxLifetime) * time.Minute)

	if err = sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err)
		sqlDB.Close()
		return Database{}, fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance := Database{
		db:   sqlDB,
		sqlx: sqlxDB,
	}
	instance = &dbInstance

	// Execute Schema Initialization
	if err = runSchemaInit(instance); err != nil {
		logger.Error("Schema initialization failed", "error", err)
		sqlDB.Close()
		instance = nil
		return Database{}, fmt.Errorf("schema initialization failed: %w", err)
	}

	logger.Info("Database initialized successfully",
		"maxOpenConns", constants.DatabaseMaxOpenConns,
		"maxIdleConns", constants.DatabaseMaxIdleConns,
		"connMaxLifetime", constants.DatabaseConnMaxLifetime,
	)

	return *instance, nil
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...any) error {
	_, err := d.sqlx.ExecContext(ctx, query, args...)
	return err
}

func (d *Database) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return d.sqlx.GetContext(ctx, dest, query, args...)
}

func (d *Database) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return d.sqlx.SelectContext(ctx, dest, query, args...)
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *Database) NamedQueryContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error) {
	return d.sqlx.NamedQueryContext(ctx, query, arg)
}

func (d *Database) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	return d.sqlx.NamedExecContext(ctx, query, arg)
}

func (d *Database) SQLx() *sqlx.DB {
	return d.sqlx
}
