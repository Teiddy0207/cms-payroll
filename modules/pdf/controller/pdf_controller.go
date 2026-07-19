package controller

import (
	"cal-salary/modules/pdf/dto"
	"cal-salary/modules/pdf/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type PDFController struct {
	Service service.PDFServiceInterface
}

func NewPDFController(svc service.PDFServiceInterface) *PDFController {
	return &PDFController{Service: svc}
}

// saveTempFile lưu multipart file ra đĩa, đóng file đúng cách và trả về đường dẫn
func saveTempFile(src io.ReadCloser, prefix string) (string, error) {
	defer src.Close()

	// Tạo file temp với tên an toàn (không dùng tên gốc để tránh ký tự đặc biệt)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("pdf_%s_%s.pdf", prefix, timestamp))

	dst, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo file tạm: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("lỗi ghi file tạm: %w", err)
	}

	// Đảm bảo flush xuống đĩa trước khi trả về path
	if err := dst.Sync(); err != nil {
		return "", fmt.Errorf("lỗi sync file tạm: %w", err)
	}

	return tmpPath, nil
}

// UploadAndMerge nhận nhiều file PDF upload lên, gộp lại và trả về file kết quả để download
func (ctrl *PDFController) UploadAndMerge(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Không thể đọc form data: " + err.Error()})
	}

	files := form.File["files"]
	if len(files) < 2 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Cần ít nhất 2 file PDF để gộp"})
	}

	outputName := c.FormValue("output_name")
	if outputName == "" {
		outputName = "merged_output"
	}

	var inputPaths []string
	for i, fh := range files {
		src, err := fh.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Lỗi đọc file: " + err.Error()})
		}
		// saveTempFile tự đóng src và dst sau khi ghi xong
		tmpPath, err := saveTempFile(src, fmt.Sprintf("merge_%d", i))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		inputPaths = append(inputPaths, tmpPath)
	}

	resp, err := ctrl.Service.MergePDFs(c.Request().Context(), dto.MergeRequest{
		InputPaths: inputPaths,
		OutputName: outputName,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// UploadAndCompress nhận 1 file PDF, nén và trả về
func (ctrl *PDFController) UploadAndCompress(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Không tìm thấy file trong request"})
	}

	src, err := fh.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	tmpPath, err := saveTempFile(src, "compress")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	outputName := strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename)) + "_compressed"
	resp, err := ctrl.Service.CompressPDF(c.Request().Context(), dto.CompressRequest{
		InputPath:  tmpPath,
		OutputName: outputName,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// UploadAndWatermark nhận 1 file PDF, thêm watermark và trả về
func (ctrl *PDFController) UploadAndWatermark(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Không tìm thấy file trong request"})
	}

	src, err := fh.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	tmpPath, err := saveTempFile(src, "watermark")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	watermarkText := c.FormValue("text")
	if watermarkText == "" {
		watermarkText = "BẢO MẬT"
	}
	outputName := strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename)) + "_watermarked"

	resp, err := ctrl.Service.WatermarkPDF(c.Request().Context(), dto.WatermarkRequest{
		InputPath:  tmpPath,
		OutputName: outputName,
		Text:       watermarkText,
		Opacity:    0.3,
		Angle:      45,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// UploadAndRotate nhận 1 file PDF, xoay trang và trả về
func (ctrl *PDFController) UploadAndRotate(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Không tìm thấy file trong request"})
	}

	src, err := fh.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	tmpPath, err := saveTempFile(src, "rotate")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	angle := 90
	if a := c.FormValue("angle"); a != "" {
		fmt.Sscanf(a, "%d", &angle)
	}
	pages := c.FormValue("pages")
	outputName := strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename)) + "_rotated"

	resp, err := ctrl.Service.RotatePDF(c.Request().Context(), dto.RotateRequest{
		InputPath:  tmpPath,
		OutputName: outputName,
		Angle:      angle,
		Pages:      pages,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// DownloadFile phục vụ download file kết quả xử lý PDF
func (ctrl *PDFController) DownloadFile(c echo.Context) error {
	filePath := c.QueryParam("path")
	if filePath == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Tham số path là bắt buộc"})
	}

	// Bảo mật: chỉ cho phép download trong thư mục pdf_output
	cleanPath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanPath, "pdf_output") {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Không được phép truy cập file này"})
	}

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "File không tồn tại"})
	}

	return c.File(cleanPath)
}
