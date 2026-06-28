package controller

import (
	"cal-salary/core/controller"
	"cal-salary/modules/payroll/service"
)

type PayrollController struct {
	controller.BaseController
	Service service.PayrollServiceInterface
}

func NewPayrollController(svc service.PayrollServiceInterface) *PayrollController {
	return &PayrollController{
		BaseController: controller.NewBaseController(),
		Service:        svc,
	}
}
