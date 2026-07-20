package dto

// MergeRequest yêu cầu gộp nhiều file PDF thành 1
type MergeRequest struct {
	// InputPaths là danh sách đường dẫn tuyệt đối đến các file PDF cần gộp
	// Thứ tự trong slice là thứ tự trang trong file kết quả
	InputPaths []string `json:"input_paths" validate:"required,min=2"`
	OutputName string   `json:"output_name" validate:"required"` // Tên file kết quả (không có extension)
}

// SplitRequest yêu cầu tách file PDF thành nhiều file theo trang
type SplitRequest struct {
	InputPath string `json:"input_path" validate:"required"`
	// SpanRange: số trang mỗi file output. Ví dụ span=3 -> mỗi file có 3 trang
	SpanRange int    `json:"span" validate:"required,min=1"`
	OutputDir string `json:"output_dir"`
}

// CompressRequest yêu cầu nén file PDF
type CompressRequest struct {
	InputPath  string `json:"input_path" validate:"required"`
	OutputName string `json:"output_name" validate:"required"`
}

// WatermarkRequest yêu cầu thêm watermark vào PDF
type WatermarkRequest struct {
	InputPath  string  `json:"input_path" validate:"required"`
	OutputName string  `json:"output_name" validate:"required"`
	Text       string  `json:"text" validate:"required"`   // Nội dung chữ watermark
	Opacity    float64 `json:"opacity"`                    // Mặc định 0.3
	Angle      float64 `json:"angle"`                      // Mặc định 45 độ
}

// RotateRequest yêu cầu xoay trang trong PDF
type RotateRequest struct {
	InputPath  string `json:"input_path" validate:"required"`
	OutputName string `json:"output_name" validate:"required"`
	Angle      int    `json:"angle" validate:"required"`      // 90, 180, 270
	// Pages: để trống = xoay tất cả trang. Ví dụ: "1-3,5" = xoay trang 1,2,3,5
	Pages      string `json:"pages"`
}

// PDFJobResponse phản hồi về trạng thái một job xử lý PDF
type PDFJobResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	OutputPath string `json:"output_path,omitempty"`
	OutputURL  string `json:"output_url,omitempty"` // URL download kết quả
}
