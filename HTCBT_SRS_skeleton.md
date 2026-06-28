# TÀI LIỆU ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS)
## HỆ THỐNG TÍNH LƯƠNG PHÂN TÁN (DISTRIBUTED PAYROLL PLATFORM)

| Phiên bản | Trạng thái | Ngày | Tác giả / Vai trò | Mô tả thay đổi |
| :--- | :--- | :--- | :--- | :--- |
| v0.1 | Dự thảo | 2026-06-28 | BA (Antigravity) | Khởi tạo tài liệu đặc tả, hoàn thiện Module Phân quyền & Tổ chức (Phiên BA-S1). |

---

# 1. Giới thiệu

## 1.1 Mục tiêu
Tài liệu đặc tả này xác định các yêu cầu chức năng, phi chức năng và thiết kế hệ thống cho **Hệ thống tính lương phân tán (Distributed Payroll Platform)**. Hệ thống giúp các doanh nghiệp tự động hóa hoàn toàn quy trình thu thập dữ liệu chấm công thời gian thực hiệu năng cao, định nghĩa công thức lương động theo mô hình 3P thông qua động cơ AST, và kiểm toán tự động bằng trí tuệ nhân tạo (AI).

## 1.2 Phạm vi
Hệ thống bao gồm các phân hệ cốt lõi:
1. Xác thực và Phân quyền mức dòng (RBAC & RLS).
2. Quản lý Hồ sơ nhân sự và Khung năng lực.
3. Quản lý ca kíp và Chấm công khuôn mặt thời gian thực (qua NATS JetStream).
4. Động cơ biên dịch công thức lương động (Expr AST Parser & Topological Sort).
5. Tính toán bảng lương hàng loạt quy mô lớn (Go Worker Pool & Redis Lock).
6. Kiểm toán bảng lương bằng AI và Trợ lý ảo tra cứu lương tự động (RAG Chatbot).

## 1.3 Định nghĩa, từ viết tắt và thuật ngữ

| Thuật ngữ | Định nghĩa trong hệ thống |
| :--- | :--- |
| **SaaS** | Software as a Service - Mô hình phân phối phần mềm dịch vụ đám mây. |
| **3P** | Mô hình lương dựa trên Vị trí (Position - P1), Năng lực (Person - P2) và Hiệu suất (Performance - P3). |
| **AST** | Abstract Syntax Tree - Cây cú pháp trừu tượng được sinh ra khi biên dịch chuỗi công thức lương. |
| **DAG** | Directed Acyclic Graph - Đồ thị có hướng không chu trình dùng để mô tả phụ thuộc giữa các công thức lương. |
| **RBAC** | Role-Based Access Control - Kiểm soát truy cập dựa trên vai trò của người dùng. |
| **RLS** | Row-Level Security - Bảo mật cấp dòng, giới hạn quyền xem dữ liệu bản ghi cụ thể theo bộ phận quản lý. |
| **NATS JetStream** | Message broker chịu tải cao được sử dụng để xếp hàng và điều tiết luồng sự kiện check-in chấm công. |
| **OLAP** | Online Analytical Processing - Hệ cơ sở dữ liệu phân tích (sử dụng ClickHouse để lưu trữ log chấm công). |
| **OLTP** | Online Transaction Processing - Hệ cơ sở dữ liệu giao dịch (sử dụng PostgreSQL để lưu dữ liệu nghiệp vụ chính). |

## 1.4 Tham khảo
* Đặc tả kiến trúc kỹ thuật ban đầu của hệ thống tính lương phân tán.
* Tiêu chuẩn quốc tế ISO/IEC/IEEE 29148:2018 về Kỹ thuật yêu cầu phần mềm.

## 1.5 Tổng quan
Tài liệu này được tổ chức thành 3 chương chính:
* **Chương 1:** Giới thiệu tổng quan dự án và các thuật ngữ kỹ thuật sử dụng.
* **Chương 2:** Mô tả tổng quan bối cảnh sản phẩm, vai trò tác nhân và các ràng buộc hệ thống.
* **Chương 3:** Chi tiết các yêu cầu chức năng (chia theo các module phát triển qua từng phiên BA) và yêu cầu phi chức năng.

---

# 2. Mô tả tổng quan

## 2.1 Bối cảnh sản phẩm
Doanh nghiệp quy mô vừa và lớn gặp khó khăn trong việc quản trị dữ liệu nhân sự phức tạp, tính toán lương Excel thủ công dễ sai sót và nghẽn hệ thống check-in chấm công vào giờ cao điểm. Nền tảng này giải quyết bài toán bằng cách chia nhỏ các thành phần xử lý (phân tán hóa), áp dụng hàng đợi thông điệp để đệm dữ liệu chấm công và dùng sức mạnh của AI để hỗ trợ kiểm soát chất lượng dữ liệu.

## 2.2 Chức năng sản phẩm
Sản phẩm được chia thành các nhóm chức năng lớn:
1. **Module Phân quyền & Tổ chức:** Định nghĩa cơ cấu và phân quyền truy cập.
2. **Module Nhân sự & Hợp đồng:** Quản lý thông tin hồ sơ và các chỉ số lương cơ sở.
3. **Module Ca kíp & Chấm công:** Lập lịch làm việc và hứng dữ liệu chấm công thời gian thực.
4. **Module Động cơ Công thức:** Cho phép người dùng viết công thức tính lương động trên giao diện UI.
5. **Module Tính lương hàng loạt:** Động cơ chạy song song tính lương quy mô lớn.
6. **Module AI Kiểm toán & Assistant:** Rà soát sai lệch và trả lời khiếu nại của nhân viên.

## 2.3 Đặc điểm người dùng

