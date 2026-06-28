package dto

type UpdateSystemSettingRequest struct {
	Value       string  `json:"value" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type SystemSettingResponse struct {
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	Description *string `json:"description,omitempty"`
}
