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
	pdf.POST("/merge", r.ctrl.UploadAndMerge)       // Gộp nhiều PDF
	pdf.POST("/compress", r.ctrl.UploadAndCompress)  // Nén PDF
	pdf.POST("/watermark", r.ctrl.UploadAndWatermark) // Thêm watermark
	pdf.POST("/rotate", r.ctrl.UploadAndRotate)      // Xoay trang

	// Download kết quả
	pdf.GET("/download", r.ctrl.DownloadFile) // Download file output
}
