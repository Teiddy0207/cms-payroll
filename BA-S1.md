# BA-S1 · Module Cơ cấu tổ chức & Xác thực phân quyền (RBAC/RLS)

> **Phiên BA:** S1 (Giai đoạn Khảo sát & Đặc tả) · **Ngày:** 2026-06-28
> **Nuôi chương đặc tả:** Chương 1, Chương 2, Chương 3.2.1
> **Master plan:** [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md) · **Walkthrough:** [ba_walkthrough.md](./ba_walkthrough.md)
> **Cần nạp context (dependency):** Không có (Phiên khởi tạo)
> **Nguồn:** Thảo luận trực tiếp với Chủ đầu tư ngày 2026-06-28

---

## 1. Mục tiêu phiên & phạm vi
Phiên này làm rõ cơ cấu tổ chức phòng ban, mô hình phân quyền vai trò (RBAC) kết hợp bảo mật mức dòng (RLS), cơ chế xác thực người dùng và yêu cầu ghi vết lịch sử hoạt động (Audit Logs).

* **Mục tiêu:** 
  * Định hình sơ đồ tổ chức 2 cấp (Công ty $\rightarrow$ Phòng ban).
  * Định nghĩa ma trận vai trò (5 vai trò mặc định) và phân quyền RLS cho cấp Trưởng phòng.
  * Thống nhất cơ chế xác thực Email/Password kết hợp sẵn sàng tích hợp Keycloak SSO.
  * Xác định các hành vi cần Audit và ràng buộc nhập lý do khi điều chỉnh số liệu.
* **Phạm vi:** Nghiệp vụ quản lý phòng ban, gán chức vụ, định danh, phân quyền và ghi log bảo mật.
* **Ngoài phạm vi:** Quản lý chi tiết hồ sơ nhân viên và hợp đồng lao động (sẽ làm rõ ở phiên BA-S2).

---

## 2. Hiện trạng (as-is — kiểm chứng bằng mã nguồn)
* Dự án khởi tạo mới từ đầu (Greenfield), chưa có mã nguồn nghiệp vụ hiện trạng.
* Toàn bộ cấu trúc sẽ được xây dựng mới theo tài liệu đặc tả này.

---

## 3. Yêu cầu (to-be — dẫn tiêu chuẩn)
Các yêu cầu được thống nhất trực tiếp với Chủ đầu tư để thiết lập hệ thống Enterprise SaaS:
* Cấu trúc tổ chức tinh gọn: **Công ty $\rightarrow$ Phòng ban**. Mỗi phòng ban có tối đa 1 Trưởng phòng tại một thời điểm để điều hành nhân sự trong bộ phận đó.
* Phân quyền bảo mật cao: Hỗ trợ cả **RBAC (Role-Based Access Control)** và **RLS (Row-Level Security)** để ngăn chặn rò rỉ thông tin lương chéo giữa các phòng ban.
* Ràng buộc kiểm toán chặt chẽ: Tất cả các thao tác sửa đổi công/lương thủ công phải được lưu vết cùng lý do giải trình bắt buộc.

---

## 4. Quy tắc nghiệp vụ & tham số (sinh REQ)

