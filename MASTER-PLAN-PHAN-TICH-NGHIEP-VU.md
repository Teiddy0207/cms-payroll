# MASTER PLAN PHÂN TÍCH NGHIỆP VỤ & ĐẶC TẢ SRS
## HỆ THỐNG TÍNH LƯƠNG PHÂN TÁN (DISTRIBUTED PAYROLL PLATFORM)

---

## 1. MỤC TIÊU & PHƯƠNG PHÁP
* **Mục tiêu:** Chuyển đổi từ mô tả ý tưởng kỹ thuật ban đầu thành bộ Tài liệu đặc tả yêu cầu phần mềm (SRS) chuẩn hóa, chi tiết và có thể chuyển giao trực tiếp cho đội ngũ phát triển.
* **Vai trò:** 
  * **BA (Antigravity):** Chịu trách nhiệm phỏng vấn, gợi mở yêu cầu, chuẩn hóa nghiệp vụ và soạn thảo tài liệu.
  * **Chủ đầu tư / Khách hàng (User):** Cung cấp nghiệp vụ thực tế, phê duyệt các quyết định thiết kế và quy tắc nghiệp vụ.
* **Phương pháp:** Tiến hành phân tích cuốn gói theo từng Phiên BA (Session). Mỗi phiên sẽ tập trung làm rõ một nhóm chức năng (Module), tạo file phân tích phiên và cập nhật trực tiếp vào chương đích của tài liệu đặc tả hệ thống.

---

## 2. KẾ HOẠCH PHÂN CHIA PHIÊN BA (SESSIONS ROADMAP)

| Phiên BA | Module trọng tâm | Mục tiêu phân tích chính | Tài liệu đặc tả |
| :--- | :--- | :--- | :--- |
| **BA-S1** | **Cơ cấu tổ chức & Phân quyền (RBAC)** | Phân quyền vai trò (Admin, HR, Kế toán, Nhân viên), cấu hình sơ đồ phòng ban, chức vụ. | Chương 3.2.1 |
| **BA-S2** | **Hồ sơ nhân sự & Khung năng lực (P1 & P2)** | Quản lý hợp đồng, cách tính toán chi tiết tiêu chuẩn chức danh (P1) và khung điểm năng lực (P2). | Chương 3.2.2 |
| **BA-S3** | **Quản lý ca trực & Chấm công khuôn mặt (NATS)** | Lịch phân ca, luồng đẩy sự kiện chấm công qua NATS JetStream, kiểm soát trùng lặp và tính công ngày. | Chương 3.2.3 |
| **BA-S4** | **Động cơ công thức & Tính lương Batch (AST)** | Thiết lập công thức động bằng AST, giải thuật sắp xếp topo, tối ưu bộ nhớ RAM, xử lý khóa Redis. | Chương 3.2.4 |
| **BA-S5** | **Kiểm toán AI (Auditor) & RAG Chatbot** | Phương pháp phát hiện bất thường ($3\sigma$ + LLM), chatbot tra cứu thông tin cá nhân bảo mật (RLS). | Chương 3.2.5 |

---

## 3. TÀI LIỆU THAM CHIẾU & HIỆN TRẠNG CỐT LÕI
* **Core Stack:** Golang (Gin/Expr), NATS JetStream, ClickHouse, PostgreSQL, Redis, Next.js.
* **Mẫu tài liệu đặc tả SRS:** [HTCBT_SRS_skeleton.md](./HTCBT_SRS_skeleton.md) (sẽ được cấu trúc lại cho dự án lương 3P).
* **Mẫu phiên phân tích BA:** [_BA-TEMPLATE.md](./_BA-TEMPLATE.md).
