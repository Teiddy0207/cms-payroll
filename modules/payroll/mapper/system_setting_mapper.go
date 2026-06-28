package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToSystemSettingDTO(item *entity.SystemSetting) *dto.SystemSettingResponse {
	return &dto.SystemSettingResponse{
		Key:         item.Key,
		Value:       item.Value,
		Description: item.Description,
	}
}

func ToSystemSettingDTOList(list []entity.SystemSetting) []dto.SystemSettingResponse {
	res := make([]dto.SystemSettingResponse, len(list))
	for i, item := range list {
		res[i] = *ToSystemSettingDTO(&item)
	}
	return res
}