| Vai trò người dùng | Mô tả quyền hạn trong hệ thống |
| :--- | :--- |
| **Admin** | Quản trị toàn bộ hệ thống, thiết lập cấu hình cơ sở và giám sát lịch sử log hệ thống. |
| **Giám đốc (Director)** | Được quyền xem báo cáo dashboard phân tích và bảng lương tổng hợp của toàn công ty. Không có quyền sửa đổi. |
| **Kế toán (Accountant)** | Được quyền cấu hình công thức lương, chạy job tính lương hàng loạt, điều chỉnh lương thực nhận và phê duyệt disbursement. |
| **Trưởng phòng (Head)** | Quản lý nhân sự trong phòng ban mình trực thuộc, phê duyệt đơn từ (nghỉ phép, tăng ca) và xem lương của nhân viên thuộc phòng ban mình. |
| **Nhân viên (Employee)** | Xem bảng công cá nhân, tạo đơn từ, tra cứu phiếu lương (payslip) cá nhân và tương tác với AI Assistant. |

## 2.4 Ràng buộc chung
* Hệ thống backend bắt buộc phát triển bằng Golang để tối ưu hiệu năng tính toán.
* Giao diện frontend phát triển bằng Next.js hỗ trợ vẽ biểu đồ công thức trực quan (React Flow).
* Dữ liệu chấm công thô bắt buộc lưu trữ tại ClickHouse để phục vụ phân tích tải cao.

## 2.5 Giả định và phụ thuộc
* Giả định các thiết bị nhận diện khuôn mặt chấm công (Camera AI) có khả năng trích xuất vector embedding 512 chiều trước khi gửi dữ liệu về Gateway.

---

# 3. Yêu cầu đặc tả

## 3.1 Yêu cầu giao diện ngoài
*(Sẽ được cập nhật chi tiết)*

---

## 3.2 Yêu cầu chức năng

| Số thứ tự | Mã yêu cầu | Tên chức năng | Mô tả ngắn |
| :--- | :--- | :--- | :--- |
| 1 | **UC-AUTH-01** | Đăng nhập hệ thống | Xác thực tài khoản người dùng qua Email/Password hoặc Keycloak SSO. |
| 2 | **UC-ORG-01**  | Quản lý phòng ban | Thiết lập sơ đồ tổ chức phòng ban trực thuộc công ty. |
| 3 | **UC-ORG-02**  | Gán chức danh và trưởng phòng | Quản lý chức danh nghề nghiệp và gán người quản lý cho bộ phận. |
| 4 | **UC-AUDIT-01**| Tra cứu lịch sử kiểm toán | Ghi nhận và hiển thị vết hoạt động của người dùng hệ thống. |
| 5 | **UC-EMP-01**  | Quản lý hồ sơ & Import Excel | Quản lý danh sách nhân sự và nạp dữ liệu hàng loạt. |
| 6 | **UC-P1-01**   | Tính toán Lương P1 | Tính lương vị trí cộng dồn các chỉ mục tiêu chuẩn. |
| 7 | **UC-P2-01**   | Chấm điểm & Tính Lương P2 | Quản lý đánh giá năng lực bán niên và quy đổi điểm ra lương P2. |
| 8 | **UC-ATT-01**  | Chấm công di động qua NATS | Thu thập dữ liệu check-in thô qua queue NATS JetStream. |
| 9 | **UC-CALC-01** | Tính toán ngày công tự động | Tổng hợp dữ liệu check-in thô thành công ngày (1.0, 0.5 công). |
| 10| **UC-EXC-01**  | Quản lý tờ trình & Duyệt OT | Cho phép nhân sự viết tờ trình giải trình và xin phê duyệt OT. |
| 11| **UC-FORM-01** | Quản lý công thức lương | Thiết lập công thức động, hỗ trợ lưu vết phiên bản công thức. |
| 12| **UC-PIPE-01** | Tính lương Batch tự động | Chạy tính lương bất đồng bộ cho toàn bộ nhân sự qua Go Worker Pool. |
| 13| **UC-ADJ-01**  | Điều chỉnh lương nháp | Kế toán đề xuất sửa lương thủ công gửi Trưởng phòng phê duyệt. |
| 14| **UC-DISB-01** | Khóa sổ & Xác thực chi trả | Khóa vĩnh viễn bảng lương tháng và quét Face ID Giám đốc phê duyệt chi trả. |
| 15| **UC-AI-01**   | Kiểm toán bảng lương AI | Tự động phân tích tìm lỗi và cảnh báo qua email/dashboard. |
| 16| **UC-AI-02**   | Sidebar AI Phân tích | Hỏi đáp phân tích lương, ngày công kết hợp RAG học quy chế nội bộ. |

---

### 3.2.1 Module 1: Phân quyền & Tổ chức (Phiên BA-S1)

#### 3.2.1.1 Chức năng Xác thực người dùng (UC-AUTH-01)
* **Mô tả:** Cho phép người dùng đăng nhập hệ thống để lấy Access Token truy cập các tài nguyên API.
* **Dữ liệu đầu vào:**
  * Phương thức 1: Email, Password.
  * Phương thức 2: Mã Authorization Code từ Keycloak SSO.
* **Xử lý:**
  1. Đối với Email/Password: Kiểm tra sự tồn tại tài khoản trong PostgreSQL, băm mật khẩu so khớp với mật khẩu đã lưu.
  2. Đối với Keycloak SSO: Xác thực chữ ký số của token phát ra từ Keycloak, tự động đồng bộ tài khoản nếu là lần đầu đăng nhập.
  3. Cấp cặp token: JWT Access Token (hết hạn sau 15 phút) và Refresh Token (hết hạn sau 7 ngày).
  4. Ghi nhận trạng thái phiên hoạt động lên Redis Cache.
* **Dữ liệu đầu ra:** Cặp JWT Token, thông tin vai trò người dùng (Roles) và phòng ban trực thuộc.

#### 3.2.1.2 Chức năng Quản lý cơ cấu phòng ban (UC-ORG-01)
* **Mô tả:** Thiết lập sơ đồ phòng ban 2 cấp (Công ty -> Phòng ban).
* **Dữ liệu đầu vào:** Tên phòng ban, mã phòng ban, mô tả.
* **Xử lý:**
  1. Cho phép tạo mới, cập nhật hoặc xóa phòng ban (chỉ cho phép xóa khi phòng ban không còn nhân viên hoạt động).
  2. Áp dụng quy tắc gán **tối đa 01 Trưởng phòng** cho mỗi phòng ban. Khi gán Trưởng phòng mới, hệ thống tự động gỡ vai trò Trưởng phòng của nhân sự cũ trong phòng ban đó.
