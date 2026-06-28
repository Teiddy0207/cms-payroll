package controller

import (
	"cal-salary/core/controller"
	"cal-salary/modules/activity_log/service"
)

type ActivityLogController struct {
	controller.BaseController
	ActivityLogService service.ActivityLogServiceInterface
}

func NewActivityLogController(activityLogService service.ActivityLogServiceInterface) *ActivityLogController {
	return &ActivityLogController{
		BaseController:     controller.NewBaseController(),
		ActivityLogService: activityLogService,
	}
}
