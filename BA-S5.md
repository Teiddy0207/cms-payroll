# BA-S5 · Module AI Kiểm toán & Phân tích RAG (DeepSeek/OpenAI & pgvector)

> **Phiên BA:** S5 (Giai đoạn Khảo sát & Đặc tả) · **Ngày:** 2026-06-28
> **Nuôi chương đặc tả:** Chương 3.2.5 (Module AI & Analytics), Chương 3.3 (Đối tượng bổ sung), Chương 3.4, 3.5, 3.6 (Yêu cầu phi chức năng & Ràng buộc)
> **Master plan:** [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md) · **Walkthrough:** [ba_walkthrough.md](./ba_walkthrough.md)
> **Cần nạp context (dependency):** BA-S1, BA-S2, BA-S3, BA-S4
> **Nguồn:** Thảo luận trực tiếp với Chủ đầu tư ngày 2026-06-28

---

## 1. Mục tiêu phiên & phạm vi
Phiên này làm rõ cơ chế kiểm toán bảng lương tự động bằng AI, gửi cảnh báo qua email/dashboard, thiết lập thanh chuyên mục "AI Phân tích" (Sidebar) hỗ trợ phân quyền RLS 3 cấp độ truy cập dữ liệu và cơ chế nạp tài liệu học quy chế lương nội bộ (RAG).

* **Mục tiêu:**
  * Đặc tả luồng kích hoạt tự động AI Auditor ngay sau khi chạy xong bảng lương Batch.
  * Thiết lập cơ chế gửi cảnh báo lỗi chấm công/lương lên Dashboard và gửi email tổng hợp cho Kế toán.
  * Thiết kế phân hệ "AI Phân tích" dưới dạng Sidebar hỗ trợ bảo mật thông tin dòng chặt chẽ cho 3 đối tượng (Employee, Trưởng phòng, Giám đốc/Kế toán).
  * Định nghĩa cơ chế nạp tài liệu định dạng PDF/DOCX cho AI học luật lương nội bộ.
* **Phạm vi:** Trình duyệt AI Auditor, báo cáo tự động, Sidebar "AI Phân tích" phân quyền RLS, và tải lên tài liệu RAG.
* **Ngoài phạm vi:** Đóng gói và bàn giao toàn bộ tài liệu đặc tả SRS chính thức.

---

## 2. Hiện trạng (as-is — kiểm chứng bằng mã nguồn)
* Dự án khởi tạo mới từ đầu, chưa có mã nguồn nghiệp vụ hiện trạng.
* Toàn bộ cấu trúc sẽ được xây dựng mới theo tài liệu đặc tả này.

---

## 3. Yêu cầu (to-be — dẫn tiêu chuẩn)
Các yêu cầu đặc thù về AI và Bảo mật dữ liệu trong môi trường doanh nghiệp:
* **Tự động hóa hoàn toàn:** Kế toán không cần thao tác thủ công, AI tự động quét ngay khi có bảng lương nháp mới.
* **Sidebar AI Phân tích bảo mật:** Thay vì dùng chatbot thả trôi ở góc màn hình, hệ thống bố trí một chuyên trang/sidebar phân tích riêng biệt. Hệ thống lọc dữ liệu (Data Guardrails) trước khi nạp vào Prompt gửi lên LLM để đảm bảo không rò rỉ dữ liệu.
* **Học quy chế số:** Đọc hiểu và đối chiếu bảng lương thực tế với các file văn bản quy chế pháp lý nội bộ để trả lời thắc mắc chuẩn xác nhất.

---

## 4. Quy tắc nghiệp vụ & tham số (sinh REQ)