* **Dữ liệu đầu ra:** Danh sách cây phòng ban được cập nhật trong PostgreSQL.

#### 3.2.1.3 Phân quyền mức dòng dữ liệu (Row-Level Security - RLS)
* **Mô tả:** Cơ chế tự động chèn bộ lọc phòng ban vào tất cả các câu lệnh SQL/ClickHouse đối với tài khoản cấp Trưởng phòng.
* **Xử lý:**
  1. Khi người dùng có vai trò `Trưởng phòng` gửi request truy vấn dữ liệu (Bảng công, Đơn từ nghỉ phép, Phiếu lương nhân viên).
  2. Hệ thống đọc thông tin phòng ban trực thuộc của Trưởng phòng từ Context Token.
  3. Backend tự động append điều kiện: `WHERE employee.department_id = context.department_id` vào câu lệnh truy vấn cơ sở dữ liệu.
  4. Nghiêm cấm hiển thị hoặc xử lý phê duyệt đơn từ đối với nhân viên có `department_id` khác với Trưởng phòng đó.
* **Dữ liệu đầu ra:** Dữ liệu trả về được lọc đúng phòng ban quản lý.

#### 3.2.1.4 Chức năng Ghi vết kiểm toán (UC-AUDIT-01)
* **Mô tả:** Ghi nhận lại các hành vi nhạy cảm trong hệ thống chấm công - tính lương.
* **Xử lý:**
  1. Hệ thống tự động ghi nhật ký (IP, Thời gian, Tác nhân, Hành động, Dữ liệu cũ, Dữ liệu mới) đối với các hành vi: Đăng nhập/Đăng xuất, Chạy tính toán lương tháng, Sửa đổi bảng công, Thay đổi tiền lương thực nhận.
  2. Đối với hành vi **Sửa đổi bảng công** và **Thay đổi tiền lương thực nhận**: Yêu cầu API nhận kèm tham số `reason` (Lý do điều chỉnh). Nếu tham số này rỗng hoặc độ dài dưới 10 ký tự, hệ thống từ chối cập nhật và trả về mã lỗi `400 Bad Request`.
* **Dữ liệu đầu ra:** Bản ghi nhật ký hoạt động được lưu vào PostgreSQL.

### 3.2.2 Module 2: Hồ sơ nhân sự & Tính toán P1, P2 (Phiên BA-S2)

#### 3.2.2.1 Chức năng Quản lý Hồ sơ nhân sự & Nhập liệu Excel (UC-EMP-01)
* **Mô tả:** Cho phép HR quản lý danh sách hồ sơ nhân viên và nạp dữ liệu hàng loạt từ file Excel.
* **Dữ liệu đầu vào:**
  * Thao tác thủ công: Hồ sơ nhân viên (Mã NV, Họ tên, Ngày sinh, Phòng ban, Email, Số điện thoại).
  * Nhập hàng loạt: File Excel và Cấu hình ánh xạ cột (Mapping Schema).
* **Xử lý:**
  1. Tạo, sửa, chuyển trạng thái "Thôi việc" cho nhân sự.
  2. Đối với chức năng Import Excel: HR chọn tải lên một file Excel bất kỳ. Giao diện frontend hiển thị các cột trong file Excel và cho phép HR chọn ánh xạ (ví dụ: Cột A khớp với trường `Mã nhân viên`, Cột B khớp với trường `Họ tên`). Backend Go nhận file stream cùng schema mapping JSON để parse và bulk insert dữ liệu vào DB PostgreSQL.
* **Dữ liệu đầu ra:** Danh sách hồ sơ nhân sự được cập nhật trong DB.

#### 3.2.2.2 Chức năng Cấu hình & Tính toán Lương vị trí P1 (UC-P1-01)
* **Mô tả:** Tính toán lương cơ sở P1 (`P1_BASE`) tự động dựa trên mức lương nền của chức danh công việc cộng dồn theo các tiêu chuẩn.
* **Dữ liệu đầu vào:** Chức danh (Role), mức lương nền chức danh, danh sách mã tiêu chuẩn (chứng chỉ, ngành nghề).
* **Xử lý:**
  1. HR cấu hình bảng tiêu chuẩn `job_standards` gồm mã tiêu chuẩn (ví dụ: `DEGREE_BACHELOR`, `CERT_CPA`, `TECH_GOLANG`) kèm theo số tiền phụ cấp được cộng thêm.
  2. Khi tính lương, hệ thống truy vấn hợp đồng của nhân viên để lấy lương nền chức danh (`Position_Base_Rate`) và danh sách các tiêu chuẩn bắt buộc thuộc vị trí công việc đó.
  3. Hệ thống cộng dồn: `P1_BASE = Position_Base_Rate + SUM(Standard_Value)`.
* **Dữ liệu đầu ra:** Số tiền `P1_BASE` được lưu vào RAM context để chuyển sang Dynamic Formula Engine.

#### 3.2.2.3 Chức năng Đánh giá & Tính toán Lương năng lực P2 (UC-P2-01)
* **Mô tả:** Quản lý chấm điểm năng lực định kỳ bán niên và quy đổi điểm năng lực thành tiền lương P2 (`P2_COMPETENCY`).
* **Dữ liệu đầu vào:** Điểm đánh giá (1-100) của từng mục năng lực do Trưởng phòng chấm cho nhân viên phòng mình.
* **Xử lý:**
  1. HR cấu hình danh sách năng lực (ví dụ: `Tư duy logic`, `Kỹ năng viết code`, `Trình độ ngoại ngữ`) kèm theo trọng số (Weight từ 1 đến 5) tương ứng với từng vị trí công việc.
  2. Định kỳ **6 tháng một lần (đánh giá bán niên)**, Trưởng phòng đăng nhập để chấm điểm (Scale từ 1 đến 100) cho từng năng lực của các nhân viên trực thuộc phòng mình quản lý (Áp dụng RLS).
  3. Khi chạy tính lương, hệ thống lấy kết quả đánh giá mới nhất của nhân viên và áp dụng công thức: `P2_COMPETENCY = SUM(Competency_Score * Competency_Weight) * Company_Point_Rate`. Trong đó, `Company_Point_Rate` là đơn giá quy đổi 1 điểm thành VND (do doanh nghiệp cấu hình).
