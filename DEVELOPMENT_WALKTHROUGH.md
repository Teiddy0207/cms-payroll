# SỔ THEO DÕI TIẾN ĐỘ PHÁT TRIỂN PHẦN MỀM (DEVELOPMENT WALKTHROUGH)

Tài liệu này dùng để theo dõi tiến độ hoàn thiện hệ thống tính lương phân tán theo các phiên nghiệp vụ trong [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md).

---

## 📊 Bảng Tổng Hợp Trạng Thái Các Module

| Module Nghiệp Vụ | Phiên BA | Hiện Trạng Mã Nguồn (As-is) | Trạng Thái | Độ Ưu Tiên |
| :--- | :---: | :--- | :---: | :---: |
| **Cơ cấu tổ chức & Phân quyền (RBAC/RLS)** | `BA-S1` | Đã có sơ đồ phòng ban, xác thực JWT, ghi Audit log. | **Hoàn thành cơ bản** | High |
| **Hồ sơ nhân sự & Tính toán P1** | `BA-S2` | Đã có Hồ sơ, Hợp đồng, tính lương vị trí P1 ổn định. | **Hoàn thành** | High |
| **Chấm công khuôn mặt & Luồng NATS** | `BA-S3` | Chưa có kết nối NATS JetStream, chưa có API check-in hứng tải. | **Chưa thực hiện** | **High (Cần làm trước)** |
| **Động cơ công thức & Sắp xếp Topo** | `BA-S4` | Sắp xếp công thức đang bị fix cứng, chưa có giải thuật đồ thị Topo. | **Chưa thực hiện** | **High (Cần làm trước)** |
| **Lương hiệu quả P3 (Performance)** | `BA-S4` | Đang bị gán cứng = 0.0 trong code, chưa liên kết với động cơ công thức. | **Chưa thực hiện** | **High (Cần làm trước)** |
| **Phiếu điều chỉnh lương & Khóa sổ** | `BA-S4` | Chưa có API duyệt phiếu chỉnh sửa lương của kế toán. | **Chưa thực hiện** | Medium |
| **Xác thực chi tiền bằng Face ID Giám đốc** | `BA-S4` | Chưa có API so khớp khuôn mặt để xác thực chi tiền. | **Chưa thực hiện** | Medium |
| **AI Kiểm toán & RAG Chatbot** | `BA-S5` | Chưa có công cụ quét bất thường và RAG Chatbot. | **Chưa thực hiện** | Low |
| **Cải tiến Lương năng lực P2** | `BA-S2` | Đang là phép cộng đơn giản, cần nâng cấp 3 cấu phần ($P_{21}, P_{22}, P_{23}$). | **Để lại cuối cùng** | Low |

---

## 🗓️ Lộ Trình Triển Khai Chi Tiết (Roadmap)

### 🚀 Giai đoạn 1: Ổn định hệ thống & Luồng chấm công bất đồng bộ (Ưu tiên hiện tại)
* **Mục tiêu:** Tích hợp NATS JetStream để hứng tải sự kiện check-in, cài đặt giải thuật sắp xếp Topo động cho công thức lương và xử lý Lương hiệu quả P3.
* **Các task chính:**
  1. Thiết lập kết nối NATS JetStream và viết Queue Producer tại API Gateway.
  2. Viết Consumer bất đồng bộ lưu log vào DB, tích hợp lọc trùng sự kiện.
  3. Xây dựng đồ thị phụ thuộc biến (Dependency Graph) và thuật toán Topological Sort cho `Formula Engine`.
  4. Liên kết tính toán P3 từ động cơ công thức thay vì gán cứng 0.0.

### 🔒 Giai đoạn 2: Nghiệp vụ điều chỉnh & Bảo mật sinh trắc học
* **Mục tiêu:** Đưa quy trình phê duyệt sửa lương và duyệt chi bằng Face ID vào vận hành.
* **Các task chính:**
  1. API đề xuất và duyệt phiếu điều chỉnh lương `PayrollAdjustment`.
  2. Tích hợp thư viện hoặc giải thuật so khớp vector khuôn mặt Giám đốc khi phê duyệt chi lương.

### 🧠 Giai đoạn 3: Trí tuệ nhân tạo (AI) & Nâng cấp lương P2
* **Mục tiêu:** Hoàn thiện AI Auditor, chatbot nội bộ (RAG Chatbot tra cứu thông tin bảo mật RLS) và nâng cấp lương năng lực P2.
* **Các task chính:**
  1. AI Auditor quét dữ liệu lương bất thường.
  2. Chatbot RAG phân quyền RLS.
  3. Cải tiến P2 (Khung năng lực $P_{21}$, Chứng chỉ $P_{22}$, thâm niên $P_{23}$ và chặn trần).
