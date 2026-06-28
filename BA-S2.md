# BA-S2 · Module Hồ sơ nhân sự & Tính toán P1, P2

> **Phiên BA:** S2 (Giai đoạn Khảo sát & Đặc tả) · **Ngày:** 2026-06-28
> **Nuôi chương đặc tả:** Chương 3.2.2 (Module Nhân sự & Hợp đồng), Chương 3.3 (Đối tượng)
> **Master plan:** [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md) · **Walkthrough:** [ba_walkthrough.md](./ba_walkthrough.md)
> **Cần nạp context (dependency):** BA-S1 (để kế thừa phân quyền RLS của Trưởng phòng khi chấm điểm năng lực)
> **Nguồn:** Thảo luận trực tiếp với Chủ đầu tư ngày 2026-06-28

---

## 1. Mục tiêu phiên & phạm vi
Phiên này làm rõ nghiệp vụ lưu trữ hồ sơ nhân sự, hợp đồng lao động, cấu hình các tiêu chuẩn vị trí (P1), tự định nghĩa khung năng lực và chấm điểm (P2), cùng tính năng nhập hồ sơ hàng loạt từ Excel.

* **Mục tiêu:**
  * Đặc tả cấu trúc lưu trữ Hồ sơ nhân viên và Hợp đồng lao động.
  * Thiết lập công thức tính Lương vị trí (`P1_BASE`) cộng dồn theo chứng chỉ, ngành nghề và tiêu chuẩn chức danh.
  * Định cấu hình lương Năng lực (`P2_COMPETENCY`) dựa trên các chỉ mục năng lực tự chọn, cơ chế đánh giá 6 tháng/lần của Trưởng phòng và quy đổi điểm ra tiền tệ theo cấu hình doanh nghiệp.
  * Thống nhất cơ chế nhập dữ liệu hàng loạt bằng Excel hỗ trợ ánh xạ cột động (Custom mapping).
* **Phạm vi:** Nghiệp vụ nhân viên, hợp đồng, thiết lập tiêu chuẩn P1, bộ chỉ số P2 và nhập liệu Excel.
* **Ngoài phạm vi:** Phân ca làm việc và tính công hàng ngày (sẽ được làm rõ ở phiên BA-S3).

---

## 2. Hiện trạng (as-is — kiểm chứng bằng mã nguồn)
* Dự án khởi tạo mới từ đầu, chưa có mã nguồn nghiệp vụ hiện trạng.
* Toàn bộ cấu trúc sẽ được xây dựng mới theo tài liệu đặc tả này.

---

## 3. Yêu cầu (to-be — dẫn tiêu chuẩn)
Các yêu cầu được thống nhất với Chủ đầu tư nhằm tự động hóa tối đa:
* Lương vị trí (`P1_BASE`) được tính tự động dựa trên mức lương nền của chức danh cộng thêm tiền thưởng của các chứng chỉ và chuyên môn đạt được.
* Lương năng lực (`P2_COMPETENCY`) được đánh giá định kỳ 6 tháng/lần bởi Trưởng phòng của bộ phận (áp dụng RLS, trưởng phòng không được đánh giá nhân sự phòng khác). Khung năng lực có thể tùy biến linh hoạt cho từng công ty và phòng ban.
* Hệ thống hỗ trợ nạp Excel tùy biến cấu hình cột để doanh nghiệp dễ dàng chuyển đổi dữ liệu từ hệ thống cũ sang.

---

## 4. Quy tắc nghiệp vụ & tham số (sinh REQ)

