package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// SeedUserProfiles seeds initial user profiles data
func SeedUserProfiles(ctx context.Context, db database.Database) error {
	userProfiles := []struct {
		Code     string
		FullName string
	}{
		{Code: "2510", FullName: "CAO XUÂN HOAN"},
		{Code: "2529", FullName: "ĐẶNG XUÂN BÌNH"},
		{Code: "2506", FullName: "ĐÀO HỒNG QUANG"},
		{Code: "2897", FullName: "TRẦN ĐẠI HỒNG DŨNG"},
		{Code: "2825", FullName: "NGUYỄN THỤY LƯU"},
		{Code: "3442", FullName: "HOÀNG THỊ ĐÔNG"},
		{Code: "3453", FullName: "PHAN TÁ NGHĨA"},
		{Code: "2519", FullName: "NGUYỄN NGỌC LINH"},
		{Code: "2592", FullName: "TRẦN NGỌC HUY"},
		{Code: "2892", FullName: "NGUYỄN MINH HÀ"},
		{Code: "1908", FullName: "LÊ THỊ VIỆT HOA"},
		{Code: "2622", FullName: "NGUYỄN TRẦN CHUNG THƯ"},
		{Code: "3656", FullName: "DƯƠNG ĐÌNH PHI"},
		{Code: "3908", FullName: "DIỆP THỊ HỒNG LOAN"},
		{Code: "3968", FullName: "ĐẶNG HOÀNG MINH"},
		{Code: "4005", FullName: "HOÀNG THỊ THU HÀ"},
		{Code: "4232", FullName: "TRẦN THỊ TUYẾT MAI"},
		{Code: "4456", FullName: "DƯƠNG VĂN TƯỞNG"},
		{Code: "4461", FullName: "TRẦN BẢO TRINH"},
		{Code: "5237", FullName: "NGUYỄN HUY HOÀNG"},
		{Code: "5266", FullName: "TRẦN NGUYỄN PHƯƠNG UY"},
		{Code: "5269", FullName: "NGUYỄN THỊ THÚY PHƯỢNG"},
		{Code: "5256", FullName: "PHẠM HOÀNG TRÚC LINH"},
		{Code: "5334", FullName: "NGUYỄN VĂN HOÀI"},
		{Code: "5235", FullName: "VŨ QUỐC CÔNG"},
		{Code: "5342", FullName: "HÀ HỒNG ĐẠT"},
	}

	// Đếm số mục đã seed thành công và số mục bị lỗi
	successCount := 0
	errorCount := 0
	skipCount := 0

	// Seed từng user profile
	for _, up := range userProfiles {
		// Kiểm tra xem đã tồn tại chưa (theo code)
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM user_profiles WHERE code = $1 LIMIT 1`, up.Code)
		if err == nil && existingID != uuid.Nil {
			skipCount++
			continue // Bỏ qua nếu đã tồn tại
		} else if err != nil && err != sql.ErrNoRows {
			logger.Error("SeedUserProfiles: Error checking existing profile", "code", up.Code, "error", err)
			errorCount++
			continue
		}

		// Tạo UUID mới cho profile
		profileID := uuid.New()

		// Insert user profile với chỉ code và full_name, các trường khác để NULL
		query := `INSERT INTO user_profiles (id, code, full_name, created_at, updated_at) 
		         VALUES ($1, $2, $3, NOW(), NOW())`

		err = db.ExecContext(ctx, query, profileID, up.Code, up.FullName)
		if err != nil {
			logger.Error("SeedUserProfiles: Failed to create user profile",
				"code", up.Code,
				"full_name", up.FullName,
				"id", profileID,
				"error", err)
			errorCount++
			continue
		}

		successCount++

	}

	if errorCount > 0 {
		return fmt.Errorf("seed completed with %d errors out of %d total items", errorCount, len(userProfiles))
	}

	return nil
}
