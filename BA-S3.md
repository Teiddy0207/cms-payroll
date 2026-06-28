# BA-S3 · Module Ca kíp & Chấm công khuôn mặt (NATS Ingestion)

> **Phiên BA:** S3 (Giai đoạn Khảo sát & Đặc tả) · **Ngày:** 2026-06-28
> **Nuôi chương đặc tả:** Chương 3.2.3 (Module Ca kíp & Chấm công), Chương 3.3 (Đối tượng bổ sung)
> **Master plan:** [MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md](./MASTER-PLAN-PHAN-TICH-NGHIEP-VU.md) · **Walkthrough:** [ba_walkthrough.md](./ba_walkthrough.md)
> **Cần nạp context (dependency):** BA-S1 (để áp dụng phân quyền RLS cho việc duyệt đơn từ của Trưởng phòng)
> **Nguồn:** Thảo luận trực tiếp với Chủ đầu tư ngày 2026-06-28

---

## 1. Mục tiêu phiên & phạm vi
Phiên này làm rõ quy định ca làm việc tiêu chuẩn, quy tắc tính giờ làm thêm (OT), cơ chế thu thập dữ liệu check-in khuôn mặt từ Mobile App qua NATS JetStream, và thuật toán tính toán ngày công thực tế (Daily Attendance Calculation) kèm theo quy trình bù công bằng đơn giải trình.

* **Mục tiêu:**
  * Định nghĩa Ca làm việc tiêu chuẩn (8:30 - 17:30) kèm thời gian ân hạn đi muộn (30 phút).
  * Quy chuẩn cách tính OT (Nhân 1.5 cho đêm, nhân 2.0 cho ngày nghỉ) và ràng buộc duyệt OT từ Trưởng phòng.
  * Thiết lập luồng gửi nhận log chấm công từ thiết bị di động (nhận diện trên máy gửi Employee_ID) qua NATS JetStream.
  * Đặc tả thuật toán tính ngày công (1.0 công, 0.5 công) và quy trình duyệt tờ trình bù công cho người đi muộn/về sớm.
* **Phạm vi:** Cấu hình ca làm việc, luồng sự kiện NATS, tính công ngày và duyệt đơn từ bù công.
* **Ngoài phạm vi:** Động cơ cấu hình công thức lương và tính lương Batch hàng loạt (sẽ làm rõ ở phiên BA-S4).

---

## 2. Hiện trạng (as-is — kiểm chứng bằng mã nguồn)
* Dự án khởi tạo mới từ đầu, chưa có mã nguồn nghiệp vụ hiện trạng.
* Toàn bộ cấu trúc sẽ được xây dựng mới theo tài liệu đặc tả này.

---

## 3. Yêu cầu (to-be — dẫn tiêu chuẩn)
Các yêu cầu được thống nhất với Chủ đầu tư nhằm tự động hóa tối đa luồng chấm công:
* Hỗ trợ chấm công linh hoạt qua Mobile App nhận diện khuôn mặt tại văn phòng (Edge AI trên điện thoại tự trích xuất và chỉ gửi mã nhân sự `employee_code`).
* Hứng tải và điều tiết dữ liệu check-in qua **NATS JetStream** để tránh nghẽn DB.
* Quy tắc công ngày nghiêm ngặt nhưng cho phép giảm trừ linh động qua tờ trình giải trình lý do được Trưởng phòng phê duyệt.

---

## 4. Quy tắc nghiệp vụ & tham số (sinh REQ)

