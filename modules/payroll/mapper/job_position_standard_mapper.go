package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToJobPositionStandardDTO(item *entity.JobPositionStandardDetail) *dto.JobPositionStandardResponse {
	return &dto.JobPositionStandardResponse{
		JobDescriptionID: item.JobDescriptionID,
		JobStandardID:    item.JobStandardID,
		StandardCode:     item.StandardCode,
		StandardName:     item.StandardName,
		AllowanceValue:   item.AllowanceValue,
	}
}

func ToJobPositionStandardDTOList(list []entity.JobPositionStandardDetail) []dto.JobPositionStandardResponse {
	res := make([]dto.JobPositionStandardResponse, len(list))
	for i, item := range list {
		res[i] = *ToJobPositionStandardDTO(&item)
	}
	return res
}
