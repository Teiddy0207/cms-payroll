package repository

import (
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/activity_log/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Insert adds a new activity log record to the database.
func (r *ActivityLogRepository) Insert(ctx context.Context, log *entity.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (
			id, user_id, user_profile_id, module, action, entity_type, entity_id,
			description, fields_changed, old_data, new_data,
			http_method, endpoint, ip_address, user_agent, status, error_message, created_at
		) VALUES (
			:id, :user_id, :user_profile_id, :module, :action, :entity_type, :entity_id,
			:description, :fields_changed, :old_data, :new_data,
			:http_method, :endpoint, :ip_address, :user_agent, :status, :error_message, :created_at
		)
	`

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	_, err := r.DB.NamedExecContext(ctx, query, log)
	return err
}

// GetActivityLogs retrieves paginated activity logs based on query params.
func (r *ActivityLogRepository) GetActivityLogs(ctx context.Context, p params.QueryParams) (*entity.PaginatedActivityLogEntity, error) {
	offset := (p.PageNumber - 1) * p.PageSize

	baseQuery := `
		FROM activity_logs al
		LEFT JOIN user_profiles up ON al.user_profile_id = up.id
		LEFT JOIN job_descriptions jd ON up.position_id = jd.id
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if p.Search != "" {
		search := "%" + p.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`
		(
			al.module ILIKE $%d OR 
			al.action ILIKE $%d OR
			al.entity_type ILIKE $%d OR
			al.description ILIKE $%d OR
			up.full_name ILIKE $%d
		)
		`, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, search)
		argIndex++
	}

	if module, exists := p.Filters["module"]; exists && module != "" {
		conditions = append(conditions, fmt.Sprintf("al.module = $%d", argIndex))
		args = append(args, module)
		argIndex++
	}
	if action, exists := p.Filters["action"]; exists && action != "" {
		conditions = append(conditions, fmt.Sprintf("al.action = $%d", argIndex))
		args = append(args, action)
		argIndex++
	}
	if entityType, exists := p.Filters["entity_type"]; exists && entityType != "" {
		conditions = append(conditions, fmt.Sprintf("al.entity_type = $%d", argIndex))
		args = append(args, entityType)
		argIndex++
	}
	if status, exists := p.Filters["status"]; exists && status != "" {
		conditions = append(conditions, fmt.Sprintf("al.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalItems int
	err := r.DB.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		logger.Error("ActivityLogRepository:GetActivityLogs - Count", err)
		return nil, err
	}

	dataQuery := `
		SELECT 
			al.id, al.user_id, al.user_profile_id, 
			up.full_name as user_profile_full_name,
			jd.code as user_profile_position_code,
			jd.name as user_profile_position_name,
			jd.main_mission as user_profile_position_description,
			al.module, al.action, al.entity_type, al.entity_id,
			al.description, al.fields_changed, al.old_data, al.new_data,
			al.http_method, al.endpoint, al.ip_address, al.user_agent, al.status, al.error_message, al.created_at
	` + baseQuery + whereClause

	orderByClause := " ORDER BY al.created_at DESC"
	if p.OrderBy != "" {
		allowedColumns := map[string]string{
			"created_at":   "al.created_at",
			"module":       "al.module",
			"action":       "al.action",
			"user_profile": "up.full_name",
			"user_id":      "al.user_id",
		}
		if col, ok := allowedColumns[p.OrderBy]; ok {
			order := "ASC"
			if p.OrderDesc {
				order = "DESC"
			}
			orderByClause = fmt.Sprintf(" ORDER BY %s %s", col, order)
		}
	}

	dataQuery += orderByClause + ` OFFSET $` + fmt.Sprintf("%d", argIndex) + ` LIMIT $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, offset, p.PageSize)

	var items []entity.ActivityLog
	err = r.DB.SelectContext(ctx, &items, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.PaginatedActivityLogEntity{
				Items:      []entity.ActivityLog{},
				TotalItems: 0,
				PageNumber: p.PageNumber,
				PageSize:   p.PageSize,
			}, nil
		}
		logger.Error("ActivityLogRepository:GetActivityLogs - Select", err)
		return nil, err
	}

	return &entity.PaginatedActivityLogEntity{
		Items:      items,
		TotalItems: totalItems,
		PageNumber: p.PageNumber,
		PageSize:   p.PageSize,
	}, nil
}
