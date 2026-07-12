-- Database Schema Initializer for cal_salary
-- All tables are created in dependency order.

-- 1. Job Descriptions (Positions)
CREATE TABLE IF NOT EXISTS job_descriptions (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    e_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    c_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    r_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    we_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000,
    wc_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000,
    wr_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000,
    salary_spread DECIMAL(5, 4) NOT NULL DEFAULT 0.0000,
    job_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    midpoint DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    min_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    max_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 2. Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    username VARCHAR(100) UNIQUE,
    password VARCHAR(255) NOT NULL,
    email_verified_at TIMESTAMP,
    position_id UUID REFERENCES job_descriptions(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 3. User Profiles (Employees)
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(255),
    date_of_birth VARCHAR(50),
    gender VARCHAR(10),
    position_id UUID REFERENCES job_descriptions(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 4. Permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 5. Roles
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 6. User Roles
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100),
    name VARCHAR(100),
    description TEXT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 7. Role Permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- 8. User Permissions
CREATE TABLE IF NOT EXISTS user_permissions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted BOOLEAN NOT NULL DEFAULT true,
    granted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP
);

-- 9. Departments (Org structure)
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    manager_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Head of Department (Trưởng phòng)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Alter user_profiles and job_descriptions to reference department if not exists
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL;

-- 10. Contracts
CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    contract_code VARCHAR(100) UNIQUE NOT NULL,
    position_base_rate DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, EXPIRED, TERMINATED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 11. Job Standards (for P1)
CREATE TABLE IF NOT EXISTS job_standards (
    id UUID PRIMARY KEY,
    standard_code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    allowance_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 12. Work Standards Seed Table
CREATE TABLE IF NOT EXISTS work_standards (
    id UUID PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL,
    "order" VARCHAR(100) NOT NULL,
    name TEXT NOT NULL,
    parent_id UUID REFERENCES work_standards(id) ON DELETE SET NULL,
    score_require BOOLEAN NOT NULL DEFAULT false,
    level INTEGER NOT NULL DEFAULT 0,
    score DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    description TEXT,
    grandted_by UUID,
    is_default BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 13. Type Competencies
CREATE TABLE IF NOT EXISTS type_competencies (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 14. Competency Dictionaries
CREATE TABLE IF NOT EXISTS competency_dictionaries (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type_id UUID NOT NULL REFERENCES type_competencies(id) ON DELETE CASCADE,
    description TEXT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    proficiency_level_max DECIMAL(5, 2) NOT NULL DEFAULT 5.00,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 15. Competency Evaluations (for P2)
CREATE TABLE IF NOT EXISTS competency_evaluations (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    evaluator_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    evaluation_period VARCHAR(50) NOT NULL, -- e.g. 2026-H1
    competency_id UUID NOT NULL REFERENCES competency_dictionaries(id) ON DELETE CASCADE,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    weight INTEGER NOT NULL CHECK (weight >= 1 AND weight <= 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 16. Daily Attendance Sheets (Posgres timekeeping details)
CREATE TABLE IF NOT EXISTS daily_attendance_sheets (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in TIMESTAMP,
    check_out TIMESTAMP,
    actual_work_day DECIMAL(3, 2) NOT NULL DEFAULT 0.00,
    ot_hours DECIMAL(5, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'ABSENT', -- NORMAL, LATE, EARLY, ABSENT
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (employee_id, date)
);

-- 17. Explanation Requests (Tờ trình bù công)
CREATE TABLE IF NOT EXISTS explanation_requests (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    reason TEXT NOT NULL,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 18. OT Requests (Đăng ký OT)
CREATE TABLE IF NOT EXISTS ot_requests (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    hours_requested DECIMAL(5, 2) NOT NULL,
    is_night_ot BOOLEAN NOT NULL DEFAULT false,
    is_holiday_ot BOOLEAN NOT NULL DEFAULT false,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 19. Payroll Formulas (Formula Versioning)
CREATE TABLE IF NOT EXISTS payroll_formulas (
    id UUID PRIMARY KEY,
    variable_name VARCHAR(100) NOT NULL,
    expression TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 20. Payroll Periods
CREATE TABLE IF NOT EXISTS payroll_periods (
    id UUID PRIMARY KEY,
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- DRAFT, CALCULATED, LOCKED, DISBURSED
    locked_at TIMESTAMP,
    disbursed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (month, year)
);

-- 21. Payroll Records
CREATE TABLE IF NOT EXISTS payroll_records (
    id UUID PRIMARY KEY,
    period_id UUID NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    p1_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    p2_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    p3_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    gross_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    tax DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    net_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- DRAFT, ADJUSTING, APPROVED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (period_id, employee_id)
);

-- 22. Payroll Adjustments
CREATE TABLE IF NOT EXISTS payroll_adjustments (
    id UUID PRIMARY KEY,
    record_id UUID NOT NULL REFERENCES payroll_records(id) ON DELETE CASCADE,
    field_modified VARCHAR(100) NOT NULL,
    old_value DECIMAL(15, 2) NOT NULL,
    new_value DECIMAL(15, 2) NOT NULL,
    reason TEXT NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 23. AI Anomalies
CREATE TABLE IF NOT EXISTS ai_anomalies (
    id UUID PRIMARY KEY,
    record_id UUID NOT NULL REFERENCES payroll_records(id) ON DELETE CASCADE,
    metric_flagged VARCHAR(100) NOT NULL,
    deviation_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    ai_explanation TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, RESOLVED, IGNORED
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 24. RAG Documents
CREATE TABLE IF NOT EXISTS rag_documents (
    id UUID PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    chunk_content TEXT NOT NULL,
    vector_embedding JSONB, -- JSON representation of embedding vector
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 25. Chat Sessions
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 26. Activity Logs (for modules/activity_log)
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    entity_name VARCHAR(100),
    entity_id VARCHAR(100),
    old_data JSONB,
    new_data JSONB,
    reason TEXT,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 27. Attendance Logs (Raw check-in data)
CREATE TABLE IF NOT EXISTS attendance_logs (
    id UUID PRIMARY KEY,
    employee_code VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    location_gps VARCHAR(255) NOT NULL,
    device_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Idempotency key for the JetStream check-in consumer (ON CONFLICT DO NOTHING dedup).
ALTER TABLE attendance_logs ADD COLUMN IF NOT EXISTS event_id UUID;
UPDATE attendance_logs SET event_id = id WHERE event_id IS NULL;
ALTER TABLE attendance_logs ALTER COLUMN event_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_attendance_logs_event_id ON attendance_logs(event_id);

-- 28. Employee Face Templates
CREATE TABLE IF NOT EXISTS employee_face_templates (
    id UUID PRIMARY KEY,
    employee_code VARCHAR(100) NOT NULL UNIQUE,
    face_data TEXT NOT NULL,
    face_embedding jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Default system admin user seeding (required for standard seeds to refer to)
INSERT INTO job_descriptions (id, code, name, description)
VALUES ('00000000-0000-0000-0000-000000000000', 'ADMIN_POS', 'Administrator', 'System administrator position')
ON CONFLICT (code) DO NOTHING;

INSERT INTO users (id, email, username, password, position_id, is_active)
VALUES ('00000000-0000-0000-0000-000000000000', 'admin@example.com', 'admin', '$2a$10$7/Zf9Y.yD8l.y6K4L.9tLeqJ2d3DkP.L6aK7zTj.b.V6D5X2uG.8K', '00000000-0000-0000-0000-000000000000', true)
ON CONFLICT (username) DO NOTHING;

-- 29. System Settings table
CREATE TABLE IF NOT EXISTS system_settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT NOT NULL,
    description TEXT,
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO system_settings (key, value, description)
VALUES 
('payroll_k_factor', '4000000', 'Hệ số quy đổi lương P1 (K factor) VND/điểm')
ON CONFLICT (key) DO NOTHING;

-- Ensure existing database has the P1 range columns
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS e_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS c_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS r_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS we_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS wc_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS wr_weight DECIMAL(5, 4) NOT NULL DEFAULT 0.0000;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS salary_spread DECIMAL(5, 4) NOT NULL DEFAULT 0.0000;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS job_score DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS midpoint DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS min_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS max_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS is_benchmark BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS market_salary DECIMAL(15, 2) NOT NULL DEFAULT 0.00;
ALTER TABLE job_descriptions ADD COLUMN IF NOT EXISTS search_keyword VARCHAR(100);

DELETE FROM system_settings WHERE key = 'company_point_rate';

