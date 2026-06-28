package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToJobPositionCompetencyDTO(item *entity.JobPositionCompetencyDetail) *dto.JobPositionCompetencyResponse {
	return &dto.JobPositionCompetencyResponse{
		JobDescriptionID: item.JobDescriptionID,
		CompetencyID:     item.CompetencyID,
		CompetencyName:   item.CompetencyName,
		CompetencyCode:   item.CompetencyCode,
		RequiredLevel:    item.RequiredLevel,
		Weight:           item.Weight,
		PointValue:       item.PointValue,
	}
}

func ToJobPositionCompetencyDTOList(list []entity.JobPositionCompetencyDetail) []dto.JobPositionCompetencyResponse {
	res := make([]dto.JobPositionCompetencyResponse, len(list))
	for i, item := range list {
		res[i] = *ToJobPositionCompetencyDTO(&item)
	}
	return res
}
