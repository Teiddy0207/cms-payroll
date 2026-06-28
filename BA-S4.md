# BA-S4 · Module Động cơ công thức & Tính lương Batch (AST & Redis Lock)

> **Phiên BA:** S4 (Giai đoạn Khảo sát & Đặc tả) · **Ngày:** 2026-06-28
> **Nuôi chương đặc tả:** Chương 3.2.4 (Module Động cơ & Tính toán Batch), Chương 3.3 (Đối tượng bổ sung)
> **Master plan:** [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md) · **Walkthrough:** [ba_walkthrough.md](./ba_walkthrough.md)
> **Cần nạp context (dependency):** BA-S1, BA-S2, BA-S3
> **Nguồn:** Thảo luận trực tiếp với Chủ đầu tư ngày 2026-06-28

---

## 1. Mục tiêu phiên & phạm vi
Phiên này làm rõ cơ chế lưu vết phiên bản công thức lương (Versioning), quy trình chạy tính toán lương hàng loạt bằng Go Goroutines và Redis Lock, cơ chế duyệt đơn chỉnh sửa lương thủ công từ kế toán gửi Trưởng phòng, và phương thức bảo mật phê duyệt chi trả lương cuối cùng bằng nhận diện khuôn mặt của Giám đốc.

* **Mục tiêu:**
  * Đặc tả cơ chế quản lý phiên bản công thức lương (`Formula Versioning`) bảo đảm tính đúng đắn khi tính toán lịch sử.
  * Thiết kế luồng xử lý bất đồng bộ tính lương batch quy mô lớn (Eager Loading, RAM Map, Redis Lock).
  * Quy trình điều chỉnh lương thủ công: Kế toán chỉnh sửa $\rightarrow$ Trưởng phòng phê duyệt duyệt (áp dụng RLS) $\rightarrow$ Lưu vết thay đổi.
  * Thiết lập cơ chế kiểm soát chi trả lương (Disbursement Approval) bắt buộc xác thực bằng khuôn mặt (Face ID) của Giám đốc.
* **Phạm vi:** Phiên bản công thức, Go Worker Pool, luồng duyệt điều chỉnh lương, và phê duyệt chi tiền bằng nhận diện khuôn mặt Giám đốc.
* **Ngoài phạm vi:** Phân tích AI cảnh báo bất thường bảng lương và chat tra cứu RAG (sẽ làm rõ ở phiên BA-S5).

---

## 2. Hiện trạng (as-is — kiểm chứng bằng mã nguồn)
* Dự án khởi tạo mới từ đầu, chưa có mã nguồn nghiệp vụ hiện trạng.
* Toàn bộ cấu trúc sẽ được xây dựng mới theo tài liệu đặc tả này.

---

## 3. Yêu cầu (to-be — dẫn tiêu chuẩn)
Các yêu cầu đặc thù phục vụ tối ưu hóa hiệu năng và bảo mật hệ thống:
* **Công thức hồi tố (Retroactive Calculation):** Lưu trữ lịch sử thay đổi của mọi công thức. Khi chạy tính toán lại lương của một tháng trong quá khứ, hệ thống tự động tải đúng phiên bản công thức hoạt động tại thời điểm đó.
* **Quy trình điều chỉnh chặt chẽ:** Kế toán được phép điều chỉnh lương trên UI dạng Spreadsheet nhưng bắt buộc phải tạo phiếu yêu cầu phê duyệt gửi Trưởng phòng duyệt mới được áp dụng vào bảng lương chính thức.
* **Bảo mật sinh trắc học:** Sử dụng chính camera nhận dạng khuôn mặt trên thiết bị di động/máy tính của Giám đốc để ký duyệt lệnh chi trả, đảm bảo tính chống từ chối (Non-repudiation).

---

## 4. Quy tắc nghiệp vụ & tham số (sinh REQ)

