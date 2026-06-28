# Thiết kế Calculator Lương P1 & P2

> **Phiên:** BA-S2 (tiếp theo)
> **Ngày:** 2026-06-28
> **Phụ thuộc:** Đã hoàn thành CRUD `contracts`, `job_standards`, `competency_evaluations`

---

## 1. Phân tích — 2 mắt xích đang thiếu

Công thức đã được đặc tả rõ trong SRS, nhưng nhìn vào database hiện tại có **2 liên kết bị đứt** chưa được triển khai:

### Mắt xích 1 — P1 thiếu bảng JOIN chức danh ↔ tiêu chuẩn

```
job_descriptions (Chức danh)  ←— ??? —→  job_standards (Tiêu chuẩn)
```

`job_standards` hiện là bảng độc lập. Không có bảng nối để biết chức danh X yêu cầu tiêu chuẩn nào.
→ **Cần thêm:** `job_position_standards`

### Mắt xích 2 — P2 thiếu cấu hình `Company_Point_Rate`

```
P2 = SUM(Score × Weight) × Company_Point_Rate  ← Con số này lưu ở đâu?
```

→ **Cần thêm:** Bảng `system_settings` dạng key-value

---

## 2. Chiến lược lưu trữ Tiêu chuẩn — Cách 3 (Kết hợp)

### So sánh 3 phương án

| Tiêu chí | Cách 1: Theo Chức danh | Cách 2: Theo Nhân viên | **Cách 3: Kết hợp ⭐** |
|---|---|---|---|
| Độ chính xác | Thấp | Cao | **Cao nhất** |
| Công sức HR | Ít | Nhiều | **Trung bình** |
| Phù hợp thực tế | Không sát | Sát | **Rất sát** |

### Cách 3 hoạt động như thế nào?

```
P1 = Position_Base_Rate
   + SUM(Tiêu chuẩn BẮT BUỘC của Chức danh)   ← cấu hình 1 lần cho cả vị trí
   + SUM(Tiêu chuẩn CÁ NHÂN của Nhân viên)     ← nhập theo từng người khi đạt được
```

**Ví dụ thực tế:**

```
Chức danh "Kế toán trưởng" — tiêu chuẩn bắt buộc:
└── DEGREE_BACHELOR  →  +500,000đ   (ai giữ chức danh này đều được)

Nhân viên A (Kế toán trưởng) — chứng chỉ cá nhân thực tế đạt được:
├── CERT_CPA         →  +1,500,000đ
└── CERT_ENGLISH_B2  →  +300,000đ

→ P1_A = 12,000,000 + 500,000 + 1,500,000 + 300,000 = 14,300,000đ

Nhân viên B (cùng chức danh, chưa có chứng chỉ gì thêm):
→ P1_B = 12,000,000 + 500,000 = 12,500,000đ
```

---

## 3. Schema Database cần bổ sung

### Bảng mới 1: `job_position_standards`

```sql
CREATE TABLE IF NOT EXISTS job_position_standards (
    job_description_id UUID NOT NULL REFERENCES job_descriptions(id) ON DELETE CASCADE,
    job_standard_id    UUID NOT NULL REFERENCES job_standards(id) ON DELETE CASCADE,
    created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (job_description_id, job_standard_id)
);
```

### Bảng mới 2: `user_profile_standards`

```sql
CREATE TABLE IF NOT EXISTS user_profile_standards (
    user_profile_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    job_standard_id UUID NOT NULL REFERENCES job_standards(id) ON DELETE CASCADE,
    achieved_at     DATE,
    note            TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_profile_id, job_standard_id)
);
```

### Bảng mới 3: `system_settings`

```sql
CREATE TABLE IF NOT EXISTS system_settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT NOT NULL,
    description TEXT,
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO system_settings (key, value, description)
VALUES ('company_point_rate', '5000', 'Đơn giá quy đổi 1 điểm năng lực thành VND')
ON CONFLICT (key) DO NOTHING;
```

### Bảng mới 4: `payroll_record_details`

```sql
CREATE TABLE IF NOT EXISTS payroll_record_details (
    id          UUID PRIMARY KEY,
    record_id   UUID NOT NULL REFERENCES payroll_records(id) ON DELETE CASCADE,
    component   VARCHAR(50)    NOT NULL, -- BASE_RATE | P1_STANDARD | P2_COMP | P2_TOTAL
    description TEXT           NOT NULL,
    source      VARCHAR(20)    NOT NULL, -- CONTRACT | POSITION | PERSONAL | EVAL | FORMULA
    amount      DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 4. Luồng tính toán P1

### Công thức
```
P1_BASE = Position_Base_Rate + SUM(Standard_Value)
```

### SQL Query (UNION tiêu chuẩn bắt buộc + cá nhân)

```sql
SELECT position_base_rate FROM contracts
WHERE employee_id = $1 AND status = 'ACTIVE' LIMIT 1;

SELECT js.standard_code, js.name, js.allowance_value, 'POSITION' AS source
FROM job_position_standards jps
JOIN job_standards js ON js.id = jps.job_standard_id
JOIN user_profiles up ON up.position_id = jps.job_description_id
WHERE up.id = $1

UNION

