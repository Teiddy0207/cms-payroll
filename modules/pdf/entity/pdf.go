package entity

import "time"

// PDFJob đại diện cho một công việc xử lý PDF
type PDFJob struct {
	ID          string    `db:"id" json:"id"`
	FileName    string    `db:"file_name" json:"file_name"`
	FilePath    string    `db:"file_path" json:"file_path"`
	Operation   string    `db:"operation" json:"operation"` // MERGE, SPLIT, COMPRESS, WATERMARK, ROTATE
	Status      string    `db:"status" json:"status"`       // PENDING, PROCESSING, COMPLETED, FAILED
	ErrorMsg    string    `db:"error_msg" json:"error_msg,omitempty"`
	OutputPath  string    `db:"output_path" json:"output_path,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// WatermarkConfig cấu hình thêm watermark vào PDF
type WatermarkConfig struct {
	Text     string  `json:"text"`      // Nội dung watermark
	Opacity  float64 `json:"opacity"`   // Độ trong suốt (0.0 - 1.0)
	Angle    float64 `json:"angle"`     // Góc nghiêng (0-360)
	FontSize int     `json:"font_size"` // Cỡ chữ
}

// RotateConfig cấu hình xoay trang PDF
type RotateConfig struct {
	Pages []int `json:"pages"`  // Danh sách số trang cần xoay (1-based)
	Angle int   `json:"angle"`  // Góc xoay: 90, 180, 270
}