* **Dữ liệu đầu ra:** Số tiền `P2_COMPETENCY` lưu vào RAM context tính toán.

---

## 3.3 Đối tượng (Entity Schemas)

Hệ thống quản lý các thực thể dữ liệu chính sau trong PostgreSQL:

### 3.3.1 Thực thể Nhân viên (Employee)
* `id` (UUID, Primary Key)
* `employee_code` (VARCHAR, Unique): Mã nhân viên.
* `full_name` (VARCHAR): Họ và tên.
* `email` (VARCHAR, Unique): Email đăng nhập.
* `password_hash` (VARCHAR): Mật khẩu băm.
* `department_id` (UUID, Foreign Key): Phòng ban trực thuộc.
* `role` (VARCHAR): Vai trò mặc định (Admin, Giám đốc, Kế toán, Trưởng phòng, Nhân viên).
* `status` (VARCHAR): Trạng thái hoạt động (ACTIVE, INACTIVE).

### 3.3.2 Thực thể Hợp đồng lao động (Contract)
* `id` (UUID, Primary Key)
* `employee_id` (UUID, Foreign Key): Liên kết nhân viên.
* `contract_code` (VARCHAR, Unique): Số hợp đồng.
* `position_base_rate` (DECIMAL): Mức lương sàn chức danh.
* `start_date` (DATE): Ngày bắt đầu hiệu lực.
* `end_date` (DATE): Ngày hết hiệu lực (nếu có).
* `status` (VARCHAR): Trạng thái hợp đồng (ACTIVE, EXPIRED, TERMINATED).

### 3.3.3 Thực thể Tiêu chuẩn công việc (JobStandard)
* `id` (UUID, Primary Key)
* `standard_code` (VARCHAR, Unique): Mã tiêu chuẩn (ví dụ: `CERT_CPA`, `DEGREE_BACHELOR`).
* `name` (VARCHAR): Tên tiêu chuẩn.
* `allowance_value` (DECIMAL): Số tiền được cộng thêm vào lương vị trí P1.

### 3.3.4 Thực thể Danh mục năng lực (Competency)
* `id` (UUID, Primary Key)
* `competency_code` (VARCHAR, Unique): Mã năng lực (ví dụ: `COMP_LOGIC`, `COMP_CODING`).
* `name` (VARCHAR): Tên năng lực.
* `description` (TEXT): Mô tả tiêu chí đánh giá.

### 3.3.5 Thực thể Đánh giá năng lực (CompetencyEvaluation)
* `id` (UUID, Primary Key)
* `employee_id` (UUID, Foreign Key): Nhân viên được đánh giá.
* `evaluator_id` (UUID, Foreign Key): Trưởng phòng thực hiện đánh giá.
* `evaluation_period` (VARCHAR): Kỳ đánh giá (ví dụ: `2026-H1`, `2026-H2`).
* `competency_id` (UUID, Foreign Key): Liên kết năng lực đánh giá.
* `score` (INTEGER): Điểm chấm (từ 1 đến 100).
* `weight` (INTEGER): Trọng số yêu cầu (từ 1 đến 5).
* `created_at` (TIMESTAMP): Thời điểm đánh giá.

### 3.2.3 Module 3: Ca kíp & Chấm công khuôn mặt (Phiên BA-S3)

#### 3.2.3.1 Chức năng Thiết lập ca và Chính sách OT (UC-ATT-01)
* **Mô tả:** Cho phép quản trị viên thiết lập ca làm việc tiêu chuẩn và các hệ số làm thêm giờ (OT).
* **Dữ liệu đầu vào:** Tên ca, giờ bắt đầu, giờ kết thúc, thời gian ân hạn đi muộn, các hệ số nhân OT.
* **Xử lý:**
  1. Cấu hình Ca hành chính tiêu chuẩn: Giờ vào 8:30 (cho phép muộn đến 9:00), Giờ ra: 17:30.
  2. Cấu hình hệ số OT: OT đêm nhân 1.5, OT ngày nghỉ/lễ nhân 2.0.
  3. Áp dụng quy tắc kiểm soát OT: Người lao động làm thêm sau 17:30 chỉ được tính OT nếu có đăng ký và được Trưởng phòng phê duyệt. Nếu không có đơn, giờ làm thêm không được ghi nhận.
* **Dữ liệu đầu ra:** Cấu hình ca làm việc và bảng hệ số OT được lưu trong DB.

#### 3.2.3.2 Chức năng Thu thập chấm công qua NATS JetStream (UC-ATT-02)
* **Mô tả:** Hứng nhận dữ liệu check-in thời gian thực từ Mobile App của nhân viên gửi về.
* **Dữ liệu đầu vào:** Thông tin sự kiện check-in `{employee_id, timestamp, location_gps, device_id}`.
* **Xử lý:**
  1. Ứng dụng di động tự chạy nhận diện khuôn mặt cục bộ (Edge AI), so khớp và trích xuất đúng `employee_id` gửi về API Gateway.
  2. API Gateway thực hiện kiểm tra Rate Limit (Redis Token Bucket) và đẩy sự kiện check-in thô vào NATS JetStream chủ đề `timekeeping.checkin`. Gateway phản hồi lập tức `202 Accepted` về điện thoại.
  3. Go Attendance Workers đăng ký lắng nghe NATS, thực hiện lọc trùng lặp bản ghi (nếu nhân viên bấm quẹt liên tục nhiều lần trong vòng 5 phút thì chỉ lấy lượt đầu tiên).
  4. Worker ghi nhận bản ghi vào cơ sở dữ liệu phân tích ClickHouse (để lưu trữ lịch sử chấm công quy mô lớn) và PostgreSQL (đồng bộ trạng thái).
* **Dữ liệu đầu ra:** Log chấm công được lưu vết an toàn trong ClickHouse.

