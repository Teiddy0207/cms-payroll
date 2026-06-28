package mapper

import (
	"cal-salary/modules/activity_log/dto"
	"cal-salary/modules/activity_log/entity"
	authDto "cal-salary/modules/auth/dto"
	"encoding/json"
)

func ToActivityLogDTO(log *entity.ActivityLog) *dto.ActivityLogDTO {
	if log == nil {
		return nil
	}

	var userProfile *authDto.LoginUserProfile
	if log.UserProfileID != nil && log.UserProfileFullName != nil {
		userProfile = &authDto.LoginUserProfile{
			ID:       *log.UserProfileID,
			FullName: log.UserProfileFullName,
		}

		if log.UserProfilePositionCode != nil && log.UserProfilePositionName != nil {
			userProfile.Position = &authDto.PositionInfo{
				Code:        *log.UserProfilePositionCode,
				Name:        *log.UserProfilePositionName,
				Description: log.UserProfilePositionDesc,
			}
		}
	}

	return &dto.ActivityLogDTO{
		ID:            log.ID,
		UserID:        log.UserID,
		UserProfileID: log.UserProfileID,
		UserProfile:   userProfile,
		Module:        log.Module,
		Action:        log.Action,
		EntityType:    log.EntityType,
		EntityID:      log.EntityID,
		Description:   log.Description,
		FieldsChanged: log.FieldsChanged,
		OldData:       (*json.RawMessage)(log.OldData),
		NewData:       (*json.RawMessage)(log.NewData),
		HTTPMethod:    log.HTTPMethod,
		Endpoint:      log.Endpoint,
		IPAddress:     log.IPAddress,
		UserAgent:     log.UserAgent,
		Status:        log.Status,
		ErrorMessage:  log.ErrorMessage,
		CreatedAt:     log.CreatedAt,
	}
}

func ToActivityLogPaginationDTO(entity *entity.PaginatedActivityLogEntity) *dto.PaginatedActivityLogDTO {
	if entity == nil {
		return nil
	}

	responses := make([]dto.ActivityLogDTO, len(entity.Items))
	for i, item := range entity.Items {
		mapped := ToActivityLogDTO(&item)
		if mapped != nil {
			responses[i] = *mapped
		}
	}

	totalPages := 0
	if entity.PageSize > 0 {
		totalPages = (entity.TotalItems + entity.PageSize - 1) / entity.PageSize
	}

	return &dto.PaginatedActivityLogDTO{
		Items:      responses,
		TotalItems: entity.TotalItems,
		TotalPages: totalPages,
		PageNumber: entity.PageNumber,
		PageSize:   entity.PageSize,
	}
}
