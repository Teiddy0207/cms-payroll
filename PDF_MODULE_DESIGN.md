# ĐẶC TẢ THIẾT KẾ KỸ THUẬT: MODULE XỬ LÝ PDF (GOLANG)

Tài liệu này đặc tả thiết kế và kiến trúc cho Module xử lý PDF chuyên sâu tích hợp trực tiếp vào hệ thống Backend sử dụng ngôn ngữ Go.

---

## 1. Các Tính Năng Chức Năng (Features)

| Tính năng | Mô tả chi tiết | Giải pháp đề xuất trên Go |
| :--- | :--- | :--- |
| **Organize Pages** | Ghép file (Merge), chia nhỏ (Split), xoay trang (Rotate), sắp xếp lại thứ tự trang (Reorder). | Sử dụng **`github.com/pdfcpu/pdfcpu`** (thư viện xử lý PDF gốc bằng Go có hiệu năng rất cao). |
| **Convert to PDF** | Chuyển đổi các định dạng văn bản (HTML, Office, Images) sang PDF chất lượng cao. | Kết hợp **`github.com/jung-kurt/gofpdf`** (tạo PDF) hoặc tích hợp CLI wrapper cho **`wkhtmltopdf`** / **`libreoffice`**. |
| **OCR PDF** | Nhận diện ký tự quang học từ tài liệu scan để chuyển thành PDF có thể tìm kiếm văn bản (Searchable PDF). | Sử dụng Go bindings cho **Tesseract OCR** (`github.com/otiai10/gosseract/v2`). |
| **PDF/A Standard** | Định dạng PDF chuẩn lưu trữ lâu dài (PDF/A-1b, PDF/A-2b...). | Cấu hình profile xuất của `pdfcpu` hoặc hậu xử lý file thông qua công cụ tối ưu hóa cấu trúc tệp. |
| **Digital Signature** | Ký số tài liệu bằng chứng thư số (chữ ký số cá nhân/doanh nghiệp - chuẩn PAdES). | Sử dụng gói mã hóa gốc của Go (`crypto/x509`, `crypto/rsa`) kết hợp ghi đè byte range vào phần Trailer của PDF. |
| **Deep Redact** | Che đen thông tin nhạy cảm (bảng lương, số CCCD, tài khoản) vĩnh viễn khỏi file (không thể khôi phục). | Thực hiện vẽ đè các khối hình chữ nhật màu đen lên vị trí text nhạy cảm và gỡ bỏ hoàn toàn stream text cũ bên dưới. |
| **Watermark** | Đóng dấu chìm hình ảnh, chữ bản quyền lên tài liệu lương/chấm công. | Sử dụng API Watermark/Stamp của **`pdfcpu`**. |
| **Compress PDF** | Nén dung lượng file PDF tối đa nhưng giữ nguyên chất lượng hiển thị để lưu trữ. | Tối ưu hóa font, nén lại luồng hình ảnh bằng thuật toán deflate thông qua **`pdfcpu/pkg/api`**. |
| **Folder-based Automation** | Lắng nghe thư mục (Watch folder). Khi có tệp mới, tự động chạy pipeline xử lý và xuất kết quả. | Sử dụng **`github.com/fsnotify/fsnotify`** để theo dõi tệp tin thời gian thực trong OS. |

---

## 2. Kiến Trúc Hệ Thống (Architecture)

Module được thiết kế theo mô hình **Pipeline Pattern** giúp dễ dàng mở rộng và tái sử dụng các bước xử lý:

```
[Watch Folder / API Request]
            │
            ▼
┌────────────────────────┐
│    Job Orchestrator    │ (Điều phối luồng công việc)
└───────────┬────────────┘
            │
            ▼
┌────────────────────────┐
│     PDF Pipeline       │
│  ├─ Step 1: Compress   │ (Sử dụng pdfcpu)
│  ├─ Step 2: Watermark  │ (Sử dụng pdfcpu)
│  └─ Step 3: Sign       │ (Sử dụng crypto/x509)
└───────────┬────────────┘
            │
            ▼
   [Output Folder / DB]
```

---

## 3. Chi Tiết Thực Thi Kỹ Thuật (Technical Implementation)

### 3.1. Cấu trúc thư mục Module trong dự án

```
modules/pdf/
├── controller/        # Hứng các yêu cầu HTTP xử lý PDF từ Frontend
├── service/           # Logic xử lý chính (Merge, OCR, Compress...)
├── repository/        # Lưu trữ lịch sử xử lý file và log
├── watcher/           # Tự động hóa giám sát thư mục (fsnotify)
└── entity/            # Thực thể lưu cấu hình watermark/chữ ký số
```

### 3.2. Đoạn mã mẫu giám sát thư mục (Folder Watcher)

```go
package watcher

import (
	"context"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type FolderWatcher struct {
	watchDir  string
	outputDir string
}

func (w *FolderWatcher) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					// Chỉ xử lý file PDF mới
					if filepath.Ext(event.Name) == ".pdf" {
						log.Printf("Phát hiện file mới: %s, đang đưa vào pipeline xử lý...", event.Name)
						// Kích hoạt PDF Pipeline xử lý ngầm (Go goroutine)
						go w.processPipeline(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Lỗi giám sát thư mục:", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	return watcher.Add(w.watchDir)
}

func (w *FolderWatcher) processPipeline(filePath string) {
	// Thực hiện: Đóng watermark -> Nén -> Ký số -> Lưu ra thư mục output
}
```

---

## 4. Kế Hoạch Triển Khai (Deployment Roadmap)

1. **Sprint 1 (Core PDF Processing):** Tích hợp `pdfcpu` thực hiện các tác vụ gộp, chia, nén, đóng dấu và xoay trang.
2. **Sprint 2 (Security & Sign):** Tích hợp chữ ký số PAdES và cơ chế Deep Redact che dấu thông tin.
3. **Sprint 3 (OCR & Convert):** Tích hợp công cụ chuyển đổi tệp và nhận diện chữ viết scan qua Tesseract OCR.
4. **Sprint 4 (Automation & API):** Xây dựng Folder Watcher dùng `fsnotify` và hoàn tất API xử lý cho Frontend.
