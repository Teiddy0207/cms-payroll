package pdf

import (
	"cal-salary/core/middleware"
	"cal-salary/modules/pdf/controller"
	"cal-salary/modules/pdf/router"
	"cal-salary/modules/pdf/service"
	"cal-salary/modules/pdf/watcher"
	"context"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	watchDir  = "./pdf_watch"  // Thư mục theo dõi file mới
	outputDir = "./pdf_output" // Thư mục lưu file kết quả
)

type PDFModule struct {
	Service service.PDFServiceInterface
	Watcher *watcher.FolderWatcher
}

func Init() *PDFModule {
	// Tạo thư mục watch và output nếu chưa tồn tại
	os.MkdirAll(watchDir, 0755)
	os.MkdirAll(outputDir, 0755)

	svc := service.NewPDFService()
	fw := watcher.NewFolderWatcher(watchDir, outputDir, svc)

	return &PDFModule{
		Service: svc,
		Watcher: fw,
	}
}

// SetupRouter đăng ký các route HTTP cho PDF Module
func (m *PDFModule) SetupRouter(e *echo.Echo, mw *middleware.Middleware) {
	ctrl := controller.NewPDFController(m.Service)
	router.NewPDFRouter(ctrl).Setup(e, mw)
}

// StartWatcher khởi động Folder Watcher trong background
func (m *PDFModule) StartWatcher(ctx context.Context) error {
	return m.Watcher.Start(ctx)
}
