package repository

import (
	"cal-salary/core/database"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type PayrollRepository struct {
	DB database.Database
}

func NewPayrollRepository(db database.Database) *PayrollRepository {
	return &PayrollRepository{DB: db}
}

type PayrollRepositoryInterface interface {
	// Departments
	CreateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error)
	GetDepartments(ctx context.Context, params params.QueryParams) ([]entity.Department, int, error)
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*entity.Department, error)
	UpdateDepartment(ctx context.Context, id uuid.UUID, dept *entity.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error

	// Job Positions
	CreateJobPosition(ctx context.Context, pos *entity.JobPosition) (*entity.JobPosition, error)
	GetJobPositions(ctx context.Context, params params.QueryParams) ([]entity.JobPosition, int, error)
	GetJobPositionByID(ctx context.Context, id uuid.UUID) (*entity.JobPosition, error)
	UpdateJobPosition(ctx context.Context, id uuid.UUID, pos *entity.JobPosition) error
	DeleteJobPosition(ctx context.Context, id uuid.UUID) error
	GetBenchmarkJobPositions(ctx context.Context) ([]entity.JobPosition, error)
	GetAllJobPositions(ctx context.Context) ([]entity.JobPosition, error)

	// User Profiles
	CreateUserProfile(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error)
	GetUserProfiles(ctx context.Context, params params.QueryParams) ([]entity.UserProfile, int, error)
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context, id uuid.UUID, profile *entity.UserProfile) error
	DeleteUserProfile(ctx context.Context, id uuid.UUID) error

	// Contracts
	CreateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error)
	GetContracts(ctx context.Context, params params.QueryParams) ([]entity.Contract, int, error)
	GetContractByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error)
	UpdateContract(ctx context.Context, id uuid.UUID, contract *entity.Contract) error
	DeleteContract(ctx context.Context, id uuid.UUID) error
	HasActiveContract(ctx context.Context, employeeID uuid.UUID, excludeContractID *uuid.UUID) (bool, error)
	GetActiveContractByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*entity.Contract, error)



	// Competency Evaluations
	CreateEvaluation(ctx context.Context, eval *entity.CompetencyEvaluation) (*entity.CompetencyEvaluation, error)
	GetEvaluations(ctx context.Context, params params.QueryParams) ([]entity.CompetencyEvaluation, int, error)
	GetEvaluationByID(ctx context.Context, id uuid.UUID) (*entity.CompetencyEvaluation, error)
	UpdateEvaluation(ctx context.Context, id uuid.UUID, eval *entity.CompetencyEvaluation) error
	DeleteEvaluation(ctx context.Context, id uuid.UUID) error
	DeleteEvaluationsByPeriodAndEmployee(ctx context.Context, employeeID uuid.UUID, period string) error
	GetManagedDepartmentID(ctx context.Context, managerID uuid.UUID) (uuid.UUID, error)
	GetEmployeeDepartmentID(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	IsAdminOrDirector(ctx context.Context, userID uuid.UUID) (bool, error)



	// Employee Competencies
	AssignEmployeeCompetencies(ctx context.Context, userProfileID uuid.UUID, competencyIDs []uuid.UUID) error
	GetCompetenciesByEmployee(ctx context.Context, userProfileID uuid.UUID) ([]entity.EmployeeCompetencyDetail, error)
	RemoveEmployeeCompetency(ctx context.Context, userProfileID uuid.UUID, competencyID uuid.UUID) error

	// Competencies (Dictionary)
	CreateCompetency(ctx context.Context, c *entity.Competency) (*entity.Competency, error)
	GetCompetencies(ctx context.Context, params params.QueryParams) ([]entity.Competency, int, error)
	GetCompetencyByID(ctx context.Context, id uuid.UUID) (*entity.Competency, error)
	UpdateCompetency(ctx context.Context, id uuid.UUID, c *entity.Competency) error
	DeleteCompetency(ctx context.Context, id uuid.UUID) error

	// System Settings
	GetSystemSetting(ctx context.Context, key string) (*entity.SystemSetting, error)
	SetSystemSetting(ctx context.Context, setting *entity.SystemSetting) error
	GetAllSystemSettings(ctx context.Context) ([]entity.SystemSetting, error)

	// Payroll Batch Calculation
	GetOrCreatePeriod(ctx context.Context, month, year int) (*entity.PayrollPeriod, error)
	GetPayrollPeriodByMonthYear(ctx context.Context, month, year int) (*entity.PayrollPeriod, error)
	UpsertPayrollRecord(ctx context.Context, record *entity.PayrollRecord, details []entity.PayrollRecordDetail) error
	GetPayrollRecords(ctx context.Context, periodID uuid.UUID) ([]entity.PayrollRecord, error)
	GetPayrollRecordDetails(ctx context.Context, recordID uuid.UUID, departmentID *uuid.UUID) ([]entity.PayrollRecordDetail, error)
	GetPayrollFormulasByPeriod(ctx context.Context, start, end time.Time) ([]entity.PayrollFormula, error)
	CreatePayrollFormula(ctx context.Context, formula *entity.PayrollFormula) (*entity.PayrollFormula, error)
	GetPayrollFormulas(ctx context.Context) ([]entity.PayrollFormula, error)
	UpdatePayrollFormula(ctx context.Context, id uuid.UUID, formula *entity.PayrollFormula) error
	DeletePayrollFormula(ctx context.Context, id uuid.UUID) error
}