| Mã | Phát biểu | Nguồn | Loại | As-is | Phiên | Nghiệm thu |
|:---|:---|:---|:---|:---|:---|:---|
| **REQ-ORG-001** | Hệ thống hỗ trợ cơ cấu tổ chức 2 cấp: Công ty và các Phòng ban trực thuộc. | Khách hàng | Functional | Mới | S1 | [ ] |
| **REQ-ORG-002** | Mỗi Phòng ban có tối đa 01 Trưởng phòng (Head of Department) quản lý và nhiều Nhân viên trực thuộc. | Khách hàng | Business Rule | Mới | S1 | [ ] |
| **REQ-AUTH-001**| Đăng nhập mặc định qua Email và Password. Token xác thực là JWT (Access Token & Refresh Token) lưu trữ trạng thái phiên trên Redis. | Khách hàng | Security | Mới | S1 | [ ] |
| **REQ-AUTH-002**| Cấu trúc hệ thống sẵn sàng tích hợp giao thức OpenID Connect (OIDC) để hỗ trợ SSO (Keycloak) khi cấu hình. | Khách hàng | Technical | Mới | S1 | [ ] |
| **REQ-RBAC-001**| Hệ thống định nghĩa sẵn 05 vai trò mặc định: Admin, Giám đốc, Kế toán, Trưởng phòng, Nhân viên. | Khách hàng | Business Rule | Mới | S1 | [ ] |
| **REQ-RLS-001** | **Bảo mật mức dòng (RLS):** Trưởng phòng chỉ được xem thông tin lương, bảng công và phê duyệt đơn từ của nhân sự thuộc phòng ban mình quản lý. Nghiêm cấm truy cập dữ liệu phòng ban khác. | Khách hàng | Security | Mới | S1 | [ ] |
| **REQ-RLS-002** | Giám đốc được xem toàn bộ bảng lương công ty. Kế toán được xem và điều chỉnh bảng lương toàn công ty. | Khách hàng | Security | Mới | S1 | [ ] |
| **REQ-AUDIT-001**| Ghi vết lịch sử (Audit Log) đối với các hành vi: Đăng nhập/Đăng xuất, Sửa đổi bảng công, Kích hoạt tính lương, Điều chỉnh lương thực nhận. | Khách hàng | Functional | Mới | S1 | [ ] |
| **REQ-AUDIT-002**| Khi người dùng thực hiện Sửa đổi bảng công hoặc Điều chỉnh lương thực nhận, hệ thống bắt buộc phải yêu cầu nhập "Lý do điều chỉnh" (Text field, tối thiểu 10 ký tự) mới cho phép lưu. | Khách hàng | Business Rule | Mới | S1 | [ ] |

---

## 5. Câu hỏi & Quyết định từ Chủ đầu tư (Q&A)

### Q-S1-01: Cơ cấu tổ chức có cần phân cấp sâu hơn không?
* **Quyết định:** Không, cấu trúc đơn giản gồm **Công ty $\rightarrow$ Phòng ban**. Trong phòng ban có 1 trưởng phòng và nhiều nhân sự. Trưởng phòng chỉ điều hành nhân sự phòng ban đó.

### Q-S1-02: Các vai trò hệ thống và phân quyền dữ liệu?
* **Quyết định:** Hệ thống gồm 5 vai trò (Admin, Giám đốc, Kế toán, Trưởng phòng, Nhân viên). Áp dụng Row-Level Security: Trưởng phòng chỉ thao tác trên nhân sự phòng ban mình.

### Q-S1-03: Ghi vết kiểm toán và ràng buộc dữ liệu?
* **Quyết định:** Ghi log các thao tác đăng nhập, sửa công, chạy lương, sửa lương. Khi sửa công và lương bắt buộc phải nhập lý do giải trình.

### Q-S1-04: Phương thức xác thực?
* **Quyết định:** Đăng nhập Email/Password chuẩn, hỗ trợ tích hợp SSO Keycloak.

---

## 6. Tiêu chí chấp nhận sơ bộ của phiên (Acceptance Criteria)

- [ ] Thiết lập cấu trúc cơ sở dữ liệu (Database Schema) hỗ trợ mối quan hệ Công ty - Phòng ban - Chức vụ và gán Trưởng phòng.
- [ ] Thiết kế API đăng nhập JWT và API xác thực phân quyền RBAC/RLS hoạt động chính xác theo phân quyền dòng.
- [ ] Cơ chế ghi log kiểm toán bắt buộc nhập trường lý do (reason) khi gọi API cập nhật công/lương.

---

## 7. Đóng góp vào chương đặc tả

| Mục phiên này | Chương đích |
|---|---|
| Quy tắc REQ-ORG-001, REQ-ORG-002 | Chương 2 (Mô tả tổng quan) & Chương 3.2.1.1 (Quản lý phòng ban) |
| Quy tắc REQ-AUTH, REQ-RBAC, REQ-RLS | Chương 3.2.1.2 (Xác thực và Phân quyền RBAC/RLS) |
| Quy tắc REQ-AUDIT | Chương 3.2.1.3 (Quản lý vết kiểm toán - Audit Logs) |
