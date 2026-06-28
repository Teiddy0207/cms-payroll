package repository

import (
	"cal-salary/core/database"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/entity"
	"context"

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

	// Job Standards
	CreateJobStandard(ctx context.Context, standard *entity.JobStandard) (*entity.JobStandard, error)
	GetJobStandards(ctx context.Context, params params.QueryParams) ([]entity.JobStandard, int, error)
	GetJobStandardByID(ctx context.Context, id uuid.UUID) (*entity.JobStandard, error)
	UpdateJobStandard(ctx context.Context, id uuid.UUID, standard *entity.JobStandard) error
	DeleteJobStandard(ctx context.Context, id uuid.UUID) error

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
}
