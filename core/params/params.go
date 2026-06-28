package params

import (
	"cal-salary/core/constants"
	"cal-salary/core/utils"

	"github.com/labstack/echo/v4"
)

type QueryParams struct {
	PageNumber int
	PageSize   int
	Search     string
	Filters    map[string]string
	OrderBy    string
	OrderDesc  bool
	KpoType    string
}

func NewQueryParams(c echo.Context) *QueryParams {
	filters := make(map[string]string)

	filterKeys := []string{
		"province_code",
		"user_id",
		"exclude_profile",
		"district_code",
		"status",
		"ids",
		"start_date",
		"parent_id",
		"only_parent",
		"type",
		"office_id",
		"type_compensation",
		"role_ids",
		"branch_ids",
		"user_ids",
		"department_ids",
		"department_id",
		"gender",
		"position_ids",
		"position_id",
		"assigned_date",
		"province_id",
		"timekeeping_sheet_id",
		"apply_at",
		"user_approve",
		"user_created",
		"schedule_id",
		"perspective_id",
		"personnel_id",
		"is_active",
		"is_default",
		"month",
		"all",
		"type_id",
		"quarter",
		"perspective_allocation_list_id",
		"kpo_allocation_list_id",
		"is_read",
		"module",
		"action",
		"entity_type",
		"employee_state",
		"user_profile_id",
		"kpi_quarter_scoring_id",
		"kpi_superior_score_id",
		"kpi_target_id",
		"leave_type",
		"work_date",
		"start_date",
		"end_date",
		"active_quarters",
	}

	for _, key := range filterKeys {
		if value := c.QueryParam(key); value != "" {
			filters[key] = value
		}
	}

	return &QueryParams{
		PageNumber: utils.ToNumberWithDefault(c.QueryParam("page_number"), constants.DefaultPageNumber),
		PageSize:   utils.ToNumberWithDefault(c.QueryParam("page_size"), constants.DefaultPageSize),
		Search:     c.QueryParam("search"),
		Filters:    filters,
		OrderBy:    c.QueryParam("order_by"),
		OrderDesc:  c.QueryParam("order_desc") == "true",
	}
}
