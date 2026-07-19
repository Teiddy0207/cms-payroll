package watcher

import (
	"cal-salary/core/logger"
	"cal-salary/modules/pdf/dto"
	"cal-salary/modules/pdf/service"
	"context"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// FolderWatcher theo dõi thư mục và tự động xử lý file PDF mới xuất hiện
type FolderWatcher struct {
	watchDir   string
	outputDir  string
	pdfService service.PDFServiceInterface
}

func NewFolderWatcher(watchDir, outputDir string, svc service.PDFServiceInterface) *FolderWatcher {
	return &FolderWatcher{
		watchDir:   watchDir,
		outputDir:  outputDir,
		pdfService: svc,
	}
}

// Start bắt đầu theo dõi thư mục trong nền (non-blocking, chạy trong goroutine)
func (w *FolderWatcher) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		logger.Info("FolderWatcher: Started watching directory", "dir", w.watchDir)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Chỉ xử lý khi file PDF mới được tạo hoặc copy vào thư mục
				if (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) &&
					filepath.Ext(event.Name) == ".pdf" {
					logger.Info("FolderWatcher: Detected new PDF file", "file", event.Name)
					// Xử lý bất đồng bộ để không block watcher
					go w.autoPipeline(ctx, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("FolderWatcher: Error watching directory", "error", err)
			case <-ctx.Done():
				logger.Info("FolderWatcher: Stopping due to context cancellation")
				return
			}
		}
	}()

	return watcher.Add(w.watchDir)
}

// autoPipeline thực hiện pipeline tự động: Watermark → Compress → ghi ra thư mục output
func (w *FolderWatcher) autoPipeline(ctx context.Context, filePath string) {
	baseName := filepath.Base(filePath)
	nameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]

	logger.Info("FolderWatcher: Running auto pipeline", "file", filePath)

	// Bước 1: Thêm watermark "BẢO MẬT - NỘI BỘ"
	wmResp, err := w.pdfService.WatermarkPDF(ctx, dto.WatermarkRequest{
		InputPath:  filePath,
		OutputName: nameWithoutExt + "_wm",
		Text:       "BẢO MẬT - NỘI BỘ",
		Opacity:    0.25,
		Angle:      45,
	})
	if err != nil {
		logger.Error("FolderWatcher: Watermark step failed", "error", err)
		return
	}

	// Bước 2: Nén file vừa đóng dấu watermark
	_, err = w.pdfService.CompressPDF(ctx, dto.CompressRequest{
		InputPath:  wmResp.OutputPath,
		OutputName: nameWithoutExt + "_final",
	})
	if err != nil {
		logger.Error("FolderWatcher: Compress step failed", "error", err)
		return
	}

	logger.Info("FolderWatcher: Auto pipeline completed successfully", "file", filePath)
}
