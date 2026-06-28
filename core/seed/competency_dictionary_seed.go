package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SeedCompetencyDictionaries seeds initial competency dictionaries data
func SeedCompetencyDictionaries(ctx context.Context, db database.Database) error {

	competencyDictionaries := []struct {
		Code                string
		Name                string
		TypeID              uuid.UUID
		Description         string
		ProficiencyLevelMax float64
	}{
		{
			Code:   "NT01",
			Name:   "Trình độ học vấn",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
		{
			Code:   "NT02",
			Name:   "Vi tính, thiết bị kỹ thuật",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
		{
			Code:   "NT03",
			Name:   "Ngoại ngữ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NT04",
			Name:   "Kinh nghiệm chuyên môn",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL01",
			Name:   "Lập kế hoạch",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL02",
			Name:   "Hiểu biết",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL03",
			Name:   "Phán quyết - ra quyết định",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL04",
			Name:   "Khả năng thuyết phục",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL05",
			Name:   "Tính sáng tạo",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL06",
			Name:   "Giải quyết vấn đề",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL07",
			Name:   "Huấn luyện đào tạo",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "NL08",
			Name:   "Năng lực lãnh đạo",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH",
			Name:   "Tiếp thị bán hàng",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH01",
			Name:   "Năng lực tiếp thị bán hàng",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH02",
			Name:   "Kiến thức về khách hàng",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH03",
			Name:   "kiến thức về thị trường",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH04",
			Name:   "kiến thức về đối thủ cạnh tranh",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BH05",
			Name:   "Năng lực tổ chức sự kiện hội nghị khách hàng",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BV",
			Name:   "Năng lực bảo vệ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BV01",
			Name:   "Năng lực phục vụ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "BV02",
			Name:   "Năng lực bảo vệ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "KT",
			Name:   "Năng lực Tài chính - Kế toán",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "KT01",
			Name:   "Năng lực Tài chính",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "KT02",
			Name:   "Năng lực kế toán",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "KT03",
			Name:   "Năng lực tính toán và nhạy bén với con số",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // NT - Tiêu chuẩn năng lực nhận thức
		},
		{
			Code:   "KT04",
			Name:   "Năng lực kiểm toán",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM - Tiêu chuẩn năng lực chuyên môn và kỹ năng
		},
		// KS - Năng lực Kiểm soát nội bộ
		{
			Code:   "KS01",
			Name:   "Năng lực phân tích tài chính",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "KS02",
			Name:   "Năng lực thông tin kế toán",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "KS03",
			Name:   "Năng lực hệ thống quản lý",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// NS - Năng lực Hành chính - Nhân sự
		{
			Code:   "NS01",
			Name:   "Năng lực đánh giá và chọn lựa nhân sự",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS02",
			Name:   "Năng lực quản lý hành chánh và luật pháp",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS03",
			Name:   "Năng lực quản lý nhân sự",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS04",
			Name:   "Quản lý ngân sách nhân sự",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS05",
			Name:   "Năng lực phát triển tổ chức",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS06",
			Name:   "Kiến thức về lưu trữ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS07",
			Name:   "Năng lực duy trì, và phát triển hệ thống quản lý chất lượng (QMS), ISO",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "NS08",
			Name:   "Năng lực về tổng hợp",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// TL - Năng lực Bảo hộ lao động
		{
			Code:   "TL01",
			Name:   "Năng lực về tổ chức cán bộ",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL02",
			Name:   "Năng lực hiểu biết về lao động – tiền lương",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL03",
			Name:   "Năng lực về Bảo hộ lao động",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL04",
			Name:   "Năng lực về Bảo vệ chăm sóc sức khỏe",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL06",
			Name:   "Năng lực pháp chế",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL07",
			Name:   "Năng lực quản lý đào tạo",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TL08",
			Name:   "Năng lực quản lý huấn luyện, đào tạo lại Thuyền viên",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "BV03",
			Name:   "Năng lực về An toàn, PCCC, PCCN",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// PL - Năng lực pháp luật
		{
			Code:   "PL01",
			Name:   "Năng lực luật thương mại",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "PL02",
			Name:   "Năng lực luật dân sự",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "PL03",
			Name:   "Năng lực luật kinh doanh",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// SC - Năng lực sửa chữa đầu tư thiết bị
		{
			Code:   "SC01",
			Name:   "Năng lực sửa chữa máy móc, thiết bị",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SC02",
			Name:   "Năng lực chọn nhà cung cấp sửa chữa",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// DM - Năng lực chuyên ngành vận tải biển
		{
			Code:   "DM02",
			Name:   "Năng lực của các công nghệ / thiết bị tàu biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DM04",
			Name:   "Năng lực quản lý an toàn lao động ngành vận tải biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DM05",
			Name:   "Năng lực quản lý môi trường ngành vận tải biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DM08",
			Name:   "Năng lực thực hành về công nghệ vận hành tàu biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DM09",
			Name:   "Năng lực thực hành về máy móc, thiết bị tàu biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// XD - Năng lực đầu tư xây dựng
		{
			Code:   "XD01",
			Name:   "Năng lực thiết kế quy hoạch",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "XD02",
			Name:   "Năng lực vẽ thiết kế và xây dựng mô hình",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "XD03",
			Name:   "Năng lực kiến trúc",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "XD04",
			Name:   "Năng lực chọn nhà thầu",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// DT - Năng lực lập dự án đầu tư
		{
			Code:   "DT01",
			Name:   "Năng lực về đặc điểm kinh doanh của ngành Vận tải và Thuê tàu biển",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DT02",
			Name:   "Kiến thức về luật chuyên ngành đầu tư",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DT03",
			Name:   "Khả năng thẩm định tính kinh tế của dự án",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DT04",
			Name:   "Tổng hợp và phân tích thông tin lĩnh vực đầu tư",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "DT05",
			Name:   "Kỹ thuật viết dự án đầu tư",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// TC - Năng lực thi công
		{
			Code:   "TC01",
			Name:   "Quản lý thời gian thực hiện dự án",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TC02",
			Name:   "Quản lý rủi ro trong thực hiện dự án",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TC03",
			Name:   "Quản lý chi phí",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TC04",
			Name:   "Quản lý nguồn lực",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TC05",
			Name:   "Khả năng Giám sát thi công",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TC06",
			Name:   "Kiến thức về kỹ thuật thi công",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// LX - Năng lực Lái xe
		{
			Code:   "LX01",
			Name:   "Năng lực lái xe văn phòng",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// LK - Năng lực khác
		{
			Code:   "DC01",
			Name:   "Năng lực Đầu tư Tài chính",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "XN01",
			Name:   "Năng lực làm chứng từ xuất nhập khẩu",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "VK01",
			Name:   "Năng lực Kinh tế - Kế hoạch",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "MH01",
			Name:   "Năng lực đánh giá và lựa chọn nhà cung cấp hàng hóa, vật tư, nguyên liệu, nhiên liệu...",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "HD01",
			Name:   "Kiểm soát hợp đồng cung cấp",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TH04",
			Name:   "Năng lực phần mềm chuyên môn",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "TV01",
			Name:   "Năng lực kế hoạch - điều hành cung ứng Thuyền viên",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SK01",
			Name:   "Kiểm soát kho",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "RR01",
			Name:   "Năng lực Quản trị Rủi ro",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// MK - Năng lực chuyên môn Marketing
		{
			Code:   "MK01",
			Name:   "Năng lực xây dựng, quản lý & phát triển thương hiệu",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "MK03",
			Name:   "Năng lực quan hệ công chúng (PR)",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// SS - Kỹ năng mềm
		{
			Code:   "SS01",
			Name:   "Giao tiếp",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS03",
			Name:   "Kiểm soát công việc cá nhân",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS04",
			Name:   "Sử dụng vi tính, thiết bị kỹ thuật",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS05",
			Name:   "Soạn thảo văn bản",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS06",
			Name:   "Hợp tác đồng đội",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS07",
			Name:   "Cam kết lòng trung thành nghề nghiệp",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS08",
			Name:   "Trung thực, đáng tin cậy",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "SS09",
			Name:   "Khả năng làm việc độc lập: chủ động trong công việc",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		// QL - Quản lý
		{
			Code:   "QL02",
			Name:   "Tổ chức và điều phối",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "QL04",
			Name:   "Ủy thác công việc",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "QL06",
			Name:   "Kiểm tra giám sát",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "QL07",
			Name:   "Đánh giá nhân viên",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
		{
			Code:   "QL09",
			Name:   "Quản lý nguồn lực: huy động và sử dụng hiệu quả tài nguyên, vật lực, nhân lực, tài chính, công nghệ, thông tin, thời gian",
			TypeID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), // CM
		},
	}

	// Seed từng competency dictionary
	for _, cd := range competencyDictionaries {
		// Kiểm tra xem đã tồn tại chưa (theo code)
		var existingID uuid.UUID
		err := db.GetContext(ctx, &existingID, `SELECT id FROM competency_dictionaries WHERE code = $1 LIMIT 1`, cd.Code)
		if err == nil && existingID != uuid.Nil {
			continue // Bỏ qua nếu đã tồn tại
		}

		// Kiểm tra type_competency có tồn tại không
		var typeExists bool
		err = db.GetContext(ctx, &typeExists, `SELECT EXISTS(SELECT 1 FROM type_competencies WHERE id = $1)`, cd.TypeID)
		if err != nil || !typeExists {
			continue // Bỏ qua nếu không tìm thấy type
		}

		// Tạo competency dictionary mới với is_default = true
		competencyID := uuid.New()
		proficiencyLevelMax := cd.ProficiencyLevelMax
		if proficiencyLevelMax == 0 {
			proficiencyLevelMax = 5.0 // Giá trị mặc định
		}
		query := `INSERT INTO competency_dictionaries (id, code, name, type_id, description, is_default, proficiency_level_max, created_at, updated_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`

		err = db.ExecContext(ctx, query, competencyID, cd.Code, cd.Name, cd.TypeID, cd.Description, true, proficiencyLevelMax)
		if err != nil {
			logger.Error("SeedCompetencyDictionaries: Failed to create competency dictionary", "code", cd.Code, "error", err)
			return fmt.Errorf("failed to seed competency dictionary %s: %w", cd.Code, err)
		}

		logger.Info("SeedCompetencyDictionaries: Created competency dictionary", "code", cd.Code, "name", cd.Name, "type_id", cd.TypeID)
	}

	logger.Info("SeedCompetencyDictionaries: Completed seeding competency dictionaries")
	return nil
}
