package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID   uuid.UUID `db:"id"`
	Code string    `db:"code"`
}

func SeedAttendanceData(ctx context.Context, db database.Database) error {
	// 1. Lấy danh sách toàn bộ nhân sự
	var profiles []Profile
	err := db.SelectContext(ctx, &profiles, `SELECT id, code FROM user_profiles`)
	if err != nil {
		return fmt.Errorf("failed to get user profiles for seeding: %w", err)
	}

	if len(profiles) == 0 {
		logger.Warn("SeedAttendanceData: No user profiles found to seed attendance for")
		return nil
	}

	// 2. Clear old mock data to start clean
	_ = db.ExecContext(ctx, `DELETE FROM daily_attendance_sheets`)
	_ = db.ExecContext(ctx, `DELETE FROM attendance_logs`)

	logger.Info(fmt.Sprintf("SeedAttendanceData: Seeding attendance for %d employees for July 2026...", len(profiles)))

	// 3. Chuẩn bị khoảng thời gian: Tháng 7/2026 (từ ngày 1 đến ngày 31)
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	if loc == nil {
		loc = time.Local
	}

	rand.Seed(time.Now().UnixNano())

	// Lặp qua từng ngày trong tháng 7 năm 2026
	for day := 1; day <= 31; day++ {
		currentDay := time.Date(2026, 7, day, 0, 0, 0, 0, loc)
		// Bỏ qua Chủ Nhật
		if currentDay.Weekday() == time.Sunday {
			continue
		}

		for _, p := range profiles {
			// Tạo ngẫu nhiên một chút:
			// 8% cơ hội vắng mặt (ABSENT)
			// 10% cơ hội đi muộn (LATE)
			// 5% cơ hội về sớm (EARLY)
			// 77% cơ hội đúng giờ (NORMAL)
			roll := rand.Intn(100)

			var status string
			var checkIn, checkOut *time.Time
			var workDay float64 = 0.0
			var otHours float64 = 0.0

			if roll < 8 {
				// ABSENT
				status = "ABSENT"
				workDay = 0.0
			} else {
				workDay = 1.0
				// Thiết lập checkin / checkout ngẫu nhiên quanh mốc 8:00 và 17:30
				inHour := 7
				inMin := 45 + rand.Intn(15) // 7:45 - 8:00
				outHour := 17
				outMin := 30 + rand.Intn(30) // 17:30 - 18:00

				if roll >= 8 && roll < 18 {
					// LATE (Checkin sau 8:15)
					status = "LATE"
					inHour = 8
					inMin = 16 + rand.Intn(30) // 8:16 - 8:46
				} else if roll >= 18 && roll < 23 {
					// EARLY (Checkout trước 17:00)
					status = "EARLY"
					outHour = 16
					outMin = rand.Intn(59) // 16:00 - 16:59
				} else {
					status = "NORMAL"
				}

				// Thỉnh thoảng làm OT (khoảng 15% số ngày bình thường)
				if status == "NORMAL" && rand.Intn(100) < 15 {
					outHour = 19 + rand.Intn(3)  // 19:00 - 22:00
					otHours = float64(outHour-17) + float64(outMin)/60.0
				}

				tIn := time.Date(2026, 7, day, inHour, inMin, rand.Intn(60), 0, loc)
				tOut := time.Date(2026, 7, day, outHour, outMin, rand.Intn(60), 0, loc)
				checkIn = &tIn
				checkOut = &tOut
			}

			sheetID := uuid.New()
			// Insert daily_attendance_sheets
			insertSheet := `INSERT INTO daily_attendance_sheets (id, employee_id, date, check_in, check_out, actual_work_day, ot_hours, status, created_at, updated_at)
			                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			                ON CONFLICT (employee_id, date) DO NOTHING`
			err = db.ExecContext(ctx, insertSheet, sheetID, p.ID, currentDay.Format("2006-01-02"), checkIn, checkOut, workDay, otHours, status)
			if err != nil {
				logger.Error("SeedAttendanceData: Failed to insert sheet", "employee", p.Code, "date", currentDay.Format("2006-01-02"), "error", err)
			}

			// Đồng thời insert log thô (attendance_logs) nếu có đi làm
			if checkIn != nil {
				logInID := uuid.New()
				insertLogIn := `INSERT INTO attendance_logs (id, employee_code, timestamp, location_gps, device_id, event_id, created_at)
				                VALUES ($1, $2, $3, $4, $5, $1, NOW())`
				_ = db.ExecContext(ctx, insertLogIn, logInID, p.Code, *checkIn, "10.01, 106.32", "FACE_KIOSK_01")
			}
			if checkOut != nil {
				logOutID := uuid.New()
				insertLogOut := `INSERT INTO attendance_logs (id, employee_code, timestamp, location_gps, device_id, event_id, created_at)
				                 VALUES ($1, $2, $3, $4, $5, $1, NOW())`
				_ = db.ExecContext(ctx, insertLogOut, logOutID, p.Code, *checkOut, "10.01, 106.32", "FACE_KIOSK_01")
			}
		}
	}

	logger.Info("SeedAttendanceData: Seeding attendance records completed successfully!")
	return nil
}