#### 3.2.3.3 Chức năng Động cơ tính công ngày tự động & Quy trình bù công (UC-CALC-01 / UC-EXC-01)
* **Mô tả:** Tổng hợp các lượt quẹt thẻ thô trong ngày thành ngày công thực tế của nhân viên.
* **Xử lý:**
  1. Hàng đêm lúc 23:59, hệ thống kích hoạt Batch Worker tự động quét log chấm công trong ngày của nhân viên.
  2. Quy tắc 1.0 công: Có ít nhất 02 lượt quẹt thẻ, lượt đầu tiên trước 09:00 và lượt cuối cùng sau 17:30.
  3. Quy tắc 0.5 công: Quẹt đầu tiên trước 09:00 và quẹt cuối cùng nằm trong khoảng 12:00 - 14:00 (về nghỉ nửa ngày).
  4. Quy trình xử lý tờ trình bù công: Nhân viên đi muộn (check-in sau 09:00) hoặc về sớm (trước 17:30 dưới 2 tiếng) có quyền làm "Tờ trình giải trình" gửi lên hệ thống. Nếu Trưởng phòng duyệt tờ trình này (áp dụng RLS), ngày công đó sẽ được ghi đè và tính đủ 1.0 công bình thường.
* **Dữ liệu đầu ra:** Bảng tổng hợp công ngày (`daily_attendance_sheets`) lưu vào PostgreSQL.

---

## 3.3 Đối tượng (Entity Schemas)

Hệ thống quản lý các thực thể dữ liệu chính sau trong PostgreSQL:

### 3.3.1 Thực thể Nhân viên (Employee)
*(Đã định nghĩa ở Chương 3.3.1)*

### 3.3.2 Thực thể Hợp đồng lao động (Contract)
*(Đã định nghĩa ở Chương 3.3.2)*

### 3.3.3 Thực thể Tiêu chuẩn công việc (JobStandard)
*(Đã định nghĩa ở Chương 3.3.3)*

### 3.3.4 Thực thể Danh mục năng lực (Competency)
*(Đã định nghĩa ở Chương 3.3.4)*

### 3.3.5 Thực thể Đánh giá năng lực (CompetencyEvaluation)
*(Đã định nghĩa ở Chương 3.3.5)*

### 3.3.6 Thực thể Log chấm công thô (AttendanceLog - Lưu trữ tại ClickHouse)
* `id` (UUID, Primary Key)
* `employee_code` (VARCHAR): Mã nhân viên chấm công.
* `timestamp` (DateTime): Thời điểm quẹt thẻ.
* `location_gps` (Point/String): Tọa độ chấm công (vĩ độ, kinh độ).
* `device_id` (VARCHAR): ID của điện thoại thực hiện chấm.

### 3.3.7 Thực thể Bảng công ngày (DailyAttendanceSheet - Lưu trữ tại PostgreSQL)
* `id` (UUID, Primary Key)
* `employee_id` (UUID, Foreign Key)
* `date` (DATE): Ngày chấm công.
* `check_in` (TIMESTAMP, Nullable): Lượt quẹt vào đầu tiên được ghi nhận.
* `check_out` (TIMESTAMP, Nullable): Lượt quẹt ra cuối cùng được ghi nhận.
* `actual_work_day` (DECIMAL): Số ngày công tổng hợp (0.0, 0.5, 1.0).
* `ot_hours` (DECIMAL): Số giờ OT được phê duyệt.
* `status` (VARCHAR): Trạng thái công (NORMAL, LATE, EARLY, ABSENT).

### 3.3.8 Thực thể Tờ trình giải trình bù công (ExplanationRequest)
* `id` (UUID, Primary Key)
* `employee_id` (UUID, Foreign Key)
* `date` (DATE): Ngày cần giải trình bù công.
* `reason` (TEXT): Lý do đi muộn/về sớm.
* `approved_by` (UUID, Foreign Key): Trưởng phòng duyệt.
* `status` (VARCHAR): Trạng thái (PENDING, APPROVED, REJECTED).

### 3.3.9 Thực thể Đăng ký làm thêm giờ (OTRequest)
* `id` (UUID, Primary Key)
* `employee_id` (UUID, Foreign Key)
* `date` (DATE): Ngày làm thêm giờ.
* `hours_requested` (DECIMAL): Số giờ đăng ký OT.
* `is_night_ot` (BOOLEAN): Làm thêm ca đêm.
* `is_holiday_ot` (BOOLEAN): Làm thêm ngày nghỉ/lễ.
* `approved_by` (UUID, Foreign Key): Trưởng phòng duyệt.
* `status` (VARCHAR): Trạng thái (PENDING, APPROVED, REJECTED).

### 3.2.4 Module 4: Động cơ Công thức & Tính lương Batch (Phiên BA-S4)

#### 3.2.4.1 Chức năng Cấu hình công thức & Quản lý phiên bản (UC-FORM-01)
* **Mô tả:** Cho phép Kế toán thiết lập công thức tính lương động thông qua biểu thức toán học và quản lý lịch sử thay đổi phiên bản.
* **Dữ liệu đầu vào:** Tên biến lương, biểu thức tính toán (Expression string), khoảng thời gian hiệu lực (`start_date`, `end_date`).
* **Xử lý:**
  1. Khi người dùng lưu công thức mới, Backend Go chạy thư viện parser (Expr) để kiểm tra lỗi cú pháp (Syntax checking).
  2. Tạo bản ghi mới trong bảng `payroll_formulas` lưu trữ thời gian hiệu lực của công thức. Các công thức cũ hết hiệu lực được gán ngày kết thúc (`end_date`).
  3. Khi thực hiện tính toán lương hoặc tính toán lại hồi tố (Retroactive calculation) của tháng trong quá khứ, hệ thống tự động lọc các công thức có khoảng hiệu lực bao phủ thời điểm cần tính.
* **Dữ liệu đầu ra:** Cấu hình công thức lương được lưu vết phiên bản trong PostgreSQL.

