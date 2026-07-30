package repository

import (
	"cal-salary/core/database"
	"cal-salary/modules/meeting/entity"
	"context"

	"github.com/google/uuid"
)

type MeetingRepositoryInterface interface {
	CreateMeeting(ctx context.Context, meeting *entity.Meeting) error
	GetMeetingByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error)
	GetMeetingsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Meeting, error)
	UpdateRSVPStatus(ctx context.Context, meetingID, userID uuid.UUID, status entity.RSVPStatus, note string) error
	UpdateMeeting(ctx context.Context, meeting *entity.Meeting) error
	UpdateMeetingStatus(ctx context.Context, meetingID uuid.UUID, status entity.MeetingStatus) error
	GetUserProfileNameByUserID(ctx context.Context, userID uuid.UUID) (string, error)
	GetUserIDFromProfileOrUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	SaveMeetingSummary(ctx context.Context, summary *entity.MeetingSummary) error
	GetMeetingSummaryByMeetingID(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingSummary, error)
}

type MeetingRepository struct {
	db database.IDatabase
}

func NewMeetingRepository(db database.IDatabase) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) GetUserIDFromProfileOrUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var targetUserID uuid.UUID
	query := `
		SELECT user_id FROM user_profiles WHERE (id = $1 OR user_id = $1) AND user_id IS NOT NULL
		UNION ALL
		SELECT id FROM users WHERE id = $1
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &targetUserID, query, id)
	if err != nil || targetUserID == uuid.Nil {
		return id, nil
	}
	return targetUserID, nil
}

func (r *MeetingRepository) CreateMeeting(ctx context.Context, meeting *entity.Meeting) error {
	query := `
		INSERT INTO meetings (id, title, description, host_id, start_time, end_time, google_event_id, room_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	err := r.db.ExecContext(ctx, query,
		meeting.ID, meeting.Title, meeting.Description, meeting.HostID,
		meeting.StartTime, meeting.EndTime, meeting.GoogleEventID, meeting.RoomURL, meeting.Status,
	)
	if err != nil {
		return err
	}

	for _, a := range meeting.Attendees {
		attQuery := `
			INSERT INTO meeting_attendees (id, meeting_id, user_id, rsvp_status, note, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (meeting_id, user_id) DO UPDATE SET rsvp_status = EXCLUDED.rsvp_status
		`
		_ = r.db.ExecContext(ctx, attQuery, a.ID, a.MeetingID, a.UserID, a.RSVPStatus, a.Note)
	}

	return nil
}

func (r *MeetingRepository) GetMeetingByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error) {
	var meeting entity.Meeting
	query := `SELECT id, title, description, host_id, start_time, end_time, google_event_id, room_url, status, created_at, updated_at FROM meetings WHERE id = $1`
	err := r.db.GetContext(ctx, &meeting, query, id)
	if err != nil {
		return nil, err
	}

	var attendees []entity.MeetingAttendee
	attQuery := `SELECT id, meeting_id, user_id, rsvp_status, note, created_at, updated_at FROM meeting_attendees WHERE meeting_id = $1`
	_ = r.db.SelectContext(ctx, &attendees, attQuery, id)
	meeting.Attendees = attendees

	return &meeting, nil
}

func (r *MeetingRepository) GetMeetingsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Meeting, error) {
	var meetings []entity.Meeting
	query := `
		SELECT DISTINCT m.id, m.title, m.description, m.host_id, m.start_time, m.end_time, m.google_event_id, m.room_url, m.status, m.created_at, m.updated_at 
		FROM meetings m
		LEFT JOIN meeting_attendees ma ON ma.meeting_id = m.id
		WHERE m.host_id = $1 OR ma.user_id = $1
		ORDER BY m.start_time DESC
	`
	err := r.db.SelectContext(ctx, &meetings, query, userID)
	if err != nil {
		return []entity.Meeting{}, nil
	}

	for i := range meetings {
		var attendees []entity.MeetingAttendee
		attQuery := `SELECT id, meeting_id, user_id, rsvp_status, note, created_at, updated_at FROM meeting_attendees WHERE meeting_id = $1`
		_ = r.db.SelectContext(ctx, &attendees, attQuery, meetings[i].ID)
		meetings[i].Attendees = attendees
	}

	return meetings, nil
}

func (r *MeetingRepository) UpdateRSVPStatus(ctx context.Context, meetingID, userID uuid.UUID, status entity.RSVPStatus, note string) error {
	query := `
		UPDATE meeting_attendees 
		SET rsvp_status = $1, note = $2, updated_at = NOW() 
		WHERE meeting_id = $3 AND user_id = $4
	`
	return r.db.ExecContext(ctx, query, status, note, meetingID, userID)
}

func (r *MeetingRepository) UpdateMeeting(ctx context.Context, meeting *entity.Meeting) error {
	query := `
		UPDATE meetings 
		SET title = $1, description = $2, start_time = $3, end_time = $4, status = $5, updated_at = NOW() 
		WHERE id = $6
	`
	return r.db.ExecContext(ctx, query, meeting.Title, meeting.Description, meeting.StartTime, meeting.EndTime, meeting.Status, meeting.ID)
}

func (r *MeetingRepository) UpdateMeetingStatus(ctx context.Context, meetingID uuid.UUID, status entity.MeetingStatus) error {
	query := `UPDATE meetings SET status = $1, updated_at = NOW() WHERE id = $2`
	return r.db.ExecContext(ctx, query, status, meetingID)
}

func (r *MeetingRepository) GetUserProfileNameByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	var name string
	query := "SELECT full_name FROM user_profiles WHERE user_id = $1 OR id = $1 LIMIT 1"
	err := r.db.GetContext(ctx, &name, query, userID)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *MeetingRepository) SaveMeetingSummary(ctx context.Context, summary *entity.MeetingSummary) error {
	query := `
		INSERT INTO meeting_summaries (id, meeting_id, summary, key_decisions, action_items, efficiency_score, sentiment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (meeting_id) DO UPDATE 
		SET summary = EXCLUDED.summary, key_decisions = EXCLUDED.key_decisions, action_items = EXCLUDED.action_items,
		    efficiency_score = EXCLUDED.efficiency_score, sentiment = EXCLUDED.sentiment, updated_at = NOW()
	`
	if summary.ID == uuid.Nil {
		summary.ID = uuid.New()
	}
	return r.db.ExecContext(ctx, query, summary.ID, summary.MeetingID, summary.Summary, summary.KeyDecisions, summary.ActionItems, summary.EfficiencyScore, summary.Sentiment)
}

func (r *MeetingRepository) GetMeetingSummaryByMeetingID(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingSummary, error) {
	var summary entity.MeetingSummary
	query := `
		SELECT id, meeting_id, summary, key_decisions, action_items, efficiency_score, sentiment, created_at, updated_at
		FROM meeting_summaries
		WHERE meeting_id = $1
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &summary, query, meetingID)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
