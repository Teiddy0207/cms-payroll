package service

import (
	"cal-salary/core/cache"
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/repository"
	"context"

	"github.com/google/uuid"
)

type PayrollService struct {
	repo  repository.PayrollRepositoryInterface
	cache *cache.Cache
}

func NewPayrollService(repo repository.PayrollRepositoryInterface, redisCache *cache.Cache) *PayrollService {
	return &PayrollService{repo: repo, cache: redisCache}
}

type PayrollServiceInterface interface {
	// Departments
	CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, *errors.AppError)
	GetDepartments(ctx context.Context, params params.QueryParams) (*dto.PaginatedDepartmentDTO, *errors.AppError)
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*dto.DepartmentResponse, *errors.AppError)
	UpdateDepartment(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartmentRequest) *errors.AppError
	DeleteDepartment(ctx context.Context, id uuid.UUID) *errors.AppError

	// Job Positions
	CreateJobPosition(ctx context.Context, req *dto.CreateJobPositionRequest) (*dto.JobPositionResponse, *errors.AppError)
	GetJobPositions(ctx context.Context, params params.QueryParams) (*dto.PaginatedJobPositionDTO, *errors.AppError)
	GetJobPositionByID(ctx context.Context, id uuid.UUID) (*dto.JobPositionResponse, *errors.AppError)
	UpdateJobPosition(ctx context.Context, id uuid.UUID, req *dto.UpdateJobPositionRequest) *errors.AppError
	DeleteJobPosition(ctx context.Context, id uuid.UUID) *errors.AppError

	// User Profiles
	CreateUserProfile(ctx context.Context, req *dto.CreateUserProfileRequest) (*dto.UserProfileResponse, *errors.AppError)
	GetUserProfiles(ctx context.Context, params params.QueryParams) (*dto.PaginatedUserProfileDTO, *errors.AppError)
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, *errors.AppError)
	UpdateUserProfile(ctx context.Context, id uuid.UUID, req *dto.UpdateUserProfileRequest) *errors.AppError
	DeleteUserProfile(ctx context.Context, id uuid.UUID) *errors.AppError

	// Contracts
	CreateContract(ctx context.Context, req *dto.CreateContractRequest) (*dto.ContractResponse, *errors.AppError)
	GetContracts(ctx context.Context, params params.QueryParams) (*dto.PaginatedContractDTO, *errors.AppError)
	GetContractByID(ctx context.Context, id uuid.UUID) (*dto.ContractResponse, *errors.AppError)
	UpdateContract(ctx context.Context, id uuid.UUID, req *dto.UpdateContractRequest) *errors.AppError
	DeleteContract(ctx context.Context, id uuid.UUID) *errors.AppError

	// Job Standards
	CreateJobStandard(ctx context.Context, req *dto.CreateJobStandardRequest) (*dto.JobStandardResponse, *errors.AppError)
	GetJobStandards(ctx context.Context, params params.QueryParams) (*dto.PaginatedJobStandardDTO, *errors.AppError)
	GetJobStandardByID(ctx context.Context, id uuid.UUID) (*dto.JobStandardResponse, *errors.AppError)
	UpdateJobStandard(ctx context.Context, id uuid.UUID, req *dto.UpdateJobStandardRequest) *errors.AppError
	DeleteJobStandard(ctx context.Context, id uuid.UUID) *errors.AppError

	// Competencies (Dictionary)
	CreateCompetency(ctx context.Context, req *dto.CreateCompetencyRequest) (*dto.CompetencyResponse, *errors.AppError)
	GetCompetencies(ctx context.Context, params params.QueryParams) (*dto.PaginatedCompetencyDTO, *errors.AppError)
	GetCompetencyByID(ctx context.Context, id uuid.UUID) (*dto.CompetencyResponse, *errors.AppError)
	UpdateCompetency(ctx context.Context, id uuid.UUID, req *dto.UpdateCompetencyRequest) *errors.AppError
	DeleteCompetency(ctx context.Context, id uuid.UUID) *errors.AppError

	// Competency Evaluations
	CreateBatchEvaluations(ctx context.Context, evaluatorID uuid.UUID, req *dto.BatchEvaluationRequest) ([]dto.EvaluationResponse, *errors.AppError)
	GetEvaluations(ctx context.Context, params params.QueryParams) (*dto.PaginatedEvaluationDTO, *errors.AppError)
	GetEvaluationByID(ctx context.Context, id uuid.UUID) (*dto.EvaluationResponse, *errors.AppError)
	UpdateEvaluation(ctx context.Context, id uuid.UUID, req *dto.CompetencyScoreInput) *errors.AppError
	DeleteEvaluation(ctx context.Context, id uuid.UUID) *errors.AppError

	// Job Position Standards
	AssignStandardToPosition(ctx context.Context, positionID uuid.UUID, req *dto.AssignStandardRequest) *errors.AppError
	GetStandardsByPosition(ctx context.Context, positionID uuid.UUID) ([]dto.JobPositionStandardResponse, *errors.AppError)
	RemoveStandardFromPosition(ctx context.Context, positionID uuid.UUID, standardID uuid.UUID) *errors.AppError

	// Employee Competencies
	AssignEmployeeCompetencies(ctx context.Context, userProfileID uuid.UUID, req *dto.AssignEmployeeCompetenciesRequest) *errors.AppError
	GetCompetenciesByEmployee(ctx context.Context, userProfileID uuid.UUID) ([]dto.EmployeeCompetencyResponse, *errors.AppError)
	RemoveEmployeeCompetency(ctx context.Context, userProfileID uuid.UUID, competencyID uuid.UUID) *errors.AppError

	// System Settings
	GetSystemSetting(ctx context.Context, key string) (*dto.SystemSettingResponse, *errors.AppError)
	UpdateSystemSetting(ctx context.Context, key string, req *dto.UpdateSystemSettingRequest) *errors.AppError
	GetAllSystemSettings(ctx context.Context) ([]dto.SystemSettingResponse, *errors.AppError)

	// Calculator
	PreviewSalary(ctx context.Context, employeeID uuid.UUID, period string) (*dto.SalaryPreviewResponse, *errors.AppError)
	RunSalaryCalculation(ctx context.Context, req *dto.SalaryCalculationRequest) (*dto.SalaryCalculationResponse, *errors.AppError)
	RunSalaryCalculationAsync(ctx context.Context, req *dto.SalaryCalculationRequest) (string, *errors.AppError)
	GetCalculationJobStatus(ctx context.Context, jobID string) (map[string]any, *errors.AppError)
	GetSavedPayrollRecords(ctx context.Context, period string) ([]dto.SalaryCalculationItem, *errors.AppError)

	// Formulas CRUD
	CreatePayrollFormula(ctx context.Context, req *dto.CreateFormulaRequest) (*dto.FormulaResponse, *errors.AppError)
	GetPayrollFormulas(ctx context.Context) ([]dto.FormulaResponse, *errors.AppError)
	UpdatePayrollFormula(ctx context.Context, id uuid.UUID, req *dto.UpdateFormulaRequest) *errors.AppError
	DeletePayrollFormula(ctx context.Context, id uuid.UUID) *errors.AppError
}
