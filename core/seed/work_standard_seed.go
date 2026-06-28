package seed

import (
	"cal-salary/core/database"
	"cal-salary/core/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Helper function to convert UUID to pointer
func uuidPtr(id string) *uuid.UUID {
	parsed := uuid.MustParse(id)
	return &parsed
}

// Helper function to convert code to order
// Maps code patterns to order values
func codeToOrder(code string) string {
	// Map các code đặc biệt sang order
	codeOrderMap := map[string]string{
		"YT19": "3",
		"YT13": "4.1",
		"YT14": "4.2",
		"YT15": "4.3",
		"YT16": "4.4",
		"YT17": "4.5",
		"YT22": "5.2",
		"YT20": "5.3",
		"YT18": "5.4",
		"YT21": "5.5",
		"YT1":  "6.2",
		"YT3":  "6.3",
		"YT2":  "6.4",
		"YT4":  "6.5",
		"YT5":  "7.1",
		"YT6":  "7.2",
		"YT7":  "7.3",
		"YT8":  "7.4",
		"YT9":  "7.5",
		"YT10": "7.6",
		"YT11": "7.7",
		"YT12": "7.8",
		"YT23": "9",
	}

	// Nếu code có trong map, trả về order tương ứng
	if order, exists := codeOrderMap[code]; exists {
		return order
	}

	// Nếu code có dạng YTXX_N (ví dụ: YT19_1, YT13_2)
	if len(code) > 5 && code[:2] == "YT" && code[4] == '_' {
		baseCode := code[:4] // YT19, YT13, ...
		if baseOrder, exists := codeOrderMap[baseCode]; exists {
			// Lấy số sau dấu _
			suffix := code[5:]
			return baseOrder + "." + suffix
		}
	}

	// Nếu code bắt đầu bằng số (như "4", "5.1", "6.1.1"), trả về chính nó
	if len(code) > 0 && (code[0] >= '0' && code[0] <= '9') {
		return code
	}

	// Mặc định trả về rỗng (sẽ là NULL)
	return ""
}

// SeedWorkStandards seeds initial work standards data
func SeedWorkStandards(ctx context.Context, db database.Database) error {
	// Tìm user có username = "admin" để dùng làm granted_by
	var adminID uuid.UUID
	err := db.GetContext(ctx, &adminID, `SELECT id FROM users WHERE username = 'admin' LIMIT 1`)
	if err != nil {
		logger.Info("SeedWorkStandards: Admin user not found, will set granted_by to NULL")
		adminID = uuid.Nil // Dùng Nil để đánh dấu là NULL
	} else {
		logger.Info("SeedWorkStandards: Found admin user", "user_id", adminID)
	}

	// Dữ liệu mẫu work standards (có thể tùy chỉnh theo nhu cầu)
	workStandards := []struct {
		ID            uuid.UUID
		Code          string
		Order         string
		Name          string
		ParentID      *uuid.UUID
		ScoreRequired bool
		Level         int
		Score         float64
		Description   string
	}{
		//3. Quan hệ làm việc (WORKING RELATIONSHIPS)
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Code:          "YT19",
			Order:         "3",
			Name:          "Quan hệ làm việc",
			ParentID:      nil,
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "WORKING RELATIONSHIPS",
		},

		{
			Code:          "YT19_1",
			Order:         "3.1",
			Name:          "Cần quan hệ bên trong Tổ, Nhóm, Ca",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000003"),
			ScoreRequired: true,
			Level:         1,
			Score:         1,
			Description:   "Need have the relation in the Team, Shift",
		},

		{
			Code:          "YT19_2",
			Order:         "3.2",
			Name:          "Cần quan hệ nội bộ giữa các Tổ, Nhóm, Ca trong Phòng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000003"),
			ScoreRequired: true,
			Level:         2,
			Score:         30.50,
			Description:   "Need internal relations between Teams, Shifts in Department  ",
		},

		{
			Code:          "YT19_3",
			Order:         "3.3",
			Name:          "Cần quan hệ nội bộ giữa các Phòng/ Ban",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000003"),
			ScoreRequired: true,
			Level:         3,
			Score:         60,
			Description:   "Need internal relations between the Departments/ Boards",
		},

		{
			Code:          "YT19_4",
			Order:         "3.4",
			Name:          "Cần quan hệ thường xuyên với bên ngoài (các cơ quan nhà nước…)",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000003"),
			ScoreRequired: true,
			Level:         3,
			Score:         89.50,
			Description:   "Need regular relations with the outside (the state agencies...)",
		},

		{
			Code:          "YT19_5",
			Order:         "3.5",
			Name:          "Đại diện cho Doanh nghiệp quan hệ thường xuyên với đối tác bên ngoài , khách hàng, doanh nhân khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000003"),
			ScoreRequired: true,
			Level:         3,
			Score:         119.00,
			Description:   "Enterprise represents regular relations with external partners, customers, other business",
		},
		// 4. TRÁCH NHIỆM (RESPONSIBILITIES)
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Code:          "4",
			Order:         "4",
			Name:          "Trách nhiệm",
			ParentID:      nil,
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Responsibilities",
		},
		// 4.1 Trách nhiệm về kết quả công việc của nhân viên
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000041"),
			Code:          "YT13",
			Order:         "4.1",
			Name:          "Trách nhiệm về kết quả công việc của nhân viên",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000004"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Responsibility for the results of the work of staff",
		},
		{
			Code:          "YT13_1",
			Order:         "4.1.1",
			Name:          "Chịu trách nhiệm công việc của chính mình",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "Responsible for their own work",
		},
		{
			Code:          "YT13_2",
			Order:         "4.1.2",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ 1 đến 10 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: false,
			Level:         2,
			Score:         63.50,
			Description:   "Responsible for the work of the department, team group or division / committee / branch from 1 to 10 people",
		},
		{
			Code:          "YT13_3",
			Order:         "4.1.3",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ 11 đến 50 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: false,
			Level:         3,
			Score:         122.50,
			Description:   "Responsible for the work of the department, team, group or division / committee / branch from 11 to 50 people",
		},
		{
			Code:          "YT13_4",
			Order:         "4.1.4",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ 51 đến 200 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: true,
			Level:         4,
			Score:         181.50,
			Description:   "Responsible for the work of the department, team, group or division/committee / branch from 51 to 200 people",
		},
		{
			Code:          "YT13_5",
			Order:         "4.1.5",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ 201 đến 500 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: false,
			Level:         5,
			Score:         240.50,
			Description:   "Responsible for the work of the department, team, group or division / committee / branch from 201 to 500 people",
		},
		{
			Code:          "YT13_6",
			Order:         "4.1.6",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ 501 đến 1000 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: true,
			Level:         6,
			Score:         299.50,
			Description:   "Responsible for the work of the department, team, group or division / committee / branch from 501 to 1000 people",
		},
		{
			Code:          "YT13_7",
			Order:         "4.1.7",
			Name:          "Chịu trách nhiệm công việc của Bộ phận, Tổ, Nhóm, Ca hoặc Phòng/ Ban/ Chi nhánh từ trên 1000 người",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000041"),
			ScoreRequired: false,
			Level:         7,
			Score:         358.50,
			Description:   "Responsible for the work of the department, team, group or division/committee / branch for over 1000 people",
		},
		// 4.2 Trách nhiệm về tài sản
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000042"),
			Code:          "YT14",
			Order:         "4.2",
			Name:          "Trách nhiệm về tài sản",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000004"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Property Responsibility",
		},
		{
			Code:          "YT14_1",
			Order:         "4.2.1",
			Name:          "Không phải chịu trách nhiệm quản lý tài sản",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: false,
			Level:         1,
			Score:         5.25,
			Description:   "Not responsible for asset management",
		},
		{
			Code:          "YT14_2",
			Order:         "4.2.2",
			Name:          "Chịu trách nhiệm quản lý một hoặc vài tài sản riêng lẻ, gián tiếp sử dụng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: true,
			Level:         2,
			Score:         76.05,
			Description:   "Responsible for management of whole the assets of enterprise",
		},
		{
			Code:          "YT14_3",
			Order:         "4.2.3",
			Name:          "Chịu trách nhiệm quản lý một hoặc vài tài sản riêng lẻ, trực tiếp sử dụng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: true,
			Level:         2,
			Score:         146.85,
		},
		{
			Code:          "YT14_4",
			Order:         "4.2.4",
			Name:          "Chịu trách nhiệm quản lý tài sản một Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: true,
			Level:         2,
			Score:         217.65,
		},
		{
			Code:          "YT14_5",
			Order:         "4.2.5",
			Name:          "Chịu trách nhiệm quản lý tài sản ít nhất  2 Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: true,
			Level:         2,
			Score:         288.45,
		},
		{
			Code:          "YT14_6",
			Order:         "4.2.6",
			Name:          " Chịu trách nhiệm quản lý tài sản toàn doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000042"),
			ScoreRequired: true,
			Level:         2,
			Score:         359.25,
		},
		// 4.3 Trách nhiệm tạo khả năng sinh lợi
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000043"),
			Code:          "YT15",
			Order:         "4.3",
			Name:          "Trách nhiệm tạo khả năng sinh lợi",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000004"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Responsible for creating profitability",
		},
		{
			Code:          "YT15_1",
			Order:         "4.3.1",
			Name:          "Không phải chịu trách nhiệm tạo khả năng sinh lợi một cách trực tiếp, mà chỉ có thể chịu trách nhiệm tạo khả năng sinh lợi một cách gián tiếp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000043"),
			ScoreRequired: false,
			Level:         1,
			Score:         6.75,
			Description:   "Not directly responsible for creating profitability, but only indirectly responsible for creating profitability",
		},
		{
			Code:          "YT15_2",
			Order:         "4.3.2",
			Name:          "Chịu trách nhiệm tạo khả năng sinh lợi trực tiếp của một Tổ, Nhóm, Ca",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000043"),
			ScoreRequired: false,
			Level:         2,
			Score:         95.25,
			Description:   "Responsible for generating profitable directly of a Team/ Group/ Shift",
		},
		{
			Code:          "YT15_3",
			Order:         "4.3.3",
			Name:          "Chịu trách nhiệm tạo sinh lợi trực tiếp trong phạm vi một Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000043"),
			ScoreRequired: false,
			Level:         3,
			Score:         183.75,
			Description:   "Responsible for creating profitable directly for a Department/ Board/ Branches",
		},
		{
			Code:          "YT15_4",
			Order:         "4.3.4",
			Name:          "Chịu trách nhiệm tạo khả năng sinh lợi cho ít nhất 2 Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000043"),
			ScoreRequired: false,
			Level:         4,
			Score:         272.25,
			Description:   "Responsible for creating profitable directly for at least 2 Department / Board/ Branches",
		},
		{
			Code:          "YT15_5",
			Order:         "4.3.5",
			Name:          "Chịu trách nhiệm tạo khả năng sinh lợi trực tiếp cho một doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000043"),
			ScoreRequired: true,
			Level:         5,
			Score:         360.75,
			Description:   "Responsible for creating profitable directly for a business",
		},
		// 4.4 Trách nhiệm kiểm soát tài chính
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000044"),
			Code:          "YT16",
			Order:         "4.4",
			Name:          "Trách nhiệm kiểm soát tài chính",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000004"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Responsibility for financial control",
		},
		{
			Code:          "YT16_1",
			Order:         "4.4.1",
			Name:          "Không phải chịu trách nhiệm kiểm soát tài chính",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000044"),
			ScoreRequired: false,
			Level:         1,
			Score:         6.75,
			Description:   "Not responsible for financial control",
		},
		{
			Code:          "YT16_2",
			Order:         "4.4.2",
			Name:          "Chịu trách nhiệm kiểm soát tài chính trong phạm vi Tổ, Nhóm, Ca",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000044"),
			ScoreRequired: false,
			Level:         2,
			Score:         95.25,
			Description:   "Responsible for financial control within a team/ shift",
		},
		{
			Code:          "YT16_3",
			Order:         "4.4.3",
			Name:          "Chịu trách nhiệm kiểm soát tài chính trong phạm vi một Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000044"),
			ScoreRequired: false,
			Level:         3,
			Score:         183.75,
			Description:   "Responsible for financial control within a Department / Board/ Branches",
		},
		{
			Code:          "YT16_4",
			Order:         "4.4.4",
			Name:          "Chịu trách nhiệm kiểm soát tài chính trong phạm vi ít nhất 2 Phòng/ Ban/ Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000044"),
			ScoreRequired: false,
			Level:         4,
			Score:         272.25,
			Description:   "Responsible for financial control within at least 2 Department/ Board/ Branches",
		},
		{
			Code:          "YT16_5",
			Order:         "4.4.5",
			Name:          "Chịu trách nhiệm kiểm soát tài chính cho một doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000044"),
			ScoreRequired: true,
			Level:         5,
			Score:         360.75,
			Description:   "Responsible for financial control for an enterprise",
		},
		// 4.5 Trách nhiệm ATLĐ, BHLĐ
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000045"),
			Code:          "YT17",
			Order:         "4.5",
			Name:          "Trách nhiệm ATLĐ, BHLĐ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000004"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Responsibility for occupational safety, labor protection",
		},
		{
			Code:          "YT17_1",
			Order:         "4.5.1",
			Name:          "Công việc không ảnh hưởng đến người khác về ATLĐ, BHLĐ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000045"),
			ScoreRequired: false,
			Level:         1,
			Score:         6.75,
			Description:   "The work does not affect others about occupational safety, labor protection",
		},
		{
			Code:          "YT17_2",
			Order:         "4.5.2",
			Name:          "Có ảnh hưởng, có thể gây TNLĐ cho người khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000045"),
			ScoreRequired: false,
			Level:         2,
			Score:         95.25,
			Description:   "Influential, can cause occupational accidents to others",
		},
		{
			Code:          "YT17_3",
			Order:         "4.5.3",
			Name:          "Đặc biệt ảnh hưởng đến NLĐ, dễ gây TNLĐ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000045"),
			ScoreRequired: false,
			Level:         3,
			Score:         183.75,
			Description:   "Particularly affect for employees, likely to cause occupational accidents",
		},
		{
			Code:          "YT17_4",
			Order:         "4.5.4",
			Name:          "Nguy hiểm, có thể ảnh hưởng đến tính mạng của NLĐ, người khác ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000045"),
			ScoreRequired: true,
			Level:         4,
			Score:         272.25,
			Description:   "Dangerous, may affect the lives of employees, and others (Especially dangerous, need to hight concentrate at work with the equipment (boiler, stamping machines, cutting machines, electric transformers ...) that affect the lives of employees, others) / or must bear legal responsibility for safety when the accident happened / or to be responsible in making policies on policies on labor protection and labor safety of the enterprise that can affect people other",
		},

		{
			Code:          "YT17_5",
			Order:         "4.5.5",
			Name:          "Đặc biệt nguy hiểm, cần tập trung cao độ nơi làm việc với các thiết bị (lò hơi, máy dập, máy cắt, biến thế điện...) có ảnh hưởng đến tính mạng của NLĐ, người khác/ hoặc phải chịu trách nhiệm pháp lý về ATLĐ khi tai nạn xảy ra/hoặc phải chịu trách nhiệm trong việc đưa ra các chính sách về BHLĐ và ATLĐ của Doanh nghiệp mà có thể gây ảnh hưởng đến người khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000045"),
			ScoreRequired: true,
			Level:         5,
			Score:         360.75,
			Description:   "Dangerous, may affect the lives of employees, and others (Especially dangerous, need to hight concentrate at work with the equipment (boiler, stamping machines, cutting machines, electric transformers ...) that affect the lives of employees, others) / or must bear legal responsibility for safety when the accident happened / or to be responsible in making policies on policies on labor protection and labor safety of the enterprise that can affect people other",
		},
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Code:          "5",
			Order:         "5",
			Name:          "Điều kiện làm việc",
			ParentID:      nil,
			ScoreRequired: true,
			Level:         0,
			Score:         0,
		},
		{
			Code:          "5.1",
			Order:         "5.1",
			Name:          "Phương tiện làm việc",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000005"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
		},
		{
			Code:          "5.1.1",
			Order:         "5.1.1",
			Name:          "Bàn làm việc , máy vi tính, điện thoại bàn và các trang thiết bị VPP cho công việc",
			ParentID:      nil, // Sẽ tự động resolve từ code "5.1"
			ScoreRequired: true,
			Level:         1,
			Score:         0,
		},
		{
			Code:          "5.1.2",
			Order:         "5.1.2",
			Name:          "Xe đưa rước theo tuyến cố định của công ty (nếu có)",
			ParentID:      nil, // Sẽ tự động resolve từ code "5.1"
			ScoreRequired: true,
			Level:         1,
			Score:         0,
		},
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000051"),
			Code:          "YT22",
			Order:         "5.2",
			Name:          "Thời gian làm việc",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000005"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
		},
		{
			Code:          "YT22_1",
			Order:         "5.2.1",
			Name:          "Làm việc theo giờ hành chánh (8 giờ / ngày)",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000051"),
			ScoreRequired: true,
			Level:         1,
			Score:         1.5,
		},
		{
			Code:          "YT22_2",
			Order:         "5.2.2",
			Name:          "Làm việc ngoài thị trường/ Phải thỉnh thoảng đi công tác ngoài tỉnh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000051"),
			ScoreRequired: true,
			Level:         2,
			Score:         40.83,
		},
		{
			Code:          "YT22_3",
			Order:         "5.2.3",
			Name:          "Làm việc theo ca (2 ca/ 8 giờ/ ngày) hoặc phải làm việc ngoài thị trường tương đương với 2 ca/ Phải thường xuyên đi công tác ngoài tỉnh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000051"),
			ScoreRequired: true,
			Level:         3,
			Score:         80.17,
		},
		{
			Code:          "YT22_4",
			Order:         "5.2.4",
			Name:          " Làm việc theo ca (3 ca/ 12 giờ/ ngày) hoặc phải làm việc ngoài thị trường tương đương với 3 ca/ Phải liên tục đi công tác ngoài tỉnh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000051"),
			ScoreRequired: true,
			Level:         4,
			Score:         119.50,
		},
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000052"),
			Code:          "YT20",
			Order:         "5.3",
			Name:          "Môi trường làm việc",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000005"),
			ScoreRequired: true,
			Level:         1,
			Score:         0,
			Description:   "Work Environment",
		},
		{
			Code:          "YT20_1",
			Order:         "5.3.1",
			Name:          "Bình thường",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000052"),
			ScoreRequired: false,
			Level:         1,
			Score:         1.5,
			Description:   "Normal",
		},
		{
			Code:          "YT20_2",
			Order:         "5.3.2",
			Name:          "Môi trường làm việc nặng nhọc (mức tiêu hao năng lượng cao) hoặc môi trường làm việc có áp lực về kinh doanh, tài chính, thị trường nhưng ở mức độ thấp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000052"),
			ScoreRequired: false,
			Level:         2,
			Score:         40.83,
			Description:   "Heavy working environment (high energy consumption) or working environment with less pressure on the business, financial, market but at a lower level",
		},
		{
			Code:          "YT20_3",
			Order:         "5.3.3",
			Name:          "Môi trường làm việc nặng nhọc, độc hại, nguy hiểm (công việc nặng nhọc, căng thẳng thị giác, tiếng ồn, bụi bẩn, hơi độc …) hoặc môi trường làm việc chịu áp lực về kinh doanh, tài chính, thị trường nhưng ở mức độ trung bình",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000052"),
			ScoreRequired: true,
			Level:         3,
			Score:         80.17,
			Description:   "Heavy work environment, harmful, dangerous (heavy work, visual stress, noise, dust, toxic fumes...) or the working environment of the pressure on the business, financial, market but the average level",
		},
		{
			Code:          "YT20_4",
			Order:         "5.3.4",
			Name:          "Môi trường làm việc đặc biệt nặng nhọc, độc hại, nguy hiểm (công việc nặng nhọc, nguy hiểm đặc biệt, tiếng ồn, nhiễm độc, áp lực công việc…) hoặc môi trường làm việc chịu nhiều áp lực về kinh doanh, tài chính, thị trường nhưng ở mức độ cao",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000052"),
			ScoreRequired: true,
			Level:         4,
			Score:         119.50,
			Description:   "Working environment extremely heavy, hazardous (heavy work, especially dangerous, noise, poisoning, work pressure...) or the working environment must be more very high pressure on the business, financial, market but at a high level",
		},
		// 5.4 Cường độ tập trung
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000053"),
			Code:          "YT18",
			Order:         "5.4",
			Name:          "Cường độ tập trung",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000005"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "The intensity of focus",
		},
		{
			Code:          "YT18_1",
			Order:         "5.4.1",
			Name:          "Không cần nỗ lực đặc biệt",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000053"),
			ScoreRequired: false,
			Level:         1,
			Score:         2,
			Description:   "No special effort required",
		},
		{
			Code:          "YT18_2",
			Order:         "5.4.2",
			Name:          "Cần tập trung quan sát, lắng nghe thường xuyên",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000053"),
			ScoreRequired: false,
			Level:         2,
			Score:         61.00,
			Description:   "Need to focus on observation and listening regularly",
		},
		{
			Code:          "YT18_3",
			Order:         "5.4.3",
			Name:          "Phải nỗ lực tập trung để quan sát, lắng nghe và phân tích được các cấp độ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000053"),
			ScoreRequired: true,
			Level:         3,
			Score:         120.00,
			Description:   "Must make efforts to concentrate on observation, listening and analyzing at various levels",
		},
		// 5.5 Rủi ro tai nạn lao động
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000054"),
			Code:          "YT21",
			Order:         "5.5",
			Name:          "Rủi ro tai nạn lao động, rủi ro nghề nghiệp và bệnh nghề nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000005"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Occupational accident risks, occupational risks and occupational diseases",
		},
		{
			Code:          "YT21_1",
			Order:         "5.5.1",
			Name:          "Không có rủi ro đáng kể",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000054"),
			ScoreRequired: false,
			Level:         1,
			Score:         1.5,
			Description:   "No significant risks",
		},
		{
			Code:          "YT21_2",
			Order:         "5.5.2",
			Name:          "Rủi ro về pháp lý của nghề nghiệp nhưng ở mức độ thấp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000054"),
			ScoreRequired: false,
			Level:         2,
			Score:         40.83,
			Description:   "Legal risks of the profession but at a low level",
		},
		{
			Code:          "YT21_3",
			Order:         "5.5.3",
			Name:          "Rủi ro nghề nghiệp hoặc bệnh nghề nghiệp dễ xảy ra hơn bình thường hoặc mắc bệnh nghề nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000054"),
			ScoreRequired: true,
			Level:         3,
			Score:         80.17,
			Description:   "Occupational risks or occupational diseases are more likely than normal or suffer from occupational diseases",
		},
		{
			Code:          "YT21_4",
			Order:         "5.5.4",
			Name:          "Rủi ro TNLĐ (lao động, giao thông) cao, thường xuyên phải tiếp xúc trực tiếp đến sức khỏe, ảnh hưởng đến sức khỏe, bệnh nghề nghiệp ở mức độ cao",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000054"),
			ScoreRequired: true,
			Level:         4,
			Score:         119.50,
			Description:   "High occupational accident risks (labor, traffic), frequent direct contact affecting health, occupational diseases at high levels",
		},
		// 6. YÊU CẦU VỀ NĂNG LỰC NHẬN THỨC
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			Code:          "6",
			Order:         "6",
			Name:          "Các yêu cầu cần có về năng lực nhận thức cho vị trí công việc này",
			ParentID:      nil,
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Requirements of cognitive capacity for this job position",
		},
		// 6.1 Đặc điểm cá nhân
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000061"),
			Code:          "6.1",
			Order:         "6.1",
			Name:          "Đặc điểm cá nhân",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000006"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Personal traits",
		},
		{
			Code:          "6.1.1",
			Order:         "6.1.1",
			Name:          "Giới tính",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000061"),
			ScoreRequired: false,
			Level:         1,
			Score:         0,
			Description:   "Sex (Male or Female)",
		},
		{
			Code:          "6.1.2",
			Order:         "6.1.2",
			Name:          "Độ tuổi",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000061"),
			ScoreRequired: false,
			Level:         1,
			Score:         0,
			Description:   "Age (e.g., At least 40 years old)",
		},
		{
			Code:          "6.1.3",
			Order:         "6.1.3",
			Name:          "Ngoại hình",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000061"),
			ScoreRequired: false,
			Level:         1,
			Score:         0,
			Description:   "Appearance (e.g., Good Looking)",
		},
		{
			Code:          "6.1.4",
			Order:         "6.1.4",
			Name:          "Sức khỏe",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000061"),
			ScoreRequired: false,
			Level:         1,
			Score:         0,
			Description:   "Health (e.g., Good)",
		},
		// 6.2 Trình độ học vấn/chuyên môn
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000062"),
			Code:          "NT01",
			Order:         "6.2",
			Name:          "Trình độ học vấn/chuyên môn",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000006"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Education/Profession",
		},
		{
			Code:          "YT1_1",
			Order:         "6.2.1",
			Name:          "Tốt nghiệp cấp 1, cấp 2",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: false,
			Level:         1,
			Score:         15,
			Description:   "Elementary or junior high school graduate",
		},
		{
			Code:          "YT1_2",
			Order:         "6.2.2",
			Name:          "Tốt nghiệp cấp 3",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: false,
			Level:         2,
			Score:         15,
			Description:   "High school graduate",
		},
		{
			Code:          "YT1_3",
			Order:         "6.2.3",
			Name:          "Đào tạo nghề hoặc kỹ thuật ngắn hạn",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: false,
			Level:         3,
			Score:         15,
			Description:   "Vocational training or short-term technical training",
		},
		{
			Code:          "YT1_4",
			Order:         "6.2.4",
			Name:          "Trung học chuyên nghiệp, Cao đẳng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: false,
			Level:         4,
			Score:         678.75,
			Description:   "Secondary school or College graduate",
		},
		{
			Code:          "YT1_5",
			Order:         "6.2.5",
			Name:          "Kỹ sư, cử nhân",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: true,
			Level:         5,
			Score:         900,
			Description:   "Engineering or Bachelor",
		},
		{
			Code:          "YT1_6",
			Order:         "6.2.6",
			Name:          "Thạc sĩ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: true,
			Level:         6,
			Score:         1121.25,
			Description:   "Masters",
		},
		{
			Code:          "YT1_7",
			Order:         "6.2.7",
			Name:          "Tiến sĩ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: true,
			Level:         7,
			Score:         1342.50,
			Description:   "PhD",
		},
		{
			Code:          "YT1_8",
			Order:         "6.2.8",
			Name:          "Chuyên ngành/Quản trị Kinh doanh/Khối kinh tế/Hàng hải/Giao thông vận tải",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000062"),
			ScoreRequired: true,
			Level:         8,
			Score:         1563.75,
			Description:   "Specialization/Business administration/Economic Division/Shipping/Transportation",
		},
		// 6.3 Trình độ ngoại ngữ - Anh văn
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000063"),
			Code:          "NT03",
			Order:         "6.3",
			Name:          "Trình độ ngoại ngữ - Anh văn",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000006"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Foreign Language skill - English",
		},
		{
			Code:          "YT3_1",
			Order:         "6.3.1",
			Name:          "Không cần thiết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000063"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "Unnecessary",
		},
		{
			Code:          "YT3_2",
			Order:         "6.2.2",
			Name:          "Hiểu, đọc, viết/chứng chỉ A, B",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000063"),
			ScoreRequired: false,
			Level:         2,
			Score:         78.25,
			Description:   "Can understand, read and write the certificates A, B",
		},
		{
			Code:          "YT3_3",
			Order:         "6.2.3",
			Name:          "Hiểu, đọc, viết, nói/chứng chỉ C",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000063"),
			ScoreRequired: true,
			Level:         3,
			Score:         152.00,
			Description:   "Can understand, read and write the certificates C",
		},
		{
			Code:          "YT3_4",
			Order:         "6.2.4",
			Name:          "Ngoại ngữ chuyên ngành + Chứng chỉ C trở lên",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000063"),
			ScoreRequired: true,
			Level:         4,
			Score:         225.75,
			Description:   "Specialized foreign language and the certificates C or higher",
		},
		{
			Code:          "YT3_5",
			Order:         "6.2.5",
			Name:          "Trên một ngoại ngữ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000063"),
			ScoreRequired: true,
			Level:         5,
			Score:         299.50,
			Description:   "Know more than 1 foreign languages",
		},
		// 6.4 Trình độ tin học, thiết bị kỹ thuật
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000064"),
			Code:          "NT02",
			Order:         "6.4",
			Name:          "Trình độ tin học, thiết bị kỹ thuật",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000006"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Computer skill, technical equipment",
		},
		{
			Code:          "YT2_1",
			Order:         "6.4.1",
			Name:          "Không cần thiết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000064"),
			ScoreRequired: false,
			Level:         1,
			Score:         6,
			Description:   "Unnecessary",
		},
		{
			Code:          "YT2_2",
			Order:         "6.4.2",
			Name:          "Vi tính văn phòng (Word, Excel, PowerPoint, Visio, Outlook), các loại máy photo, fax, điện thoại, tổng đài, sử dụng thành thạo internet, email,...",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000064"),
			ScoreRequired: false,
			Level:         2,
			Score:         104.33,
			Description:   "Office computer (Word, Excel, PowerPoint, Visio, Outlook), photo, fax, and phone",
		},
		{
			Code:          "YT2_3",
			Order:         "6.2.3",
			Name:          "Vi tính chuyên ngành (AutoCad, Macro Of Excel, Access...)",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000064"),
			ScoreRequired: true,
			Level:         3,
			Score:         202.67,
			Description:   "Specialized computer (AutoCad, Macro Of Excel, Access...)",
		},
		{
			Code:          "YT2_4",
			Order:         "6.2.4",
			Name:          "Khả năng lập trình công nghệ thông tin",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000064"),
			ScoreRequired: true,
			Level:         4,
			Score:         301.00,
			Description:   "Programming capabilities of information technology",
		},
		// 6.5 Số năm kinh nghiệm
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000065"),
			Code:          "NT04",
			Order:         "6.5",
			Name:          "Số năm kinh nghiệm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000006"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Years of working experience",
		},
		{
			Code:          "YT4_1",
			Order:         "6.5.1",
			Name:          "Từ 0 đến 6 tháng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000065"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "From 0 to 6 months",
		},
		{
			Code:          "YT4_2",
			Order:         "6.5.2",
			Name:          "Từ > 6 tháng đến 3 năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000065"),
			ScoreRequired: false,
			Level:         2,
			Score:         78.25,
			Description:   "From > 6 months to 3 years",
		},
		{
			Code:          "YT4_3",
			Order:         "6.5.3",
			Name:          "Từ > 3 năm - 6 năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000065"),
			ScoreRequired: false,
			Level:         3,
			Score:         152.00,
			Description:   "From > 3 years to 6 years",
		},
		{
			Code:          "YT4_4",
			Order:         "6.5.4",
			Name:          "Từ > 6 năm - 9 năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000065"),
			ScoreRequired: true,
			Level:         4,
			Score:         255.75,
			Description:   "From > 6 years to 9 years (At least 08 years of working experience in the same position)",
		},
		{
			Code:          "YT4_5",
			Order:         "6.5.5",
			Name:          "Từ > 9 năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000065"),
			ScoreRequired: true,
			Level:         5,
			Score:         299.50,
			Description:   "From > 9 years",
		},
		// 7. YÊU CẦU VỀ NĂNG LỰC HÀNH NGHỀ
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000007"),
			Code:          "7",
			Order:         "7",
			Name:          "Các yêu cầu cần có về năng lực hành nghề cho vị trí công việc này",
			ParentID:      nil,
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Requirements of practice capacity for this job position",
		},
		// 7.1 Lập kế hoạch
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000071"),
			Code:          "NL01",
			Order:         "7.1",
			Name:          "Lập kế hoạch",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Planning",
		},
		{
			Code:          "YT5_1",
			Order:         "7.1.1",
			Name:          "Không cần",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000071"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "Unnecessary",
		},
		{
			Code:          "YT5_2",
			Order:         "7.1.2",
			Name:          "Lập kế hoạch tuần, tháng",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000071"),
			ScoreRequired: false,
			Level:         2,
			Score:         78.25,
			Description:   "Week, month planning",
		},
		{
			Code:          "YT5_3",
			Order:         "7.1.3",
			Name:          "Lập kế hoạch 3 tháng/6 tháng/năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000071"),
			ScoreRequired: true,
			Level:         3,
			Score:         152.00,
			Description:   "3 months/6 months/year planning",
		},
		{
			Code:          "YT5_4",
			Order:         "7.1.4",
			Name:          "Lập kế hoạch chiến lược năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000071"),
			ScoreRequired: true,
			Level:         4,
			Score:         225.75,
			Description:   "Strategic planning in year",
		},
		{
			Code:          "YT5_5",
			Order:         "7.1.5",
			Name:          "Lập kế hoạch chiến lược trên 1 năm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000071"),
			ScoreRequired: true,
			Level:         5,
			Score:         299.50,
			Description:   "Strategic planning over year",
		},
		// 7.2 Hiểu biết
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000072"),
			Code:          "NL02",
			Order:         "7.2",
			Name:          "Hiểu biết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Understanding",
		},
		{
			Code:          "YT6_1",
			Order:         "7.2.1",
			Name:          "Có thể hiểu rõ các mệnh lệnh, chỉ thị liên quan đến công việc",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000072"),
			ScoreRequired: false,
			Level:         1,
			Score:         3,
			Description:   "Understand the orders and directives relating to work",
		},
		{
			Code:          "YT6_2",
			Order:         "7.2.2",
			Name:          "Có thể hiểu rõ các mệnh lệnh và chỉ thị liên quan tới công việc của Phòng/Ban/Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000072"),
			ScoreRequired: false,
			Level:         2,
			Score:         52.17,
			Description:   "Understand the orders and directives relating to works of Department/Board/Branches",
		},
		{
			Code:          "YT6_3",
			Order:         "7.2.3",
			Name:          "Có thể nắm được bản chất thông tin mới liên quan đến công việc, truyền đạt được trong Phòng/Ban/Chi nhánh để giải quyết công việc",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000072"),
			ScoreRequired: true,
			Level:         3,
			Score:         101.34,
			Description:   "To grasp the essence of new informations relating to works, transmitting in Department/Board/Branches for solving task",
		},
		{
			Code:          "YT6_4",
			Order:         "7.2.4",
			Name:          "Có thể liên kết thông tin, khái quát nhiều nguồn thông tin để nhận biết, hiểu rõ, phân tích giải quyết vấn đề và đào tạo được người khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000072"),
			ScoreRequired: true,
			Level:         4,
			Score:         150.51,
			Description:   "Information link, generalize information resources to recognize, understand, analyse, solving problems and training others",
		},
		// 7.3 Phán quyết – Ra quyết định
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000073"),
			Code:          "NL03",
			Order:         "7.3",
			Name:          "Phán quyết – Ra quyết định",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Judgement - Decision making",
		},
		{
			Code:          "YT7_1",
			Order:         "7.3.1",
			Name:          "Công việc không cần phán quyết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000073"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "Work without judgement",
		},
		{
			Code:          "YT7_2",
			Order:         "7.3.2",
			Name:          "Phải phán quyết các điểm nhỏ trong phạm vi các chỉ thị tương đối chi tiết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000073"),
			ScoreRequired: false,
			Level:         2,
			Score:         78.25,
			Description:   "Must judge small points within limits of relative detail directives",
		},
		{
			Code:          "YT7_3",
			Order:         "7.3.3",
			Name:          "Khi có các hướng dẫn/chỉ thị chung, phải đưa ra các quyết định tác động tới kết quả làm việc của Phòng/Ban/Chi nhánh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000073"),
			ScoreRequired: true,
			Level:         3,
			Score:         152.00,
			Description:   "Make decisions effecting to working result of department/Board/Branches when receive general instructions and directives",
		},
		{
			Code:          "YT7_4",
			Order:         "7.3.4",
			Name:          "Khi có các hướng dẫn/chỉ thị chung, phải đưa ra các quyết định tác động tới kết quả làm việc của các Phòng/Ban/Chi nhánh và khối",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000073"),
			ScoreRequired: true,
			Level:         4,
			Score:         225.75,
			Description:   "Make decisions effecting to working result of Departments and division when receive general instructions and directives",
		},
		{
			Code:          "YT7_5",
			Order:         "7.3.5",
			Name:          "Khi có các hướng dẫn/chỉ thị chung, phải đưa ra các quyết định tác động tới kết quả làm việc của cả doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000073"),
			ScoreRequired: true,
			Level:         5,
			Score:         299.50,
			Description:   "Make decisions effecting to working result of whole company when receive general instructions and directives",
		},
		// 7.4 Khả năng thuyết phục
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000074"),
			Code:          "NL04",
			Order:         "7.4",
			Name:          "Khả năng thuyết phục",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Persuasive Ability",
		},
		{
			Code:          "YT8_1",
			Order:         "7.4.1",
			Name:          "Công việc không cần thuyết phục người khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000074"),
			ScoreRequired: false,
			Level:         1,
			Score:         3,
			Description:   "The work without persuade others",
		},
		{
			Code:          "YT8_2",
			Order:         "7.4.2",
			Name:          "Cần thuyết phục các thành viên trong Phòng/Ban/Chi nhánh, Tổ/nhóm, cấp dưới",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000074"),
			ScoreRequired: false,
			Level:         2,
			Score:         52.17,
			Description:   "Need to persuade members in Department, team/group, subordinate",
		},
		{
			Code:          "YT8_3",
			Order:         "7.4.3",
			Name:          "Cần thuyết phục các thành viên trong và ngoài Phòng/Ban/Chi nhánh, Tổ/nhóm, các khối khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000074"),
			ScoreRequired: true,
			Level:         3,
			Score:         101.34,
			Description:   "Need to persuade inside and outside members of department, team, group, other divisions",
		},
		{
			Code:          "YT8_4",
			Order:         "7.4.4",
			Name:          "Cần thuyết phục ban lãnh đạo, đối tác, khách hàng khó tính",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000074"),
			ScoreRequired: true,
			Level:         4,
			Score:         150.1,
			Description:   "Need to persuade executive board, partners, fastidious customers",
		},
		// 7.5 Tính sáng tạo
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000075"),
			Code:          "NL05",
			Order:         "7.5",
			Name:          "Tính sáng tạo",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Creativity",
		},
		{
			Code:          "YT9_1",
			Order:         "7.5.1",
			Name:          "Không cần tính sáng tạo",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000075"),
			ScoreRequired: false,
			Level:         1,
			Score:         1.5,
			Description:   "Without creativity",
		},
		{
			Code:          "YT9_2",
			Order:         "7.5.2",
			Name:          "Tạo ra những cải tiến nhỏ trong phạm vi công việc của Phòng/Ban/Chi nhánh, Tổ/nhóm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000075"),
			ScoreRequired: false,
			Level:         2,
			Score:         38.38,
			Description:   "Create small innovations within job scope of department, team, group",
		},
		{
			Code:          "YT9_3",
			Order:         "7.5.3",
			Name:          "Tạo ra những cải tiến, giá trị gia tăng có thể áp dụng cho các Phòng/Ban/Chi nhánh/Khối",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000075"),
			ScoreRequired: true,
			Level:         3,
			Score:         75.26,
			Description:   "Create improvements and added values applied for Department/Division",
		},
		{
			Code:          "YT9_4",
			Order:         "7.5.4",
			Name:          "Tạo ra những quy trình mới, những sản phẩm mới mang lại hiệu quả cao hơn",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000075"),
			ScoreRequired: true,
			Level:         4,
			Score:         112.14,
			Description:   "Create new processes, new products that bring greater efficiency",
		},
		{
			Code:          "YT9_5",
			Order:         "7.5.5",
			Name:          "Tạo ra các ý tưởng mới về loại hình kinh doanh, về quản lý, tổ chức",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000075"),
			ScoreRequired: true,
			Level:         5,
			Score:         149.02,
			Description:   "Create new ideas on types of business, management, organization",
		},
		// 7.6 Nhận diện vấn đề
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000076"),
			Code:          "NL06",
			Order:         "7.6",
			Name:          "Nhận diện vấn đề",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Identifying the problem",
		},
		{
			Code:  "YT10_1",
			Order: "7.6.1",

			Name:          "Không cần thiết",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000076"),
			ScoreRequired: false,
			Level:         1,
			Score:         6,
			Description:   "Unnecessary",
		},
		{
			Code:          "YT10_2",
			Order:         "7.6.2",
			Name:          "Cần tư duy để phân biệt, nhưng ở mức độ đơn giản",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000076"),
			ScoreRequired: false,
			Level:         2,
			Score:         104.33,
			Description:   "Thinking to distinguish, but at simple level",
		},
		{
			Code:          "YT10_3",
			Order:         "7.6.3",
			Name:          "Cần tư duy nhiều, suy nghĩ sâu để phân biệt",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000076"),
			ScoreRequired: true,
			Level:         3,
			Score:         202.67,
			Description:   "Need more thinking, deep thinking to distinguish",
		},
		{
			Code:          "YT10_4",
			Order:         "7.6.4",
			Name:          "Đòi hỏi sự tìm tòi, phân tích, phán đoán, suy luận",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000076"),
			ScoreRequired: true,
			Level:         4,
			Score:         301.00,
			Description:   "Requirement of finding stuff, analysis, judgement, inference",
		},
		// 7.7 Huấn luyện đào tạo
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000077"),
			Code:          "NL07",
			Order:         "7.7",
			Name:          "Huấn luyện đào tạo",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Coaching - Training",
		},
		{
			Code:          "YT11_1",
			Order:         "7.7.1",
			Name:          "Thực hiện theo quy trình, không cần phải đào tạo người khác",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000077"),
			ScoreRequired: false,
			Level:         1,
			Score:         3,
			Description:   "Follow the process, do not need to train another person",
		},
		{
			Code:          "YT11_2",
			Order:         "7.7.2",
			Name:          "Huấn luyện được các thành viên trong Phòng/Ban/Chi nhánh, Tổ, nhóm, ca (từ 2-20 người)",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000077"),
			ScoreRequired: false,
			Level:         2,
			Score:         52.17,
			Description:   "Train members in Department, team, group (from 2-20 people)",
		},
		{
			Code:          "YT11_3",
			Order:         "7.7.3",
			Name:          "Huấn luyện được các thành viên trong một Phòng/Ban/Chi nhánh/Bộ phận nghiệp vụ, dự án",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000077"),
			ScoreRequired: true,
			Level:         3,
			Score:         101.34,
			Description:   "Train members in Department, project",
		},
		{
			Code:          "YT11_4",
			Order:         "7.7.4",
			Name:          "Soạn thảo và huấn luyện được các thành viên trong một Phòng/Ban/Chi nhánh/Bộ phận nghiệp vụ, dự án",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000077"),
			ScoreRequired: true,
			Level:         4,
			Score:         150.51,
			Description:   "Compile and coaching members in Department, project",
		},
		// 7.8 Năng lực lãnh đạo
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000078"),
			Code:          "NL08",
			Order:         "7.8",
			Name:          "Năng lực lãnh đạo",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000007"),
			ScoreRequired: true,
			Level:         0,
			Score:         0,
			Description:   "Leadership abilities",
		},
		{
			Code:          "YT12_1",
			Order:         "7.8.1",
			Name:          "Không cần năng lực lãnh đạo",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000078"),
			ScoreRequired: false,
			Level:         1,
			Score:         4.5,
			Description:   "No need to demonstrate leadership",
		},
		{
			Code:          "YT12_2",
			Order:         "7.8.2",
			Name:          "Phải lãnh đạo một Tổ, Nhóm, Ca",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000078"),
			ScoreRequired: false,
			Level:         2,
			Score:         78.25,
			Description:   "Must lead a group or team",
		},
		{
			Code:          "YT12_3",
			Order:         "7.8.3",
			Name:          "Phải lãnh đạo ít nhất một Phòng hoặc ít nhất theo một chức năng hoặc theo một tuyến sản phẩm",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000078"),
			ScoreRequired: true,
			Level:         3,
			Score:         152.00,
			Description:   "Have led at least one Department or at least a function or a product line",
		},
		{
			Code:          "YT12_4",
			Order:         "7.8.4",
			Name:          "Phải lãnh đạo ít nhất hai Phòng hoặc ít nhất theo hai chức năng hoặc theo hai tuyến sản phẩm trực thuộc doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000078"),
			ScoreRequired: true,
			Level:         4,
			Score:         225.75,
			Description:   "Have led at least two Departments or at least 2 functions or 2 product lines belong to the enterprise",
		},
		{
			Code:          "YT12_5",
			Order:         "7.8.5",
			Name:          "Phải lãnh đạo một doanh nghiệp",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000078"),
			ScoreRequired: true,
			Level:         5,
			Score:         299.50,
			Description:   "Must lead a business",
		},
		// 8. YÊU CẦU KHÁC
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000008"),
			Code:          "8",
			Order:         "8",
			Name:          "Các yêu cầu khác, nếu có",
			ParentID:      nil,
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Other requirements, if any",
		},
		{
			Code:          "8.1",
			Order:         "8.1",
			Name:          "Nhanh nhẹn, điềm tĩnh, chịu áp lực công việc, có tinh thần trách nhiệm cao, có ý thức chấp hành kỷ luật,...",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000008"),
			ScoreRequired: false,
			Level:         1,
			Score:         0,
			Description:   "Agile, calm, accept high working pressure, high sense of responsibility, a sense of observance of discipline",
		},
		{
			Code:          "8.2",
			Order:         "8.2",
			Name:          "Có kiến thức về hệ thống ISO, QMS, ATLĐ, Môi trường, ERP,...",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000008"),
			ScoreRequired: false,
			Level:         2,
			Score:         0,
			Description:   "Knowledge systems ISO, QMS, Labour Safety, Environment, ERP",
		},
		{
			Code:          "8.3",
			Order:         "8.3",
			Name:          "Hiểu biết về ngành hàng hải",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000008"),
			ScoreRequired: false,
			Level:         3,
			Score:         0,
			Description:   "Understanding maritime industry",
		},
		{
			Code:          "8.4",
			Order:         "8.4",
			Name:          "Có kiến thức về văn thư lưu trữ",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000008"),
			ScoreRequired: false,
			Level:         4,
			Score:         0,
			Description:   "Knowledge of the archives",
		},
		{
			Code:          "8.5",
			Order:         "8.5",
			Name:          "Đi công tác theo yêu cầu",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000008"),
			ScoreRequired: false,
			Level:         5,
			Score:         0,
			Description:   "Business travel as required",
		},
		// 9. CẠNH TRANH THỊ TRƯỜNG LAO ĐỘNG
		{
			ID:            uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			Code:          "YT23",
			Order:         "9",
			Name:          "Cạnh tranh của thị trường lao động, nếu có",
			ParentID:      nil,
			ScoreRequired: false,
			Level:         0,
			Score:         0,
			Description:   "Competition of the labor market, if any",
		},
		{
			Code:          "YT23_1",
			Order:         "9.1",
			Name:          "Thị trường không cạnh tranh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: false,
			Level:         1,
			Score:         2.5,
			Description:   "The market is not competitive",
		},
		{
			Code:          "YT23_2",
			Order:         "9.2",
			Name:          "Thị trường ít cạnh tranh",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: false,
			Level:         2,
			Score:         72.50,
			Description:   "The market is less competitive",
		},
		{
			Code:          "YT23_3",
			Order:         "9.3",
			Name:          "Thị trường cạnh tranh vừa",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: false,
			Level:         3,
			Score:         226.50,
			Description:   "Medium competitive market",
		},
		{
			Code:          "YT23_4",
			Order:         "9.4",
			Name:          "Thị trường cạnh tranh nhiều",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: false,
			Level:         4,
			Score:         380.50,
			Description:   "Market more competitive",
		},
		{
			Code:          "YT23_5",
			Order:         "9.5",
			Name:          "Thị trường cạnh tranh rất nhiều",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: false,
			Level:         5,
			Score:         534.50,
			Description:   "Competitive market lot",
		},
		{
			Code:          "YT23_6",
			Order:         "9.6",
			Name:          "Thị trường cạnh tranh gay gắt",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: true,
			Level:         6,
			Score:         688.50,
			Description:   "The market is fiercely competitive",
		},
		{
			Code:          "YT23_7",
			Order:         "9.7",
			Name:          "Thị trường cạnh tranh khốc liệt (chảy máu chất xám)",
			ParentID:      uuidPtr("00000000-0000-0000-0000-000000000009"),
			ScoreRequired: true,
			Level:         7,
			Score:         1002.50,
			Description:   "The competitive market is very fierce (brain drain)",
		},
	}

	// Map để lưu code -> ID cho việc resolve ParentID
	codeToIDMap := make(map[string]uuid.UUID)
	// Map để lưu UUID -> code (để tìm parent code từ UUID cố định)
	uuidToCodeMap := make(map[uuid.UUID]string)

	// Đếm số mục đã seed thành công và số mục bị lỗi
	successCount := 0
	errorCount := 0
	skipCount := 0

	// Seed từng work standard
	for _, ws := range workStandards {
		// Kiểm tra xem đã tồn tại chưa (theo ID hoặc code) - giống như type_competency_seed.go
		var existingID uuid.UUID
		var checkID uuid.UUID = ws.ID
		if checkID == uuid.Nil {
			// Nếu không có ID cố định, chỉ kiểm tra theo code
			err := db.GetContext(ctx, &existingID, `SELECT id FROM work_standards WHERE code = $1 LIMIT 1`, ws.Code)
			if err == nil && existingID != uuid.Nil {
				// Lưu vào map để dùng cho các work standard khác
				codeToIDMap[ws.Code] = existingID
				skipCount++
				continue // Bỏ qua nếu đã tồn tại
			} else if err != nil && err != sql.ErrNoRows {
				// Nếu có lỗi khác ngoài "không tìm thấy", log và tiếp tục
			}
		} else {
			// Nếu có ID cố định, kiểm tra cả ID và code
			err := db.GetContext(ctx, &existingID, `SELECT id FROM work_standards WHERE id = $1 OR code = $2 LIMIT 1`, checkID, ws.Code)
			if err == nil && existingID != uuid.Nil {
				// Lưu vào map để dùng cho các work standard khác
				codeToIDMap[ws.Code] = existingID
				uuidToCodeMap[existingID] = ws.Code
				// Nếu ID thực tế khác với ID cố định, cũng lưu mapping từ ID cố định
				if existingID != checkID {
					uuidToCodeMap[checkID] = ws.Code
				}
				skipCount++
				continue // Bỏ qua nếu đã tồn tại
			} else if err != nil && err != sql.ErrNoRows {
				// Nếu có lỗi khác ngoài "không tìm thấy", log và tiếp tục
			}
		}

		// Xác định ID: sử dụng ID đã định nghĩa sẵn nếu có, nếu không thì tạo UUID mới
		var workStandardID uuid.UUID
		if ws.ID != uuid.Nil {
			workStandardID = ws.ID
		} else {
			workStandardID = uuid.New()
		}

		// Sử dụng ParentID đã được định nghĩa sẵn
		// Nếu ParentID = nil, chỉ tìm trong codeToIDMap (parent đã được seed trước đó)
		resolvedParentID := ws.ParentID
		if resolvedParentID == nil {
			// Tìm parent từ code (ví dụ: "5.1.1" -> parent là "5.1")
			lastDotIndex := -1
			for i := len(ws.Code) - 1; i >= 0; i-- {
				if ws.Code[i] == '.' {
					lastDotIndex = i
					break
				}
			}
			if lastDotIndex > 0 {
				parentCode := ws.Code[:lastDotIndex]
				if parentID, exists := codeToIDMap[parentCode]; exists {
					resolvedParentID = &parentID
				}
			}
		}

		var query string
		var args []interface{}

		// Xử lý description: nếu rỗng thì set NULL
		var descriptionPtr *string
		if ws.Description != "" {
			descriptionPtr = &ws.Description
		}

		// Xử lý order: nếu có trong data thì dùng, nếu không thì tính từ code
		orderValue := ws.Order
		if orderValue == "" {
			orderValue = codeToOrder(ws.Code)
		}
		var orderPtr *string
		if orderValue != "" {
			orderPtr = &orderValue
		}

		// Nếu score khác 0 thì tự động set scoreRequired = true
		scoreRequired := ws.ScoreRequired
		if ws.Score != 0 {
			scoreRequired = true
		}

		// Nếu có admin user thì dùng, không thì set NULL
		// Dữ liệu seeding luôn có is_default = true
		if adminID != uuid.Nil {
			query = `INSERT INTO work_standards (id, code, name, parent_id, score_require, level, score, description, "order", grandted_by, is_default) 
			         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
			args = []interface{}{
				workStandardID,
				ws.Code,
				ws.Name,
				resolvedParentID,
				scoreRequired,
				ws.Level,
				ws.Score,
				descriptionPtr,
				orderPtr,
				adminID,
				true, // is_default = true cho dữ liệu seeding
			}
		} else {
			query = `INSERT INTO work_standards (id, code, name, parent_id, score_require, level, score, description, "order", grandted_by, is_default) 
			         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULL, $10)`
			args = []interface{}{
				workStandardID,
				ws.Code,
				ws.Name,
				resolvedParentID,
				scoreRequired,
				ws.Level,
				ws.Score,
				descriptionPtr,
				orderPtr,
				true, // is_default = true cho dữ liệu seeding
			}
		}

		err = db.ExecContext(ctx, query, args...)
		if err != nil {
			logger.Error("SeedWorkStandards: Failed to create work standard",
				"code", ws.Code,
				"name", ws.Name,
				"id", workStandardID,
				"parent_id", resolvedParentID,
				"error", err)
			// Log thêm thông tin để debug
			if resolvedParentID != nil {
				var parentCode string
				if code, exists := uuidToCodeMap[*resolvedParentID]; exists {
					parentCode = code
				} else {
					// Thử tìm parent code từ database
					var parentCodeFromDB string
					if err2 := db.GetContext(ctx, &parentCodeFromDB, `SELECT code FROM work_standards WHERE id = $1 LIMIT 1`, *resolvedParentID); err2 == nil {
						parentCode = parentCodeFromDB
					}
				}
				logger.Error("SeedWorkStandards: Parent info",
					"parent_id", *resolvedParentID,
					"parent_code", parentCode)
			}
			// Log thêm thông tin về query để debug
			logger.Error("SeedWorkStandards: Query details",
				"query", query,
				"args_count", len(args))
			// Tiếp tục seed các mục khác thay vì return error ngay
			logger.Warn("SeedWorkStandards: Skipping work standard due to error, continuing with next items", "code", ws.Code)
			errorCount++
			continue
		}

		// Lưu vào map để dùng cho các work standard khác
		codeToIDMap[ws.Code] = workStandardID
		// Lưu UUID -> code nếu có ID cố định
		if ws.ID != uuid.Nil {
			uuidToCodeMap[workStandardID] = ws.Code
		}

		successCount++
		logger.Info("SeedWorkStandards: Created work standard",
			"code", ws.Code,
			"name", ws.Name,
			"id", workStandardID,
			"parent_id", resolvedParentID,
			"granted_by", adminID)
	}

	logger.Info("SeedWorkStandards: Completed seeding work standards",
		"total", len(workStandards),
		"success", successCount,
		"skipped", skipCount,
		"errors", errorCount)

	if errorCount > 0 {
		return fmt.Errorf("seed completed with %d errors out of %d total items", errorCount, len(workStandards))
	}

	return nil
}
