package service

import (
	"cal-salary/core/logger"
	"cal-salary/modules/pdf/dto"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

const outputBaseDir = "./pdf_output"

// PDFServiceInterface định nghĩa tất cả các tác vụ xử lý PDF
type PDFServiceInterface interface {
	MergePDFs(ctx context.Context, req dto.MergeRequest) (*dto.PDFJobResponse, error)
	SplitPDF(ctx context.Context, req dto.SplitRequest) (*dto.PDFJobResponse, error)
	CompressPDF(ctx context.Context, req dto.CompressRequest) (*dto.PDFJobResponse, error)
	WatermarkPDF(ctx context.Context, req dto.WatermarkRequest) (*dto.PDFJobResponse, error)
	RotatePDF(ctx context.Context, req dto.RotateRequest) (*dto.PDFJobResponse, error)
}

type PDFService struct{}

func NewPDFService() PDFServiceInterface {
	if err := os.MkdirAll(outputBaseDir, 0755); err != nil {
		logger.Error("PDFService: Cannot create output directory", "error", err)
	}
	return &PDFService{}
}

func generateOutputPath(name string) string {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	safeName := strings.ReplaceAll(name, " ", "_")
	return filepath.Join(outputBaseDir, fmt.Sprintf("%s_%s.pdf", safeName, timestamp))
}

// MergePDFs gộp nhiều file PDF thành một file duy nhất
func (s *PDFService) MergePDFs(ctx context.Context, req dto.MergeRequest) (*dto.PDFJobResponse, error) {
	outputPath := generateOutputPath(req.OutputName)
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	logger.Info("PDFService: Merging PDFs", "files", len(req.InputPaths), "output", outputPath)

	if err := pdfcpuapi.MergeCreateFile(req.InputPaths, outputPath, false, conf); err != nil {
		logger.Error("PDFService: Merge failed", "error", err)
		return nil, fmt.Errorf("không thể gộp PDF: %w", err)
	}

	logger.Info("PDFService: Merge completed", "output", outputPath)
	return &dto.PDFJobResponse{
		Success:    true,
		Message:    fmt.Sprintf("Đã gộp %d file PDF thành công", len(req.InputPaths)),
		OutputPath: outputPath,
		OutputURL:  "/api/v1/private/pdf/download?path=" + outputPath,
	}, nil
}

// SplitPDF tách file PDF thành nhiều file nhỏ theo số trang mỗi file
func (s *PDFService) SplitPDF(ctx context.Context, req dto.SplitRequest) (*dto.PDFJobResponse, error) {
	outDir := req.OutputDir
	if outDir == "" {
		outDir = outputBaseDir
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục output: %w", err)
	}

	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	logger.Info("PDFService: Splitting PDF", "input", req.InputPath, "span", req.SpanRange)

	if err := pdfcpuapi.SplitFile(req.InputPath, outDir, req.SpanRange, conf); err != nil {
		logger.Error("PDFService: Split failed", "error", err)
		return nil, fmt.Errorf("không thể tách PDF: %w", err)
	}

	logger.Info("PDFService: Split completed", "output_dir", outDir)
	return &dto.PDFJobResponse{
		Success:    true,
		Message:    fmt.Sprintf("Đã tách PDF thành các file %d trang/file, lưu vào: %s", req.SpanRange, outDir),
		OutputPath: outDir,
	}, nil
}

// CompressPDF tối ưu hóa và nén dung lượng file PDF
func (s *PDFService) CompressPDF(ctx context.Context, req dto.CompressRequest) (*dto.PDFJobResponse, error) {
	outputPath := generateOutputPath(req.OutputName)
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	conf.Cmd = model.OPTIMIZE

	logger.Info("PDFService: Compressing PDF", "input", req.InputPath)

	fInput, err := os.Open(req.InputPath)
	if err != nil {
		return nil, fmt.Errorf("không thể mở file PDF: %w", err)
	}
	defer fInput.Close()

	pdfCtx, err := pdfcpuapi.ReadContext(fInput, conf)
	if err != nil {
		logger.Error("PDFService: Compress ReadContext failed", "error", err)
		return nil, fmt.Errorf("không thể đọc PDF: %w", err)
	}

	pdfCtx.EnsureVersionForWriting()
	_ = pdfcpuapi.ValidateContext(pdfCtx)
	_ = pdfcpuapi.OptimizeContext(pdfCtx)

	if err := pdfcpuapi.WriteContextFile(pdfCtx, outputPath); err != nil {
		logger.Error("PDFService: Compress WriteContextFile failed", "error", err)
		return nil, fmt.Errorf("không thể nén PDF: %w", err)
	}

	srcInfo, _ := os.Stat(req.InputPath)
	dstInfo, _ := os.Stat(outputPath)
	var ratio float64
	if srcInfo != nil && srcInfo.Size() > 0 {
		ratio = float64(dstInfo.Size()) / float64(srcInfo.Size()) * 100
	}

	logger.Info("PDFService: Compress completed", "original_size", srcInfo.Size(), "compressed_size", dstInfo.Size(), "ratio", ratio)
	return &dto.PDFJobResponse{
		Success:    true,
		Message:    fmt.Sprintf("Nén PDF thành công. Kích thước còn %.1f%% so với bản gốc", ratio),
		OutputPath: outputPath,
		OutputURL:  "/api/v1/private/pdf/download?path=" + outputPath,
	}, nil
}

// WatermarkPDF thêm chữ watermark lên tất cả các trang PDF
func (s *PDFService) WatermarkPDF(ctx context.Context, req dto.WatermarkRequest) (*dto.PDFJobResponse, error) {
	outputPath := generateOutputPath(req.OutputName)

	opacity := req.Opacity
	if opacity <= 0 {
		opacity = 0.3
	}
	angle := req.Angle
	if angle == 0 {
		angle = 45
	}

	logger.Info("PDFService: Adding watermark", "text", req.Text, "opacity", opacity, "angle", angle)

	wmDesc := fmt.Sprintf("op:%.2f, rot:%.0f, scale:.7 abs, pos:c", opacity, angle)
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	conf.Cmd = model.ADDWATERMARKS

	wm, err := pdfcpuapi.TextWatermark(req.Text, wmDesc, true, false, types.POINTS)
	if err != nil {
		logger.Error("PDFService: Create watermark config failed", "error", err)
		return nil, fmt.Errorf("không thể tạo cấu hình watermark: %w", err)
	}

	fInput, err := os.Open(req.InputPath)
	if err != nil {
		return nil, fmt.Errorf("không thể mở file PDF đầu vào: %w", err)
	}
	defer fInput.Close()

	// 1. Đọc context PDF mà chưa chạy validator phiên bản strict
	pdfCtx, err := pdfcpuapi.ReadContext(fInput, conf)
	if err != nil {
		logger.Error("PDFService: ReadContext failed", "error", err)
		return nil, fmt.Errorf("không thể đọc cấu trúc PDF: %w", err)
	}

	// 2. Tự động nâng Version phiên bản xuất lên PDF 1.7 trước khi validate
	// Điều này giúp tránh lỗi "unsupported in version 1.4" khi PDF gốc chứa annotation PDF 1.6+
	pdfCtx.EnsureVersionForWriting()

	// 3. Chạy validate nhẹ (nếu có cảnh báo cũng không dừng)
	_ = pdfcpuapi.ValidateContext(pdfCtx)

	// 4. Lấy toàn bộ danh sách trang
	pages, err := pdfcpuapi.PagesForPageSelection(pdfCtx.PageCount, nil, true, true)
	if err != nil {
		return nil, fmt.Errorf("lỗi chọn trang PDF: %w", err)
	}

	// 5. Thêm watermark vào context
	if err := pdfcpu.AddWatermarks(pdfCtx, pages, wm); err != nil {
		logger.Error("PDFService: AddWatermarks failed", "error", err)
		return nil, fmt.Errorf("không thể thêm watermark vào PDF: %w", err)
	}

	// 6. Ghi file kết quả
	if err := pdfcpuapi.WriteContextFile(pdfCtx, outputPath); err != nil {
		logger.Error("PDFService: WriteContextFile failed", "error", err)
		return nil, fmt.Errorf("không thể ghi file PDF kết quả: %w", err)
	}

	logger.Info("PDFService: Watermark completed", "output", outputPath)
	return &dto.PDFJobResponse{
		Success:    true,
		Message:    "Đã thêm watermark vào PDF thành công",
		OutputPath: outputPath,
		OutputURL:  "/api/v1/private/pdf/download?path=" + outputPath,
	}, nil
}

// RotatePDF xoay các trang trong file PDF
func (s *PDFService) RotatePDF(ctx context.Context, req dto.RotateRequest) (*dto.PDFJobResponse, error) {
	outputPath := generateOutputPath(req.OutputName)
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	conf.Cmd = model.ROTATE

	var pageSelection []string
	if req.Pages != "" {
		pageSelection = strings.Split(req.Pages, ",")
	}

	logger.Info("PDFService: Rotating PDF", "angle", req.Angle, "pages", req.Pages)

	fInput, err := os.Open(req.InputPath)
	if err != nil {
		return nil, fmt.Errorf("không thể mở file PDF: %w", err)
	}
	defer fInput.Close()

	pdfCtx, err := pdfcpuapi.ReadContext(fInput, conf)
	if err != nil {
		logger.Error("PDFService: Rotate ReadContext failed", "error", err)
		return nil, fmt.Errorf("không thể đọc PDF: %w", err)
	}

	pdfCtx.EnsureVersionForWriting()
	_ = pdfcpuapi.ValidateContext(pdfCtx)

	pages, err := pdfcpuapi.PagesForPageSelection(pdfCtx.PageCount, pageSelection, true, true)
	if err != nil {
		return nil, fmt.Errorf("lỗi chọn trang: %w", err)
	}

	if err := pdfcpu.RotatePages(pdfCtx, pages, req.Angle); err != nil {
		logger.Error("PDFService: RotatePages failed", "error", err)
		return nil, fmt.Errorf("không thể xoay trang PDF: %w", err)
	}

	if err := pdfcpuapi.WriteContextFile(pdfCtx, outputPath); err != nil {
		logger.Error("PDFService: WriteContextFile failed", "error", err)
		return nil, fmt.Errorf("không thể ghi file PDF: %w", err)
	}

	pageDesc := "tất cả các trang"
	if req.Pages != "" {
		pageDesc = "trang " + req.Pages
	}

	logger.Info("PDFService: Rotate completed", "output", outputPath)
	return &dto.PDFJobResponse{
		Success:    true,
		Message:    fmt.Sprintf("Đã xoay %s thêm %d độ thành công", pageDesc, req.Angle),
		OutputPath: outputPath,
		OutputURL:  "/api/v1/private/pdf/download?path=" + outputPath,
	}, nil
}
