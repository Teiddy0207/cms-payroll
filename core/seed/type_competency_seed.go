package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SeedTypeCompetencies seeds initial type competencies data
func SeedTypeCompetencies(ctx context.Context, db database.Database) error {
	typeCompetencies := []struct {
		ID          uuid.UUID
		Code        string
		Name        string
		Description string
	}{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"), // ID cố định cho NT
			Code:        "NT",
			Name:        "Tiêu chuẩn năng lực nhận thức",
			Description: "Tiêu chuẩn đánh giá năng lực nhận thức của nhân viên",
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"), // ID cố định cho NL
			Code:        "NL",
			Name:        "Tiêu chuẩn năng lực hành nghề",
			Description: "Tiêu chuẩn đánh giá năng lực hành nghề của nhân viên",
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"), // ID cố định cho CM
			Code:        "CM",
			Name:        "Tiêu chuẩn năng lực chuyên môn và kỹ năng",
			Description: "Tiêu chuẩn đánh giá năng lực chuyên môn và kỹ năng của nhân viên",
		},
	}

	// Seed từng type competency
	for _, tc := range typeCompetencies {
		// Kiểm tra xem đã tồn tại chưa (theo ID hoặc code)
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM type_competencies WHERE id = $1 OR code = $2 LIMIT 1`, tc.ID, tc.Code)
		if err == nil && existingID != uuid.Nil {
			logger.Info("SeedTypeCompetencies: Type competency already exists", "code", tc.Code, "id", tc.ID)
			continue // Bỏ qua nếu đã tồn tại
		}

		// Sử dụng ID cố định từ struct
		query := `INSERT INTO type_competencies (id, code, name, description, is_active, is_default, created_at, updated_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

		err = db.ExecContext(ctx, query, tc.ID, tc.Code, tc.Name, tc.Description, true, true)
		if err != nil {
			logger.Error("SeedTypeCompetencies: Failed to create type competency", "code", tc.Code, "error", err)
			return fmt.Errorf("failed to seed type competency %s: %w", tc.Code, err)
		}

		logger.Info("SeedTypeCompetencies: Created type competency", "code", tc.Code, "name", tc.Name, "id", tc.ID)
	}

	logger.Info("SeedTypeCompetencies: Completed seeding type competencies")
	return nil
}
