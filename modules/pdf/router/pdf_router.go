package router

import (
	"cal-salary/core/middleware"
	"cal-salary/modules/pdf/controller"

	"github.com/labstack/echo/v4"
)

type PDFRouter struct {
	ctrl *controller.PDFController
}

func NewPDFRouter(ctrl *controller.PDFController) *PDFRouter {
	return &PDFRouter{ctrl: ctrl}
}

func (r *PDFRouter) Setup(e *echo.Echo, mw *middleware.Middleware) {
	// Nhóm API PDF - yêu cầu xác thực JWT
	pdf := e.Group("/api/v1/private/pdf", mw.AuthMiddleware())

	// Xử lý tệp PDF
	pdf.POST("/merge", r.ctrl.UploadAndMerge, mw.PermissionMiddleware("storage::edit"))
	pdf.POST("/compress", r.ctrl.UploadAndCompress, mw.PermissionMiddleware("storage::edit"))
	pdf.POST("/watermark", r.ctrl.UploadAndWatermark, mw.PermissionMiddleware("storage::edit"))
	pdf.POST("/rotate", r.ctrl.UploadAndRotate, mw.PermissionMiddleware("storage::edit"))

	// Download kết quả
	pdf.GET("/download", r.ctrl.DownloadFile, mw.PermissionMiddleware("storage::edit"))
}
