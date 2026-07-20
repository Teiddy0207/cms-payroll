package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToEmployeeCompetencyDTO(item *entity.EmployeeCompetencyDetail) *dto.EmployeeCompetencyResponse {
	return &dto.EmployeeCompetencyResponse{
		UserProfileID:  item.UserProfileID,
		CompetencyID:   item.CompetencyID,
		CompetencyName: item.CompetencyName,
		CompetencyCode: item.CompetencyCode,
		PointValue:     item.PointValue,
	}
}

func ToEmployeeCompetencyDTOList(list []entity.EmployeeCompetencyDetail) []dto.EmployeeCompetencyResponse {
	res := make([]dto.EmployeeCompetencyResponse, len(list))
	for i, item := range list {
		res[i] = *ToEmployeeCompetencyDTO(&item)
	}
	return res
}