#### 3.2.4.2 Chức năng Chạy tính lương Batch bất đồng bộ (UC-PIPE-01)
* **Mô tả:** Kích hoạt Worker tính toán song song bảng lương hàng loạt cho toàn bộ nhân sự công ty.
* **Dữ liệu đầu vào:** Chu kỳ tính lương (Tháng/Năm).
* **Xử lý:**
  1. Kế toán nhấn "⚡ Tính Lương", Backend Go sinh mã `Job_ID` và tạo job chạy ngầm, phản hồi ngay lập tức về UI để người dùng không cần chờ đợi.
  2. Tiến độ công việc (% Progress) được cập nhật thời gian thực lên màn hình qua kết nối WebSocket.
  3. Backend sử dụng chiến lược **Eager Loading**: Query gộp toàn bộ hồ sơ nhân sự, bảng công ngày, chỉ số P1, P2 của chu kỳ đó lên RAM (lọc theo lô 1,000 người dùng toán tử `IN` SQL) để tạo Context Map, triệt tiêu lỗi N+1 Query.
  4. Hệ thống chạy thuật toán Topological Sort sắp xếp thứ tự tính các cột lương, sau đó phân chia tính toán song song qua **Go Goroutines (Worker Pool)**. Mỗi bản ghi nhân viên khi tính toán được bảo vệ bởi Redis Distributed Lock tránh xung đột ghi đè dữ liệu.
* **Dữ liệu đầu ra:** Kết quả lương tạm tính được lưu vào bảng nháp `payroll_records` ở trạng thái `DRAFT`.

#### 3.2.4.3 Chức năng Điều chỉnh lương nháp & Khóa sổ (UC-ADJ-01 / UC-DISB-01)
* **Mô tả:** Kế toán điều chỉnh lương thủ công và tiến hành khóa sổ bảng lương tháng.
* **Xử lý:**
  1. Kế toán xem bảng lương nháp dưới dạng bảng tính (Spreadsheet). Cho phép sửa đổi số tiền của một nhân sự.
  2. Việc sửa đổi tạo ra một yêu cầu `PayrollAdjustment` yêu cầu nhập lý do. Trưởng phòng của nhân sự đó duyệt phiếu yêu cầu này (áp dụng RLS) thì số liệu lương mới chính thức được cập nhật ghi đè.
  3. Khi bảng lương đạt độ chính xác, Kế toán chọn "Khóa sổ". Trạng thái chu kỳ chuyển sang `LOCKED`. Hệ thống kích hoạt cơ chế khóa bảo mật: Cấm mọi thao tác sửa đổi bảng công, đơn từ hay bất kỳ trường dữ liệu lương nào của tháng này.
* **Dữ liệu đầu ra:** Trạng thái bảng lương cập nhật thành `LOCKED`.

#### 3.2.4.4 Chức năng Xác thực chi trả lương bằng Face ID Giám đốc (UC-DISB-02)
* **Mô tả:** Giám đốc phê duyệt chuyển khoản chi tiền lương ảo/thực tế thông qua xác thực nhận diện khuôn mặt sinh trắc học.
* **Dữ liệu đầu vào:** Luồng hình ảnh camera quét trực tiếp khuôn mặt Giám đốc trên thiết bị.
* **Xử lý:**
  1. Giám đốc đăng nhập tài khoản quyền Director, nhấn nút "Phê duyệt chi trả" cho chu kỳ lương đã khóa.
  2. Hệ thống mở camera, trích xuất vector khuôn mặt trực tiếp và thực hiện so khớp với Vector Face ID gốc của Giám đốc được lưu trữ bảo mật trong database.
  3. Nếu độ tương đồng đạt > 95%, hệ thống ký số xác nhận lệnh chi trả, đổi trạng thái chu kỳ lương sang `DISBURSED` (Đã chi trả) và gửi thông báo payslip tự động cho toàn bộ nhân sự.
* **Dữ liệu đầu ra:** Phiếu lương được xác thực thành công, trạng thái chuyển sang `DISBURSED`.

---

## 3.3 Đối tượng (Entity Schemas)

Hệ thống quản lý các thực thể dữ liệu chính sau trong PostgreSQL:

### 3.3.1 Thực thể Nhân viên (Employee)
*(Đã định nghĩa ở Chương 3.3.1)*

### 3.3.2 Thực thể Hợp đồng lao động (Contract)
*(Đã định nghĩa ở Chương 3.3.2)*

### 3.3.3 Thực thể Tiêu chuẩn công việc (JobStandard)
*(Đã định nghĩa ở Chương 3.3.3)*

### 3.3.4 Thực thể Danh mục năng lực (Competency)
*(Đã định nghĩa ở Chương 3.3.4)*

### 3.3.5 Thực thể Đánh giá năng lực (CompetencyEvaluation)
*(Đã định nghĩa ở Chương 3.3.5)*

### 3.3.6 Thực thể Log chấm công thô (AttendanceLog - Lưu trữ tại ClickHouse)
*(Đã định nghĩa ở Chương 3.3.6)*

### 3.3.7 Thực thể Bảng công ngày (DailyAttendanceSheet - Lưu trữ tại PostgreSQL)
*(Đã định nghĩa ở Chương 3.3.7)*

### 3.3.8 Thực thể Tờ trình giải trình bù công (ExplanationRequest)
*(Đã định nghĩa ở Chương 3.3.8)*

### 3.3.9 Thực thể Đăng ký làm thêm giờ (OTRequest)
*(Đã định nghĩa ở Chương 3.3.9)*

### 3.3.10 Thực thể Công thức lương (PayrollFormula)
* `id` (UUID, Primary Key)
* `variable_name` (VARCHAR, Unique): Tên biến lương (ví dụ: `NET_SALARY`).
* `expression` (TEXT): Chuỗi công thức toán học (ví dụ: `GROSS_SALARY - TAX`).
* `start_date` (DATE): Ngày bắt đầu có hiệu lực.
* `end_date` (DATE, Nullable): Ngày hết hiệu lực.
* `description` (TEXT): Mô tả tác dụng công thức.

