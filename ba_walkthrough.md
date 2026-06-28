# BA WALKTHROUGH · SỔ THEO DÕI TIẾN TRÌNH PHÂN TÍCH NGHIỆP VỤ

---

## 1. BẢNG THEO DÕI TIẾN ĐỘ PHIÊN BA

| Phiên BA | Tên phiên | Ngày thực hiện | Trạng thái | Đầu ra |
|:---|:---|:---|:---|:---|
| **BA-S1** | Cơ cấu tổ chức & Phân quyền (RBAC) | 2026-06-28 | [x] Hoàn thành | [BA-S1.md](./BA-S1.md) |
| **BA-S2** | Hồ sơ nhân sự & Tính toán P1, P2 | 2026-06-28 | [x] Hoàn thành | [BA-S2.md](./BA-S2.md) |
| **BA-S3** | Chấm công khuôn mặt & Luồng NATS | 2026-06-28 | [x] Hoàn thành | [BA-S3.md](./BA-S3.md) |
| **BA-S4** | Động cơ công thức & Tính lương Batch | 2026-06-28 | [x] Hoàn thành | [BA-S4.md](./BA-S4.md) |
| **BA-S5** | AI Kiểm toán & RAG Chatbot Tra cứu | 2026-06-28 | [x] Hoàn thành | [BA-S5.md](./BA-S5.md) |

---

## 2. SỔ CÂU HỎI MỞ & QUYẾT ĐỊNH CẦN CHỦ ĐẦU TƯ

| Mã câu hỏi | Phiên phát sinh | Nội dung câu hỏi / Điểm cần làm rõ | Trạng thái | Quyết định của Chủ đầu tư |
|:---|:---|:---|:---|:---|
| **Q-S1-01** | BA-S1 | Cơ cấu phòng ban có giới hạn số cấp phân cấp hay không? | [x] Đã giải quyết | Cấu trúc tinh gọn 2 cấp: Công ty -> Phòng ban. |
| **Q-S1-02** | BA-S1 | Cơ chế phân quyền RBAC có cần hỗ trợ phân quyền tùy biến tới từng API hay chỉ phân quyền theo nhóm chức năng lớn? | [x] Đã giải quyết | Phân quyền 5 vai trò hệ thống và áp dụng RLS (Trưởng phòng chỉ xem dữ liệu phòng mình). |
| **Q-S1-03** | BA-S1 | Cơ chế ghi vết và ràng buộc nhập liệu? | [x] Đã giải quyết | Ghi log đăng nhập, sửa công, tính lương, sửa lương. Bắt buộc nhập lý do khi sửa đổi. |
| **Q-S1-04** | BA-S1 | Cơ chế xác thực người dùng? | [x] Đã giải quyết | Đăng nhập Email/Password cơ bản, hỗ trợ SSO Keycloak. |
| **Q-S2-01** | BA-S2 | Cơ chế tính lương vị trí P1_BASE? | [x] Đã giải quyết | Lương nền vị trí cộng dồn các tiêu chuẩn đi kèm (bằng cấp, chứng chỉ, ngành nghề). |
| **Q-S2-02** | BA-S2 | Cơ chế chấm điểm và tính phụ cấp năng lực P2? | [x] Đã giải quyết | Khung năng lực tự chọn, Trưởng phòng đánh giá định kỳ 6 tháng/lần, đơn giá điểm do công ty quy định. |
| **Q-S2-03** | BA-S2 | Cơ chế import hồ sơ bằng file Excel? | [x] Đã giải quyết | Hỗ trợ tải Excel lên và ánh xạ cột động trên giao diện. |
| **Q-S3-01** | BA-S3 | Cách tính ca và hệ số OT? | [x] Đã giải quyết | Ca hành chính 8:30-17:30, grace 30 phút. OT đêm 1.5, lễ nghỉ 2.0. Phải có Trưởng phòng duyệt mới tính lương OT. |
| **Q-S3-02** | BA-S3 | Thiết bị chấm công và cơ chế IN/OUT? | [x] Đã giải quyết | Chấm công qua Mobile App nhận diện khuôn mặt gửi Employee_ID về. Quẹt tối thiểu 2 lần/ngày để ghi nhận Check-in/out. |
| **Q-S3-03** | BA-S3 | Quy tắc tính công ngày và bù công? | [x] Đã giải quyết | Sáng vào trước 9h, chiều ra sau 17h30 tính 1 công. Ra từ 12h-14h tính 0.5 công. Đi muộn/về sớm viết giải trình Trưởng phòng duyệt để bù thành 1 công. |
| **Q-S4-01** | BA-S4 | Cơ chế lưu trữ và quản lý công thức lương khi thay đổi? | [x] Đã giải quyết | Sử dụng Formula Versioning. Khi tính lại lương tháng cũ bắt buộc sử dụng đúng công thức và dữ liệu của thời điểm cũ đó. |
| **Q-S4-02** | BA-S4 | Cơ chế điều chỉnh lương thủ công của kế toán? | [x] Đã giải quyết | Cho phép điều chỉnh lương thực nhận nhưng phải viết lý do và gửi lên Trưởng phòng của nhân sự đó duyệt. Lưu lại toàn bộ vết điều chỉnh. Lương trả vào ngày 15 hàng tháng. |
| **Q-S4-03** | BA-S4 | Khóa sổ bảng lương và duyệt lệnh chi trả? | [x] Đã giải quyết | Bảng lương đã khóa là cấm sửa đổi tuyệt đối. Lệnh chi trả (Disbursement) bắt buộc quét khuôn mặt Giám đốc để xác nhận phê duyệt. |
| **Q-S5-01** | BA-S5 | Cơ chế vận hành và báo cáo của AI Auditor? | [x] Đã giải quyết | Tự động chạy quét bảng lương khi tính xong. Báo cáo xuất ra Dashboard và gửi Email cho kế toán. Cho phép tự chỉnh số tháng lịch sử làm Baseline. |
| **Q-S5-02** | BA-S5 | Phân quyền truy cập thông tin trên AI? | [x] Đã giải quyết | RLS 3 cấp: Nhân viên hỏi của mình. Trưởng phòng hỏi của nhân viên trong phòng. Giám đốc/Kế toán hỏi toàn công ty. |
| **Q-S5-03** | BA-S5 | Kênh hội thoại AI và nạp dữ liệu? | [x] Đã giải quyết | Thiết kế Sidebar riêng "AI Phân tích". Hỗ trợ tải văn bản quy chế lương nội bộ lên làm dữ liệu RAG. |

