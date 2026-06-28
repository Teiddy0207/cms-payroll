package controller

import (
	"cal-salary/core/controller"
	"cal-salary/modules/auth/repository"
	"cal-salary/modules/auth/service"
)

type AuthController struct {
	controller.BaseController
	AuthService service.AuthServiceInterface
	AuthRepo    repository.AuthRepositoryInterface
	// PayrollRepo  payrollRepo.PayrollRepositoryInterface
}

func NewAuthController(service service.AuthServiceInterface, repo repository.AuthRepositoryInterface) *AuthController {
	return &AuthController{
		BaseController: controller.NewBaseController(),
		AuthService:    service,
		AuthRepo:       repo,
		// PayrollRepo:    payrollRepo,
	}
}