### 3.3.11 Thực thể Chu kỳ tính lương (PayrollPeriod)
* `id` (UUID, Primary Key)
* `month` (INTEGER): Tháng tính lương.
* `year` (INTEGER): Năm tính lương.
* `status` (VARCHAR): Trạng thái chu kỳ (`DRAFT`, `CALCULATED`, `LOCKED`, `DISBURSED`).
* `locked_at` (TIMESTAMP, Nullable): Thời điểm khóa sổ.
* `disbursed_at` (TIMESTAMP, Nullable): Thời điểm chi trả.

### 3.3.12 Thực thể Bảng ghi lương chi tiết (PayrollRecord)
* `id` (UUID, Primary Key)
* `period_id` (UUID, Foreign Key): Liên kết chu kỳ lương.
* `employee_id` (UUID, Foreign Key): Liên kết nhân viên.
* `p1_value` (DECIMAL): Mức lương P1 thực nhận của tháng.
* `p2_value` (DECIMAL): Mức lương P2 thực nhận của tháng.
* `p3_value` (DECIMAL): Mức lương P3 thực nhận của tháng.
* `gross_salary` (DECIMAL): Tổng thu nhập chịu thuế.
* `tax` (DECIMAL): Tiền thuế TNCN.
* `net_salary` (DECIMAL): Tiền thực nhận chuyển khoản.
* `status` (VARCHAR): Trạng thái dòng lương (`DRAFT`, `ADJUSTING`, `APPROVED`).

### 3.3.13 Thực thể Yêu cầu điều chỉnh lương (PayrollAdjustment)
* `id` (UUID, Primary Key)
* `record_id` (UUID, Foreign Key): Liên kết dòng lương cần điều chỉnh.
* `field_modified` (VARCHAR): Tên cột được sửa (ví dụ: `net_salary`).
* `old_value` (DECIMAL): Số tiền cũ.
* `new_value` (DECIMAL): Số tiền mới.
* `reason` (TEXT): Lý do kế toán chỉnh sửa số liệu.
* `requested_by` (UUID, Foreign Key): Kế toán đề xuất sửa.
* `approved_by` (UUID, Foreign Key, Nullable): Trưởng phòng duyệt.
* `status` (VARCHAR): Trạng thái phê duyệt (`PENDING`, `APPROVED`, `REJECTED`).

### 3.2.5 Module 5: AI Kiểm toán & Phân tích RAG (Phiên BA-S5)

#### 3.2.5.1 Chức năng Kiểm toán bảng lương bằng AI (UC-AI-01)
* **Mô tả:** Tự động rà soát tìm lỗi bảng lương thô dựa trên lịch sử dữ liệu và trí tuệ nhân tạo.
* **Dữ liệu đầu vào:** Bảng lương chi tiết ở trạng thái `DRAFT` vừa tính xong.
* **Xử lý:**
  1. AI Auditor khởi động tự động khi Job tính toán hoàn tất.
  2. Worker Go tính toán độ lệch chuẩn của các chỉ mục lương (Giờ làm thêm, Thưởng KPI, Tiền thuế) so với baseline lịch sử $N$ tháng của nhân viên (mặc định $N = 6$, cho phép HR tùy chỉnh giá trị $N$).
  3. Với các dòng bản ghi vượt ngưỡng độ lệch an toàn ($3\sigma$), hệ thống đóng gói context dưới dạng JSON (bao gồm thông tin lịch sử của nhân viên và chi tiết dòng lương hiện tại).
  4. Gửi context kèm prompt hướng dẫn lên LLM API (OpenAI/DeepSeek) để phân tích hành vi và giải nghĩa bất thường.
  5. Xuất các bản ghi bị gắn cờ cảnh báo kèm lý do phân tích chi tiết của AI lên Dashboard kiểm toán của Kế toán, đồng thời gửi email thông báo tự động.
* **Dữ liệu đầu ra:** Danh sách anomalies được gắn tag cảnh báo và ghi nhận vào cơ sở dữ liệu ClickHouse.

#### 3.2.5.2 Chức năng Thanh AI Phân tích & Kiến trúc RAG bảo mật (UC-AI-02)
* **Mô tả:** Tích hợp Sidebar riêng biệt cho phép người dùng hỏi đáp phân tích dữ liệu công/lương cá nhân hoặc phòng ban kết hợp đối chiếu quy chế lương.
* **Dữ liệu đầu vào:** Câu hỏi tự nhiên của người dùng, file tài liệu quy chế lương do HR tải lên.
* **Xử lý:**
  1. **Nạp tài liệu RAG:** Khi HR tải lên các tài liệu quy chế lương (`.pdf`, `.docx`), hệ thống chạy thư viện bóc tách text, cắt nhỏ thành các chunk văn bản, chạy mô hình embedding và lưu vào cơ sở dữ liệu Vector.
  2. **Tiếp nhận câu hỏi và áp dụng RLS:** Khi người dùng đặt câu hỏi trên Sidebar "AI Phân tích", hệ thống thực hiện xác thực và phân quyền truy vấn:
     * *Nhân viên:* Hệ thống chỉ cho phép truy vấn dữ liệu chấm công, bảng công và phiếu lương của chính nhân viên đó.
     * *Trưởng phòng:* Cho phép truy vấn dữ liệu cá nhân và dữ liệu tổng hợp của toàn bộ nhân viên thuộc phòng ban mình quản lý.
     * *Giám đốc/Kế toán:* Được phép truy vấn dữ liệu tổng hợp của toàn công ty.
  3. **Truy xuất & Trả lời:** Hệ thống tìm kiếm ngữ cảnh quy chế phù hợp trong Vector DB kết hợp với số liệu lương thô tương ứng từ DB PostgreSQL/ClickHouse để làm context gửi LLM sinh câu trả lời tiếng Việt chính xác.
* **Dữ liệu đầu ra:** Câu trả lời tự nhiên của AI giải thích cặn kẽ số liệu hiển thị trên Sidebar.

---

## 3.3 Đối tượng (Entity Schemas)

Hệ thống quản lý các thực thể dữ liệu chính sau trong PostgreSQL:

### 3.3.1 Thực thể Nhân viên (Employee)
*(Đã định nghĩa ở Chương 3.3.1)*