| Mã | Phát biểu | Nguồn | Loại | As-is | Phiên | Nghiệm thu |
|:---|:---|:---|:---|:---|:---|:---|
| **REQ-EMP-001** | Quản lý thông tin cá nhân nhân viên (Họ tên, ngày sinh, mã nhân viên, phòng ban, vai trò). | Khách hàng | Functional | Mới | S2 | [ ] |
| **REQ-EMP-002** | **Nhập liệu hàng loạt (Bulk Import):** Hệ thống cho phép HR tải lên file Excel chứa danh sách nhân viên. Giao diện Next.js cung cấp bộ ánh xạ (Mapping Tool) để người dùng tự chọn Cột Excel khớp với Trường dữ liệu tương ứng trong DB (Custom Template Mapping). | Khách hàng | Functional | Mới | S2 | [ ] |
| **REQ-CON-001** | Mỗi nhân viên có thể có nhiều Hợp đồng lao động theo thời gian, nhưng tại một thời điểm chỉ có tối đa 01 Hợp đồng có trạng thái "Đang kích hoạt". | Khách hàng | Business Rule | Mới | S2 | [ ] |
| **REQ-P1-001**  | HR có thể định nghĩa danh mục "Tiêu chuẩn vị trí" (ví dụ: Chứng chỉ Tiếng Anh = 1,000,000đ, Ngành công nghệ cao = 2,000,000đ). | Khách hàng | Functional | Mới | S2 | [ ] |
| **REQ-P1-002**  | **Công thức P1_BASE:** `P1_BASE = Position_Base_Rate + SUM(Standard_Value)`. Trong đó, `Position_Base_Rate` là lương sàn của vị trí chức danh (ví dụ: Kế toán = 10,000,000đ) và `Standard_Value` là các tiêu chuẩn bắt buộc/đạt được đi kèm vị trí đó. | Khách hàng | Business Rule | Mới | S2 | [ ] |
| **REQ-P2-001**  | HR có thể tự tạo mới, chỉnh sửa các chỉ số năng lực trong "Khung năng lực" của công ty (ví dụ: Kỹ năng viết code, Tư duy logic, Giao tiếp khách hàng). | Khách hàng | Functional | Mới | S2 | [ ] |
| **REQ-P2-002**  | Trưởng phòng (Head of Department) thực hiện chấm điểm năng lực (thang điểm 1-100) cho nhân sự trong phòng ban mình quản lý. Tần suất đánh giá bắt buộc là **06 tháng một lần (đánh giá bán niên)**. | Khách hàng | Business Rule | Mới | S2 | [ ] |
| **REQ-P2-003**  | **Công thức P2_COMPETENCY:** `P2_COMPETENCY = SUM(Competency_Score * Competency_Weight) * Company_Point_Rate`. Trong đó, `Competency_Weight` là trọng số của năng lực đối với chức danh (từ 1 đến 5), và `Company_Point_Rate` là đơn giá quy đổi 1 điểm năng lực thành VND do công ty cấu hình (ví dụ: 5,000đ/điểm). | Khách hàng | Business Rule | Mới | S2 | [ ] |

---

## 5. Câu hỏi & Quyết định từ Chủ đầu tư (Q&A)

### Q-S2-01: Cách cộng các tiêu chuẩn vào P1_BASE?
* **Quyết định:** P1_BASE sẽ bằng tổng các tiêu chuẩn (chứng chỉ, ngành nghề chuyên môn) cộng với mức lương nền của vị trí chức danh đó. Ví dụ vị trí kế toán sẽ có tổng lương nền chức danh + tiêu chuẩn X + tiêu chuẩn Y.

### Q-S2-02: Cơ chế thiết lập khung năng lực P2 và chấm điểm?
* **Quyết định:** Khung năng lực được HR tự do cấu hình thêm bớt chỉ số trên UI. Trưởng phòng là người đánh giá điểm cho nhân viên trực thuộc định kỳ 6 tháng một lần. Tiền cơ bản quy đổi trên mỗi điểm năng lực do công ty tự thiết lập.

### Q-S2-03: Yêu cầu về nhập hồ sơ bằng file Excel?
* **Quyết định:** Hệ thống phải hỗ trợ nhập dữ liệu nhân sự hàng loạt bằng file Excel, đồng thời cung cấp giao diện ánh xạ cột linh động để phù hợp với từng mẫu Excel khác nhau của các công ty.

---

## 6. Tiêu chí chấp nhận sơ bộ của phiên (Acceptance Criteria)

- [ ] Thiết kế bảng cơ sở dữ liệu `employees`, `contracts`, `job_standards`, `competencies` và `competency_evaluations` có lưu vết ngày đánh giá.
- [ ] Thiết kế API và giao diện import Excel hỗ trợ ánh xạ cột động.
- [ ] Xây dựng luồng tính toán thử nghiệm `P1_BASE` và `P2_COMPETENCY` trả về kết quả chính xác theo các tham số cấu hình động.

---

## 7. Đóng góp vào chương đặc tả

| Mục phiên này | Chương đích |
|---|---|
| Quy tắc REQ-EMP-001, REQ-EMP-002, REQ-CON-001 | Chương 3.2.2.1 (Quản lý thông tin nhân viên & Import Excel) |
| Quy tắc REQ-P1-001, REQ-P1-002 | Chương 3.2.2.2 (Cấu hình và tính toán Lương vị trí P1) |
| Quy tắc REQ-P2-001, REQ-P2-002, REQ-P2-003 | Chương 3.2.2.3 (Cơ cấu đánh giá và tính toán Lương năng lực P2) |
| Cấu trúc dữ liệu Đối tượng | Chương 3.3 (Cơ cấu đối tượng chính - Entity schemas) |