| Mã | Phát biểu | Nguồn | Loại | As-is | Phiên | Nghiệm thu |
|:---|:---|:---|:---|:---|:---|:---|
| **REQ-SCHED-001** | **Ca hành chính chuẩn:** Giờ vào: 08:30, Giờ ra: 17:30. Thời gian đi muộn cho phép (Grace Period) là 30 phút (từ 08:30 đến 09:00 vẫn tính là đúng giờ). | Khách hàng | Business Rule | Mới | S3 | [ ] |
| **REQ-OT-001**    | **Hệ số làm thêm giờ (OT):** Làm thêm ban đêm nhân hệ số 1.5. Làm thêm vào ngày nghỉ/lễ nhân hệ số 2.0. | Khách hàng | Business Rule | Mới | S3 | [ ] |
| **REQ-OT-002**    | **Kiểm soát OT:** Mọi giờ làm thêm sau 17:30 chỉ được tính là OT chịu lương nếu có Đơn đăng ký OT được Trưởng phòng phê duyệt. Trường hợp ở lại làm muộn không có đơn duyệt sẽ không tính OT. | Khách hàng | Business Rule | Mới | S3 | [ ] |
| **REQ-ING-001**   | **Chấm công qua Mobile App:** Thiết bị di động của nhân viên chạy mô hình nhận diện khuôn mặt cục bộ, sau đó gửi sự kiện `{employee_id, timestamp, device_info, location_gps}` về API Gateway của công ty. | Khách hàng | Functional | Mới | S3 | [ ] |
| **REQ-ING-002**   | **Hứng tải qua NATS:** API Gateway ném thẳng sự kiện check-in vào queue của NATS JetStream chủ đề `timekeeping.checkin` để phản hồi tức thì về thiết bị trong <5ms. Worker Go sẽ consume bất đồng bộ để ghi vào ClickHouse và Postgres. | Khách hàng | Technical | Mới | S3 | [ ] |
| **REQ-CALC-001**  | **Quy tắc 1.0 công:** Nhân viên có ít nhất 02 lần quẹt thẻ trong ngày: Lượt quẹt đầu tiên từ 08:30 đến 09:00 (hoặc sớm hơn) VÀ lượt quẹt cuối cùng từ 17:30 trở đi. | Khách hàng | Business Rule | Mới | S3 | [ ] |
| **REQ-CALC-002**  | **Quy tắc 0.5 công:** Nhân viên làm việc nửa ngày, thực hiện check-in buổi sáng trước 09:00 và thực hiện check-out (lượt quẹt cuối) trong khoảng thời gian từ 12:00 đến 14:00. | Khách hàng | Business Rule | Mới | S3 | [ ] |
| **REQ-EXC-001**   | **Tờ trình giải trình (Bù công):** Nhân viên đi muộn (quẹt thẻ sau 09:00) hoặc về sớm trước giờ (về trước 17:30 từ 1-2 tiếng) phải tạo "Tờ trình giải trình". Nếu được Trưởng phòng phê duyệt, ngày làm việc đó vẫn tính tròn 1.0 công bình thường. | Khách hàng | Business Rule | Mới | S3 | [ ] |

---

## 5. Câu hỏi & Quyết định từ Chủ đầu tư (Q&A)

### Q-S3-01: Cách thức xác định ca trực và hệ số làm thêm giờ (OT)?
* **Quyết định:** Giờ hành chính là 8:30 - 17:30. Cho phép đi muộn 30 phút (đến 9h vẫn tính đúng giờ). OT đêm nhân hệ số 1.5, ngày nghỉ nhân hệ số 2.0. Việc tính OT bắt buộc phải được Trưởng phòng duyệt, nếu không có đơn duyệt thì không tính tiền OT dù ở lại muộn.

### Q-S3-02: Thiết bị và cơ chế thu thập chấm công?
* **Quyết định:** Cho phép chấm công bằng Mobile App định vị GPS tại công ty. Thiết bị di động nhận dạng khuôn mặt xong chỉ gửi `Employee_ID` và thời gian về Backend. Không phân định camera IN/OUT, chỉ cần 2 lần quẹt một ngày là được.

### Q-S3-03: Quy tắc tổng hợp công ngày và xử lý đi muộn/về sớm?
* **Quyết định:** Quẹt vào trước 9h và quẹt ra sau 17h30 tính 1 công. Quẹt ra từ 12h đến 14h tính 0.5 công. Các trường hợp đi muộn, về sớm trước 1-2 tiếng phải viết tờ trình giải trình gửi Trưởng phòng duyệt để được giữ nguyên 1 công.

---

## 6. Tiêu chí chấp nhận sơ bộ của phiên (Acceptance Criteria)

- [ ] Thiết kế bảng cơ sở dữ liệu `attendance_logs` tại ClickHouse (lưu log thô) và bảng `daily_attendance_sheets` tại PostgreSQL (lưu kết quả tính công ngày).
- [ ] Xây dựng NATS JetStream Consumer nhận diện sự kiện, lọc trùng tin nhắn dựa trên mã nhân sự và phút quẹt thẻ.
- [ ] Viết hàm tính toán công ngày tự động chạy định kỳ mỗi đêm lúc 23:59 dựa trên các quy tắc 1.0 công, 0.5 công và bảng đơn từ giải trình đã duyệt.

---

## 7. Đóng góp vào chương đặc tả

| Mục phiên này | Chương đích |
|---|---|
| Quy tắc REQ-SCHED-001, REQ-OT-001, REQ-OT-002 | Chương 3.2.3.1 (Thiết lập ca và Chính sách OT) |
| Quy tắc REQ-ING-001, REQ-ING-002 | Chương 3.2.3.2 (Luồng thu thập chấm công qua NATS JetStream) |
| Quy tắc REQ-CALC-001, REQ-CALC-002, REQ-EXC-001 | Chương 3.2.3.3 (Động cơ tính công ngày tự động & Quy trình bù công) |
| Cấu trúc dữ liệu Đối tượng | Chương 3.3 (Cơ cấu đối tượng bổ sung - Log Chấm công & Đơn từ) |