### 3.3.2 Thực thể Hợp đồng lao động (Contract)
*(Đã định nghĩa ở Chương 3.3.2)*

### 3.3.3 Thực thể Tiêu chuẩn công việc (JobStandard)
*(Đã định nghĩa ở Chương 3.3.3)*

### 3.3.4 Thực thể Danh mục năng lực (Competency)
*(Đã định nghĩa ở Chương 3.3.4)*

### 3.3.5 Thực thể Đánh giá năng lực (CompetencyEvaluation)
*(Đã định nghĩa ở Chương 3.3.5)*

### 3.3.6 Thực thể Log chấm công thô (AttendanceLog - Lưu trữ tại ClickHouse)
*(Đã định nghĩa ở Chương 3.3.6)*

### 3.3.7 Thực thể Bảng công ngày (DailyAttendanceSheet - Lưu trữ tại PostgreSQL)
*(Đã định nghĩa ở Chương 3.3.7)*

### 3.3.8 Thực thể Tờ trình giải trình bù công (ExplanationRequest)
*(Đã định nghĩa ở Chương 3.3.8)*

### 3.3.9 Thực thể Đăng ký làm thêm giờ (OTRequest)
*(Đã định nghĩa ở Chương 3.3.9)*

### 3.3.10 Thực thể Công thức lương (PayrollFormula)
*(Đã định nghĩa ở Chương 3.3.10)*

### 3.3.11 Thực thể Chu kỳ tính lương (PayrollPeriod)
*(Đã định nghĩa ở Chương 3.3.11)*

### 3.3.12 Thực thể Bảng ghi lương chi tiết (PayrollRecord)
*(Đã định nghĩa ở Chương 3.3.12)*

### 3.3.13 Thực thể Yêu cầu điều chỉnh lương (PayrollAdjustment)
*(Đã định nghĩa ở Chương 3.3.13)*

### 3.3.14 Thực thể Cảnh báo bất thường AI (AiAnomaly - Lưu tại ClickHouse)
* `id` (UUID, Primary Key)
* `record_id` (UUID, Foreign Key): Liên kết dòng lương bị lỗi.
* `metric_flagged` (VARCHAR): Chỉ mục bị lỗi (ví dụ: `OT_HOURS`, `NET_SALARY`).
* `deviation_value` (DECIMAL): Trị tuyệt đối số tiền/giờ bị lệch so với baseline.
* `ai_explanation` (TEXT): Lý do và phân tích giải nghĩa của AI.
* `status` (VARCHAR): Trạng thái xử lý (`PENDING`, `RESOLVED`, `IGNORED`).

### 3.3.15 Thực thể Tài liệu RAG Quy chế (RagDocument)
* `id` (UUID, Primary Key)
* `file_name` (VARCHAR): Tên tài liệu.
* `chunk_content` (TEXT): Nội dung đoạn văn bản được cắt nhỏ.
* `vector_embedding` (Array of Float): Vector đặc trưng của đoạn văn bản phục vụ so khớp ngữ nghĩa.

### 3.3.16 Thực thể Phiên hội thoại AI Phân tích (ChatSession)
* `id` (UUID, Primary Key)
* `user_id` (UUID, Foreign Key): Người thực hiện chat.
* `message` (TEXT): Câu hỏi của người dùng.
* `ai_response` (TEXT): Câu trả lời của AI.
* `created_at` (TIMESTAMP): Thời gian trao đổi.

---

## 3.4 Yêu cầu phi chức năng

### 3.4.1 Hiệu năng
* Tốc độ xử lý của API Gateway nhận check-in từ thiết bị camera di động phải đạt độ trễ < 10ms dưới mức chịu tải 10,000 request đồng thời nhờ cơ chế đệm NATS JetStream.
* Động cơ tính toán lương Batch (Go Goroutines) phải hoàn tất tính toán cho 10,000 nhân sự trong thời gian < 5 giây (chưa tính thời gian quét AI kiểm toán).

### 3.4.2 Độ tin cậy
* Cơ sở dữ liệu ClickHouse phải duy trì hoạt động liên tục đảm bảo không bị thất thoát log chấm công thô của nhân viên kể cả khi hệ thống PostgreSQL chính tạm dừng phục vụ bảo trì.

### 3.4.3 Tính sẵn sàng
* Hệ thống triển khai trên nền tảng Docker Container, hỗ trợ phục hồi và khởi động lại tự động khi phát sinh lỗi hệ điều hành. Mức độ sẵn sàng đạt 99.9% (Uptime).

### 3.4.4 Bảo mật
* Toàn bộ JWT token được mã hóa HMAC-SHA256, được lưu trữ vết thu hồi phiên trên Redis.
* Bảo mật RLS lớp API ngăn cấm hoàn toàn hành vi khai thác API thô từ phía client để xem trộm thông tin lương phòng ban khác.

### 3.4.5 Khả năng bảo trì
* Động cơ công thức chạy trên AST tự do cập nhật phép tính không cần build lại dự án.
* Codebase tuân thủ cấu trúc Clean Architecture của Golang và chia nhỏ Module nghiệp vụ độc lập.

### 3.4.6 Khả năng mở rộng
* Cho phép mở rộng quy mô xử lý tính lương Batch bằng cách gán thêm các worker node vật lý trong Go Worker Pool mà không làm xáo trộn kiến trúc Core Orchestrator.

---

## 3.5 Yêu cầu nghịch đảo
* Hệ thống KHÔNG được phép hiển thị bất kỳ số liệu nháp (`DRAFT`) hoặc lịch sử điều chỉnh lương chưa duyệt nào của các nhân viên thuộc phòng ban khác cho tài khoản cấp Trưởng phòng.

---

## 3.6 Ràng buộc thiết kế
* Cơ chế so khớp khuôn mặt phải hoạt động độc lập trên thiết bị di động (Client-side Edge AI) để giảm tải hoàn toàn tác vụ xử lý đồ họa/ảnh thô cho server Backend Go. Backend chỉ đảm nhận xác thực chữ ký token và so khớp dữ liệu nghiệp vụ.



