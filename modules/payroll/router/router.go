package router

import (
	"cal-salary/core/middleware"
	activitylog "cal-salary/modules/activity_log/service"
	"cal-salary/modules/payroll/controller"

	"github.com/labstack/echo/v4"
)

type PayrollRouter struct {
	PayrollController *controller.PayrollController
}

func NewPayrollRouter(payrollController *controller.PayrollController) *PayrollRouter {
	return &PayrollRouter{
		PayrollController: payrollController,
	}
}

func (r *PayrollRouter) Setup(e *echo.Echo, middlewareInstance *middleware.Middleware, activityLogSvc activitylog.ActivityLogServiceInterface) {
	v1 := e.Group("/api/v1")
	privateRoutes := v1.Group("/private")

	// Apply AuthMiddleware to all private routes in payroll
	privateRoutes.Use(middlewareInstance.AuthMiddleware())

	// Apply ActivityLogMiddleware
	payrollRoutes := privateRoutes.Group("/payroll")
	payrollRoutes.Use(middlewareInstance.ActivityLogMiddleware(activityLogSvc, "payroll"))

	// Departments endpoints
	payrollRoutes.POST("/departments", r.PayrollController.CreateDepartment)
	payrollRoutes.GET("/departments", r.PayrollController.GetDepartments)
	payrollRoutes.GET("/departments/:id", r.PayrollController.GetDepartmentByID)
	payrollRoutes.PUT("/departments/:id", r.PayrollController.UpdateDepartment)
	payrollRoutes.DELETE("/departments/:id", r.PayrollController.DeleteDepartment)

	// Job Positions endpoints
	payrollRoutes.POST("/job-positions", r.PayrollController.CreateJobPosition)
	payrollRoutes.GET("/job-positions", r.PayrollController.GetJobPositions)
	payrollRoutes.GET("/job-positions/:id", r.PayrollController.GetJobPositionByID)
	payrollRoutes.PUT("/job-positions/:id", r.PayrollController.UpdateJobPosition)
	payrollRoutes.DELETE("/job-positions/:id", r.PayrollController.DeleteJobPosition)
	payrollRoutes.POST("/user-profiles/:id/competencies", r.PayrollController.AssignEmployeeCompetencies)
	payrollRoutes.GET("/user-profiles/:id/competencies", r.PayrollController.GetCompetenciesByEmployee)
	payrollRoutes.DELETE("/user-profiles/:id/competencies/:competency_id", r.PayrollController.RemoveEmployeeCompetency)
	payrollRoutes.POST("/job-positions/:id/standards", r.PayrollController.AssignStandardToPosition)
	payrollRoutes.GET("/job-positions/:id/standards", r.PayrollController.GetStandardsByPosition)
	payrollRoutes.DELETE("/job-positions/:id/standards/:standard_id", r.PayrollController.RemoveStandardFromPosition)

	// User Profiles endpoints
	payrollRoutes.POST("/user-profiles", r.PayrollController.CreateUserProfile)
	payrollRoutes.GET("/user-profiles", r.PayrollController.GetUserProfiles)
	payrollRoutes.GET("/user-profiles/:id", r.PayrollController.GetUserProfileByID)
	payrollRoutes.PUT("/user-profiles/:id", r.PayrollController.UpdateUserProfile)
	payrollRoutes.DELETE("/user-profiles/:id", r.PayrollController.DeleteUserProfile)

	// Contracts endpoints
	payrollRoutes.POST("/contracts", r.PayrollController.CreateContract)
	payrollRoutes.GET("/contracts", r.PayrollController.GetContracts)
	payrollRoutes.GET("/contracts/:id", r.PayrollController.GetContractByID)
	payrollRoutes.PUT("/contracts/:id", r.PayrollController.UpdateContract)
	payrollRoutes.DELETE("/contracts/:id", r.PayrollController.DeleteContract)

	// Job Standards endpoints
	payrollRoutes.POST("/job-standards", r.PayrollController.CreateJobStandard)
	payrollRoutes.GET("/job-standards", r.PayrollController.GetJobStandards)
	payrollRoutes.GET("/job-standards/:id", r.PayrollController.GetJobStandardByID)
	payrollRoutes.PUT("/job-standards/:id", r.PayrollController.UpdateJobStandard)
	payrollRoutes.DELETE("/job-standards/:id", r.PayrollController.DeleteJobStandard)

	// Competency Evaluations endpoints
	payrollRoutes.POST("/evaluations", r.PayrollController.CreateBatchEvaluations)
	payrollRoutes.GET("/evaluations", r.PayrollController.GetEvaluations)
	payrollRoutes.GET("/evaluations/:id", r.PayrollController.GetEvaluationByID)
	payrollRoutes.PUT("/evaluations/:id", r.PayrollController.UpdateEvaluation)
	payrollRoutes.DELETE("/evaluations/:id", r.PayrollController.DeleteEvaluation)

	// Competency Dictionaries endpoints
	payrollRoutes.POST("/competencies", r.PayrollController.CreateCompetency)
	payrollRoutes.GET("/competencies", r.PayrollController.GetCompetencies)
	payrollRoutes.GET("/competencies/:id", r.PayrollController.GetCompetencyByID)
	payrollRoutes.PUT("/competencies/:id", r.PayrollController.UpdateCompetency)
	payrollRoutes.DELETE("/competencies/:id", r.PayrollController.DeleteCompetency)

	// System Settings endpoints
	payrollRoutes.GET("/settings", r.PayrollController.GetAllSystemSettings)
	payrollRoutes.GET("/settings/:key", r.PayrollController.GetSystemSetting)
	payrollRoutes.PUT("/settings/:key", r.PayrollController.UpdateSystemSetting)

	// Calculator Engine endpoints
	payrollRoutes.GET("/calculator/preview/:employee_id", r.PayrollController.PreviewSalary)
	payrollRoutes.POST("/calculator/run", r.PayrollController.RunSalaryCalculation)
	payrollRoutes.GET("/calculator/job/:job_id", r.PayrollController.GetCalculationJobStatus)
	payrollRoutes.GET("/calculator/records", r.PayrollController.GetSavedPayrollRecords)

	// Formulas CRUD
	payrollRoutes.POST("/calculator/formulas", r.PayrollController.CreatePayrollFormula)
	payrollRoutes.GET("/calculator/formulas", r.PayrollController.GetPayrollFormulas)
	payrollRoutes.PUT("/calculator/formulas/:id", r.PayrollController.UpdatePayrollFormula)
	payrollRoutes.DELETE("/calculator/formulas/:id", r.PayrollController.DeletePayrollFormula)
}