| Mã | Phát biểu | Nguồn | Loại | As-is | Phiên | Nghiệm thu |
|:---|:---|:---|:---|:---|:---|:---|
| **REQ-FORM-001** | **Formula Versioning:** Mỗi công thức lương được lưu kèm khoảng thời gian hiệu lực (`start_date` đến `end_date`). Khi chạy tính lương chu kỳ $T$, hệ thống phải load các công thức có hiệu lực bao phủ chu kỳ $T$. | Khách hàng | Business Rule | Mới | S4 | [ ] |
| **REQ-PIPE-001** | **Chạy tính toán bất đồng bộ:** Kế toán bấm kích hoạt tính lương. API trả về lập tức `Job_ID` và mã `202 Accepted`. Tiến độ chạy song song (Goroutines) được cập nhật liên tục qua WebSocket lên UI Progress Bar của người dùng. | Khách hàng | Technical | Mới | S4 | [ ] |
| **REQ-PIPE-002** | **Tránh N+1 Query:** Động cơ nạp dữ liệu gộp theo lô 1,000 nhân viên (sử dụng toán tử `IN` nạp lên RAM Map), sau đó tính toán cuốn gói theo đồ thị Topological Sort trực tiếp trên RAM thay vì thực hiện đọc ghi liên tục xuống Database. | Khách hàng | Technical | Mới | S4 | [ ] |
| **REQ-ADJ-001**  | **Chỉnh sửa lương thủ công:** Kế toán có quyền đề xuất chỉnh sửa các cột lương thực nhận của cá nhân. Việc sửa đổi này sẽ tạo ra một phiếu `PayrollAdjustment` ghi nhận: Mã nhân sự, cột sửa, giá trị cũ, giá trị mới, lý do sửa và người đề xuất. | Khách hàng | Functional | Mới | S4 | [ ] |
| **REQ-ADJ-002**  | **Phê duyệt sửa lương:** Phiếu `PayrollAdjustment` phải được Trưởng phòng của bộ phận nhân viên đó phê duyệt (áp dụng RLS) mới được áp dụng và ghi đè vào bảng lương tổng hợp. | Khách hàng | Business Rule | Mới | S4 | [ ] |
| **REQ-LOCK-001** | **Khóa sổ bảng lương:** Sau khi kết thúc phê duyệt điều chỉnh, Kế toán thực hiện khóa bảng lương. Khi bảng lương ở trạng thái `LOCKED`, hệ thống cấm tuyệt đối mọi thao tác sửa đổi bảng công, đơn từ và bảng lương của tháng đó. | Khách hàng | Business Rule | Mới | S4 | [ ] |
| **REQ-DISB-001** | **Chi trả xác thực bằng Face ID:** Lệnh chi trả (Disbursement) chỉ được kích hoạt khi Giám đốc (Director) đăng nhập và thực hiện quét khuôn mặt trên ứng dụng. Hệ thống so khớp vector khuôn mặt của Giám đốc đạt độ tương đồng > 95% mới thực hiện mã hóa và phát lệnh chi trả. | Khách hàng | Security | Mới | S4 | [ ] |

---

## 5. Câu hỏi & Quyết định từ Chủ đầu tư (Q&A)

### Q-S4-01: Quản lý công thức khi có sự thay đổi chính sách thuế hoặc lương?
* **Quyết định:** Hệ thống phải lưu trữ lịch sử phiên bản của công thức. Khi tính toán lại lương của một tháng cũ trong quá khứ, hệ thống bắt buộc phải dùng các công thức và thông số có hiệu lực tại đúng thời điểm tháng cũ đó.

### Q-S4-02: Quy trình kế toán chỉnh sửa thủ công số liệu sau khi tính lương xong?
* **Quyết định:** Cho phép kế toán chỉnh sửa số liệu lương của cá nhân, nhưng không được tự ý áp dụng ngay. Kế toán phải nhập lý do và gửi yêu cầu phê duyệt lên Trưởng phòng của nhân sự đó duyệt. Hệ thống phải lưu lại đầy đủ vết thay đổi và số tiền sửa đổi. Lương sẽ được chi trả vào ngày 15 hàng tháng.

### Q-S4-03: Khóa sổ bảng lương và cơ chế duyệt chi tiền?
* **Quyết định:** Khi đã khóa bảng lương thì cấm tuyệt đối, không cho phép chỉnh sửa bất kỳ số liệu nào của tháng đó nữa. Lệnh chi tiền cuối cùng bắt buộc phải quét nhận diện khuôn mặt của Giám đốc để xác nhận phê duyệt mới được chạy.

---

## 6. Tiêu chí chấp nhận sơ bộ của phiên (Acceptance Criteria)

- [ ] Cấu trúc bảng `payroll_formulas` lưu trữ thời gian hiệu lực và mã phiên bản.
- [ ] Thiết kế luồng Go Worker Pool thực hiện tính toán độc lập, lưu trữ kết quả tạm vào bảng nháp `payroll_records_draft`.
- [ ] Thiết kế API và giao diện xác thực khuôn mặt (Face Recognition Verification API) dành cho Giám đốc khi phê duyệt chi trả lương.

---

## 7. Đóng góp vào chương đặc tả

| Mục phiên này | Chương đích |
|---|---|
| Quy tắc REQ-FORM-001 | Chương 3.2.4.1 (Quản lý phiên bản công thức lương) |
| Quy tắc REQ-PIPE-001, REQ-PIPE-002 | Chương 3.2.4.2 (Kiến trúc tính lương Batch hiệu năng cao) |
| Quy tắc REQ-ADJ-001, REQ-ADJ-002, REQ-LOCK-001 | Chương 3.2.4.3 (Quy trình điều chỉnh và Khóa sổ bảng lương) |
| Quy tắc REQ-DISB-001 | Chương 3.2.4.4 (Cơ chế duyệt chi tiền bằng Face ID của Giám đốc) |
| Cấu trúc dữ liệu Đối tượng | Chương 3.3 (Cơ cấu đối tượng bổ sung - Công thức & Phiếu điều chỉnh) |
