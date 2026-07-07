package controller

import (
	"cal-salary/core/constants"
	"cal-salary/core/controller"
	"cal-salary/core/errors"
	"cal-salary/core/utils"
	"cal-salary/modules/timekeeping/dto"
	"cal-salary/modules/timekeeping/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TimekeepingController struct {
	controller.BaseController
	Service service.TimekeepingService
}

func NewTimekeepingController(svc service.TimekeepingService) *TimekeepingController {
	return &TimekeepingController{
		BaseController: controller.NewBaseController(),
		Service:        svc,
	}
}

func (ctrl *TimekeepingController) getUserClaims(c echo.Context) (*utils.TokenClaims, error) {
	userData := c.Get(constants.ContextTokenData)
	if userData == nil {
		return nil, echo.ErrUnauthorized
	}
	claims, ok := userData.(*utils.TokenClaims)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	return claims, nil
}

func (ctrl *TimekeepingController) ProcessCheckIn(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CheckinRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu check-in không hợp lệ", nil)
	}

	resp, appErr := ctrl.Service.ProcessCheckIn(ctx, req)
	if appErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"code": appErr.Code, "message": appErr.Message})
	}
	return ctrl.SuccessResponse(c, resp, "Ghi nhận check-in thành công")
}

func (ctrl *TimekeepingController) GetDailyAttendanceSheets(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	period := c.QueryParam("period")
	if period == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Thiếu tham số chu kỳ period (YYYY-MM)", nil)
	}

	resp, appErr := ctrl.Service.GetDailyAttendanceSheets(ctx, claims.UserID, period)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Tải bảng công thành công")
}

func (ctrl *TimekeepingController) CalculateTimesheets(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.CalculateTimesheetRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu yêu cầu không hợp lệ", nil)
	}

	if req.StartDate == "" || req.EndDate == "" {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Ngày bắt đầu và ngày kết thúc không được trống", nil)
	}

	appErr := ctrl.Service.CalculateTimesheets(ctx, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Tính toán bảng công thành công")
}

func (ctrl *TimekeepingController) CreateExplanationRequest(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	req := new(dto.CreateExplanationRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu yêu cầu không hợp lệ", nil)
	}

	resp, appErr := ctrl.Service.CreateExplanationRequest(ctx, claims.UserID, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Tạo đơn giải trình thành công")
}

func (ctrl *TimekeepingController) GetExplanationRequests(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	resp, appErr := ctrl.Service.GetExplanationRequests(ctx, claims.UserID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Tải danh sách đơn giải trình thành công")
}

func (ctrl *TimekeepingController) UpdateExplanationRequestStatus(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	reqID, errParse := uuid.Parse(c.Param("id"))
	if errParse != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "ID đơn không hợp lệ", nil)
	}

	req := new(dto.UpdateRequestStatus)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu yêu cầu không hợp lệ", nil)
	}

	appErr := ctrl.Service.UpdateExplanationRequestStatus(ctx, claims.UserID, reqID, req)
	if appErr != nil {
		if appErr.Code == errors.ErrForbidden {
			return ctrl.Forbidden(appErr.Code, appErr.Message, nil)
		}
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Phê duyệt đơn giải trình thành công")
}

func (ctrl *TimekeepingController) CreateOTRequest(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	req := new(dto.CreateOTRequest)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu yêu cầu không hợp lệ", nil)
	}

	resp, appErr := ctrl.Service.CreateOTRequest(ctx, claims.UserID, req)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Tạo đơn đăng ký OT thành công")
}

func (ctrl *TimekeepingController) GetOTRequests(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	resp, appErr := ctrl.Service.GetOTRequests(ctx, claims.UserID)
	if appErr != nil {
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, resp, "Tải danh sách đăng ký OT thành công")
}

func (ctrl *TimekeepingController) UpdateOTRequestStatus(c echo.Context) error {
	ctx := c.Request().Context()
	claims, err := ctrl.getUserClaims(c)
	if err != nil {
		return ctrl.Unauthorized(errors.ErrUnauthorized, "Yêu cầu đăng nhập", nil)
	}

	reqID, errParse := uuid.Parse(c.Param("id"))
	if errParse != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "ID đơn không hợp lệ", nil)
	}

	req := new(dto.UpdateRequestStatus)
	if err := c.Bind(req); err != nil {
		return ctrl.BadRequest(errors.ErrInvalidRequestData, "Dữ liệu yêu cầu không hợp lệ", nil)
	}

	appErr := ctrl.Service.UpdateOTRequestStatus(ctx, claims.UserID, reqID, req)
	if appErr != nil {
		if appErr.Code == errors.ErrForbidden {
			return ctrl.Forbidden(appErr.Code, appErr.Message, nil)
		}
		if appErr.Code == errors.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"code": appErr.Code, "message": appErr.Message})
		}
		return ctrl.InternalServerError(appErr.Code, appErr.Message, appErr)
	}
	return ctrl.SuccessResponse(c, nil, "Phê duyệt đơn đăng ký OT thành công")
}
