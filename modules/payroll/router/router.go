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
	payrollRoutes.POST("/departments", r.PayrollController.CreateDepartment, middlewareInstance.PermissionMiddleware("department::edit"))
	payrollRoutes.GET("/departments", r.PayrollController.GetDepartments, middlewareInstance.PermissionMiddleware("department::read"))
	payrollRoutes.GET("/departments/:id", r.PayrollController.GetDepartmentByID, middlewareInstance.PermissionMiddleware("department::read"))
	payrollRoutes.PUT("/departments/:id", r.PayrollController.UpdateDepartment, middlewareInstance.PermissionMiddleware("department::edit"))
	payrollRoutes.DELETE("/departments/:id", r.PayrollController.DeleteDepartment, middlewareInstance.PermissionMiddleware("department::delete"))

	// Job Positions endpoints
	payrollRoutes.POST("/job-positions", r.PayrollController.CreateJobPosition, middlewareInstance.PermissionMiddleware("jobPosition::edit"))
	payrollRoutes.GET("/job-positions", r.PayrollController.GetJobPositions, middlewareInstance.PermissionMiddleware("jobPosition::read"))
	payrollRoutes.POST("/job-positions/calculate-k", r.PayrollController.CalculateKAndUpdate, middlewareInstance.PermissionMiddleware("jobPosition::edit"))
	payrollRoutes.POST("/job-positions/:id/fetch-market-salary", r.PayrollController.ScrapeMarketSalary, middlewareInstance.PermissionMiddleware("jobPosition::edit"))
	payrollRoutes.POST("/job-positions/fetch-market-salary", r.PayrollController.ScrapeMarketSalary, middlewareInstance.PermissionMiddleware("jobPosition::edit"))
	payrollRoutes.GET("/job-positions/:id", r.PayrollController.GetJobPositionByID, middlewareInstance.PermissionMiddleware("jobPosition::read"))
	payrollRoutes.PUT("/job-positions/:id", r.PayrollController.UpdateJobPosition, middlewareInstance.PermissionMiddleware("jobPosition::edit"))
	payrollRoutes.DELETE("/job-positions/:id", r.PayrollController.DeleteJobPosition, middlewareInstance.PermissionMiddleware("jobPosition::delete"))
	payrollRoutes.POST("/user-profiles/:id/competencies", r.PayrollController.AssignEmployeeCompetencies, middlewareInstance.PermissionMiddleware("userProfileCompetency::edit"))
	payrollRoutes.GET("/user-profiles/:id/competencies", r.PayrollController.GetCompetenciesByEmployee, middlewareInstance.PermissionMiddleware("userProfileCompetency::read"))
	payrollRoutes.DELETE("/user-profiles/:id/competencies/:competency_id", r.PayrollController.RemoveEmployeeCompetency, middlewareInstance.PermissionMiddleware("userProfileCompetency::delete"))

	// User Profiles endpoints
	payrollRoutes.POST("/user-profiles", r.PayrollController.CreateUserProfile, middlewareInstance.PermissionMiddleware("userProfile::edit"))
	payrollRoutes.GET("/user-profiles", r.PayrollController.GetUserProfiles, middlewareInstance.PermissionMiddleware("userProfile::read"))
	payrollRoutes.GET("/user-profiles/:id", r.PayrollController.GetUserProfileByID, middlewareInstance.PermissionMiddleware("userProfile::read"))
	payrollRoutes.PUT("/user-profiles/:id", r.PayrollController.UpdateUserProfile, middlewareInstance.PermissionMiddleware("userProfile::edit"))
	payrollRoutes.DELETE("/user-profiles/:id", r.PayrollController.DeleteUserProfile, middlewareInstance.PermissionMiddleware("userProfile::delete"))

	// Contracts endpoints (uses userProfile for contract authorization)
	payrollRoutes.POST("/contracts", r.PayrollController.CreateContract, middlewareInstance.PermissionMiddleware("userProfile::edit"))
	payrollRoutes.GET("/contracts", r.PayrollController.GetContracts, middlewareInstance.PermissionMiddleware("userProfile::read"))
	payrollRoutes.GET("/contracts/:id", r.PayrollController.GetContractByID, middlewareInstance.PermissionMiddleware("userProfile::read"))
	payrollRoutes.PUT("/contracts/:id", r.PayrollController.UpdateContract, middlewareInstance.PermissionMiddleware("userProfile::edit"))
	payrollRoutes.DELETE("/contracts/:id", r.PayrollController.DeleteContract, middlewareInstance.PermissionMiddleware("userProfile::delete"))

	// Competency Evaluations endpoints
	payrollRoutes.POST("/evaluations", r.PayrollController.CreateBatchEvaluations, middlewareInstance.PermissionMiddleware("userProfileCompetency::edit"))
	payrollRoutes.GET("/evaluations", r.PayrollController.GetEvaluations, middlewareInstance.PermissionMiddleware("userProfileCompetency::read"))
	payrollRoutes.GET("/evaluations/:id", r.PayrollController.GetEvaluationByID, middlewareInstance.PermissionMiddleware("userProfileCompetency::read"))
	payrollRoutes.PUT("/evaluations/:id", r.PayrollController.UpdateEvaluation, middlewareInstance.PermissionMiddleware("userProfileCompetency::edit"))
	payrollRoutes.DELETE("/evaluations/:id", r.PayrollController.DeleteEvaluation, middlewareInstance.PermissionMiddleware("userProfileCompetency::delete"))

	// Competency Dictionaries endpoints
	payrollRoutes.POST("/competencies", r.PayrollController.CreateCompetency, middlewareInstance.PermissionMiddleware("jobCapability::edit"))
	payrollRoutes.GET("/competencies", r.PayrollController.GetCompetencies, middlewareInstance.PermissionMiddleware("jobCapability::read"))
	payrollRoutes.GET("/competencies/:id", r.PayrollController.GetCompetencyByID, middlewareInstance.PermissionMiddleware("jobCapability::read"))
	payrollRoutes.PUT("/competencies/:id", r.PayrollController.UpdateCompetency, middlewareInstance.PermissionMiddleware("jobCapability::edit"))
	payrollRoutes.DELETE("/competencies/:id", r.PayrollController.DeleteCompetency, middlewareInstance.PermissionMiddleware("jobCapability::delete"))

	// System Settings endpoints
	payrollRoutes.GET("/settings", r.PayrollController.GetAllSystemSettings, middlewareInstance.PermissionMiddleware("setting::read"))
	payrollRoutes.GET("/settings/:key", r.PayrollController.GetSystemSetting, middlewareInstance.PermissionMiddleware("setting::read"))
	payrollRoutes.PUT("/settings/:key", r.PayrollController.UpdateSystemSetting, middlewareInstance.PermissionMiddleware("setting::edit"))

	// Calculator Engine endpoints
	payrollRoutes.GET("/calculator/preview/:employee_id", r.PayrollController.PreviewSalary, middlewareInstance.PermissionMiddleware("salary::read"))
	payrollRoutes.POST("/calculator/run", r.PayrollController.RunSalaryCalculation, middlewareInstance.PermissionMiddleware("salary::edit"))
	payrollRoutes.GET("/calculator/job/:job_id", r.PayrollController.GetCalculationJobStatus, middlewareInstance.PermissionMiddleware("salary::read"))
	payrollRoutes.GET("/calculator/records", r.PayrollController.GetSavedPayrollRecords, middlewareInstance.PermissionMiddleware("salary::read"))

	// Formulas CRUD
	payrollRoutes.POST("/calculator/formulas", r.PayrollController.CreatePayrollFormula, middlewareInstance.PermissionMiddleware("formulaDynamic::edit"))
	payrollRoutes.GET("/calculator/formulas", r.PayrollController.GetPayrollFormulas, middlewareInstance.PermissionMiddleware("formulaDynamic::read"))
	payrollRoutes.PUT("/calculator/formulas/:id", r.PayrollController.UpdatePayrollFormula, middlewareInstance.PermissionMiddleware("formulaDynamic::edit"))
	payrollRoutes.DELETE("/calculator/formulas/:id", r.PayrollController.DeletePayrollFormula, middlewareInstance.PermissionMiddleware("formulaDynamic::delete"))
}