---

## 3. NHẬT KÝ PHIÊN (SESSION LOG)

### Phiên BA-S1: Khởi động & Làm rõ Module Phân quyền & Tổ chức
* **Ngày:** 2026-06-28
* **Tóm tắt nội dung:** Khởi tạo Master Plan, thống nhất cách làm việc cuốn gói và bắt đầu phỏng vấn thu thập yêu cầu cho Module Cơ cấu tổ chức & Phân quyền.

### Phiên BA-S2: Làm rõ Hồ sơ nhân viên & Cơ chế tính toán lương P1, P2
* **Ngày:** 2026-06-28
* **Tóm tắt nội dung:** Khảo sát yêu cầu Module Hồ sơ & Hợp đồng. Làm rõ cách thức tính toán Lương vị trí (P1) cộng dồn theo tiêu chuẩn và Lương năng lực (P2) đánh giá bán niên bởi cấp quản lý.

### Phiên BA-S3: Chấm công khuôn mặt, Cấu hình OT & Phân phối qua NATS
* **Ngày:** 2026-06-28
* **Tóm tắt nội dung:** Khảo sát yêu cầu Module Ca kíp & Chấm công. Thống nhất cơ chế check-in qua điện thoại đẩy tải qua NATS JetStream, quy tắc tính công ngày và chính sách tính lương OT.

### Phiên BA-S4: Thiết kế Động cơ AST, Tính toán Batch & Bảo mật duyệt chi
* **Ngày:** 2026-06-28
* **Tóm tắt nội dung:** Khảo sát yêu cầu Module Động cơ công thức & Tính lương Batch. Định hình cơ chế quản lý phiên bản công thức lương, luồng gửi duyệt điều chỉnh lương lên Trưởng phòng và xác thực chi tiền bằng nhận diện khuôn mặt Giám đốc.

### Phiên BA-S5: Tự động hóa AI Kiểm toán, Tra cứu RAG & Phân quyền Sidebar
* **Ngày:** 2026-06-28
* **Tóm tắt nội dung:** Khảo sát yêu cầu Module AI & Analytics. Định hình cơ chế chạy tự động AI Auditor báo cáo email/dashboard, thiết kế sidebar "AI Phân tích" phân quyền RLS 3 cấp độ và tải tài liệu RAG.