SELECT js.standard_code, js.name, js.allowance_value, 'PERSONAL' AS source
FROM user_profile_standards ups
JOIN job_standards js ON js.id = ups.job_standard_id
WHERE ups.user_profile_id = $1;
```

### Kết quả trả về

```json
{
  "employee_id": "uuid-A",
  "p1": {
    "position_base_rate": 12000000,
    "standards": [
      { "code": "DEGREE_BACHELOR", "source": "POSITION", "value": 500000 },
      { "code": "CERT_CPA",        "source": "PERSONAL", "value": 1500000 },
      { "code": "CERT_ENGLISH_B2", "source": "PERSONAL", "value": 300000 }
    ],
    "total_standards": 2300000,
    "p1_base": 14300000
  }
}
```

---

## 5. Luồng tính toán P2

### Công thức
```
P2_COMPETENCY = SUM(Competency_Score × Competency_Weight) × Company_Point_Rate
```

### Ví dụ chi tiết

```
Kỳ đánh giá 2026-H1, Nhân viên A:
┌────────────────────┬───────┬────────┬──────────────┐
│ Năng lực           │ Score │ Weight │ Score×Weight │
├────────────────────┼───────┼────────┼──────────────┤
│ Tư duy logic       │  85   │   3    │     255      │
│ Kỹ năng viết code  │  90   │   5    │     450      │
│ Giao tiếp KH       │  70   │   2    │     140      │
├────────────────────┼───────┼────────┼──────────────┤
│ TỔNG               │       │        │     845      │
└────────────────────┴───────┴────────┴──────────────┘

Company_Point_Rate = 5,000đ/điểm
P2 = 845 × 5,000 = 4,225,000đ
```

---

## 6. Chiến lược lưu trữ — Snapshot Pattern

> **Nguyên tắc:** Lưu **số tiền đã tính ra** (baked values), không lưu tham chiếu công thức.
> Lương tháng cũ **bất biến hoàn toàn** dù cấu hình sau này thay đổi.

### Sơ đồ quan hệ

```
payroll_periods (1 tháng)
    │
    └── payroll_records (1 dòng / nhân viên)
            │   ├── p1_value  = 14,300,000  ← snapshot
            │   ├── p2_value  =  4,225,000  ← snapshot
            │   └── status: DRAFT → ADJUSTING → APPROVED
            │
            └── payroll_record_details (breakdown chi tiết)
                    ├── BASE_RATE   = 12,000,000  (CONTRACT)
                    ├── P1_STANDARD =    500,000  (POSITION)
                    ├── P1_STANDARD =  1,500,000  (PERSONAL)
                    ├── P2_COMP     =        255  (EVAL)
                    └── P2_TOTAL    =  4,225,000  (FORMULA)
```

### Vòng đời trạng thái

```
DRAFT ──► DRAFT ──► ADJUSTING ──► APPROVED ──► LOCKED (vĩnh viễn bất biến)
```

---

## 7. API cần triển khai

### Nhóm 1 — Tiêu chuẩn Chức danh

| Method | Endpoint | Mô tả |
|---|---|---|
| `POST` | `/job-positions/:id/standards` | Gắn tiêu chuẩn vào chức danh |
| `GET` | `/job-positions/:id/standards` | Xem tiêu chuẩn chức danh |
| `DELETE` | `/job-positions/:id/standards/:std_id` | Gỡ tiêu chuẩn |

### Nhóm 2 — Chứng chỉ Cá nhân

| Method | Endpoint | Mô tả |
|---|---|---|
| `POST` | `/user-profiles/:id/standards` | Thêm chứng chỉ cho nhân viên |
| `GET` | `/user-profiles/:id/standards` | Xem chứng chỉ nhân viên |
| `DELETE` | `/user-profiles/:id/standards/:std_id` | Xóa chứng chỉ |

### Nhóm 3 — Cấu hình Hệ thống

| Method | Endpoint | Mô tả |
|---|---|---|
| `GET` | `/settings` | Lấy toàn bộ cấu hình |
| `PUT` | `/settings/:key` | Cập nhật giá trị (vd: `company_point_rate`) |

### Nhóm 4 — Calculator Engine

| Method | Endpoint | Mô tả |
|---|---|---|
| `GET` | `/calculator/p1/:employee_id` | Tính P1 với breakdown |
| `GET` | `/calculator/p2/:employee_id?period=2026-H1` | Tính P2 với breakdown |
| `GET` | `/calculator/preview/:employee_id?period=2026-H1` | Xem trước P1+P2 |

---

## 8. Kế hoạch triển khai

- [ ] **Bước 1:** Bổ sung 4 bảng mới vào `db/init_schema.sql`
- [ ] **Bước 2:** Tạo API Nhóm 1 — Tiêu chuẩn Chức danh
- [ ] **Bước 3:** Tạo API Nhóm 2 — Chứng chỉ Cá nhân
- [ ] **Bước 4:** Tạo API Nhóm 3 — Cấu hình Hệ thống
- [ ] **Bước 5:** Viết `P1Calculator` Service
- [ ] **Bước 6:** Viết `P2Calculator` Service
- [ ] **Bước 7:** Expose endpoint `/calculator/preview`
- [ ] **Bước 8:** Test toàn bộ với dữ liệu mẫu