| Mã | Phát biểu | Nguồn | Loại | As-is | Phiên | Nghiệm thu |
|:---|:---|:---|:---|:---|:---|:---|
| **REQ-AI-AUDIT-001** | **Kích hoạt tự động:** AI Auditor tự động khởi chạy ngầm ngay sau khi Job tính lương Batch hoàn thành công việc thành công. | Khách hàng | Business Rule | Mới | S5 | [ ] |
| **REQ-AI-AUDIT-002** | ** Baseline cấu hình động:** HR được quyền cấu hình số tháng lịch sử dùng làm mốc so sánh độ lệch chuẩn (mặc định là 06 tháng gần nhất). | Khách hàng | Functional | Mới | S5 | [ ] |
| **REQ-AI-AUDIT-003** | **Báo cáo kiểm toán:** Các bản ghi bất thường được xuất ra màn hình Dashboard kiểm toán của kế toán và gửi email thông báo tự động đính kèm danh sách chi tiết. | Khách hàng | Functional | Mới | S5 | [ ] |
| **REQ-AI-RAG-001**   | **Nạp quy chế nội bộ:** Cho phép HR tải lên các tài liệu quy chế (`.pdf`, `.docx`) qua giao diện Admin. Hệ thống thực hiện bóc tách văn bản, cắt nhỏ (Chunking), trích xuất vector embedding và lưu trữ vào Vector DB để phục vụ tra cứu ngữ cảnh. | Khách hàng | Functional | Mới | S5 | [ ] |
| **REQ-AI-CHAT-001**  | **Sidebar AI Phân tích:** Hệ thống tích hợp một Sidebar riêng biệt mang tên "AI Phân tích" hiển thị ở màn hình làm việc để tương tác và hỏi đáp số liệu trực quan. | Khách hàng | Functional | Mới | S5 | [ ] |
| **REQ-AI-CHAT-002**  | **Phân quyền RLS 3 Cấp độ trong AI Phân tích:**<br>1. **Nhân viên:** Chỉ được hỏi đáp dữ liệu lương/công của chính mình.<br>2. **Trưởng phòng:** Được hỏi đáp dữ liệu của bản thân và của toàn bộ nhân viên thuộc phòng mình quản lý.<br>3. **Giám đốc / Kế toán:** Được phép hỏi dữ liệu phân tích, tổng hợp của toàn công ty. | Khách hàng | Security | Mới | S5 | [ ] |

---

## 5. Câu hỏi & Quyết định từ Chủ đầu tư (Q&A)

### Q-S5-01: Cách thức vận hành và báo cáo của AI Auditor?
* **Quyết định:** AI Auditor sẽ tự chạy tự động ngay sau khi tính lương xong. Kết quả cảnh báo sẽ xuất lên màn hình Dashboard và gửi báo cáo về hòm thư Email của kế toán. Mốc lịch sử so sánh có thể tự cấu hình số tháng.

### Q-S5-02: Cơ chế phân quyền hỏi đáp dữ liệu trên Chatbot AI?
* **Quyết định:** Phân quyền chặt chẽ 3 cấp. Nhân viên chỉ được xem của mình. Trưởng phòng được xem của nhân viên phòng mình. Giám đốc/Kế toán được hỏi đáp phân tích số liệu toàn công ty.

### Q-S5-03: Giao diện tương tác AI và nạp tài liệu học?
* **Quyết định:** Không làm widget chat tròn ở góc màn hình mà thiết kế dạng **Sidebar riêng biệt có tên "AI Phân tích"**. Cho phép tải các file văn bản quy chế lên để RAG học luật tính lương nội bộ.

---

## 6. Tiêu chí chấp nhận sơ bộ của phiên (Acceptance Criteria)

- [ ] Thiết kế bảng `ai_anomalies` và bảng `rag_documents` (lưu trữ vector văn bản đã chia nhỏ).
- [ ] Xây dựng bộ lọc trung gian (Guardrail middleware) chặn và kiểm tra phân quyền RLS trước khi gửi payload câu hỏi + context lên các API LLM (OpenAI/DeepSeek).
- [ ] Kiểm chứng việc gửi email thông báo kết quả kiểm toán bằng mock SMTP server.

---

## 7. Đóng góp vào chương đặc tả

| Mục phiên này | Chương đích |
|---|---|
| Quy tắc REQ-AI-AUDIT-001, 002, 003 | Chương 3.2.5.1 (Động cơ kiểm toán AI Auditor) |
| Quy tắc REQ-AI-RAG-001, REQ-AI-CHAT-001, 002 | Chương 3.2.5.2 (Thanh AI Phân tích & Kiến trúc RAG bảo mật) |
| Cấu trúc dữ liệu Đối tượng | Chương 3.3 (Cơ cấu đối tượng bổ sung - Log Chat & Tài liệu RAG) |
| Yêu cầu phi chức năng, ràng buộc thiết kế | Chương 3.4, 3.5, 3.6 |
