package entity

import (
	"time"

	"github.com/google/uuid"
)

type MeetingStatus string

const (
	MeetingStatusScheduled MeetingStatus = "SCHEDULED"
	MeetingStatusCancelled MeetingStatus = "CANCELLED"
	MeetingStatusCompleted MeetingStatus = "COMPLETED"
)

type RSVPStatus string

const (
	RSVPStatusPending  RSVPStatus = "PENDING"
	RSVPStatusAccepted RSVPStatus = "ACCEPTED"
	RSVPStatusDeclined RSVPStatus = "DECLINED"
)

type Meeting struct {
	ID            uuid.UUID         `db:"id" json:"id"`
	Title         string            `db:"title" json:"title"`
	Description   string            `db:"description" json:"description"`
	HostID        uuid.UUID         `db:"host_id" json:"host_id"`
	StartTime     time.Time         `db:"start_time" json:"start_time"`
	EndTime       time.Time         `db:"end_time" json:"end_time"`
	GoogleEventID string            `db:"google_event_id" json:"google_event_id,omitempty"`
	RoomURL       string            `db:"room_url" json:"room_url,omitempty"`
	Status        MeetingStatus     `db:"status" json:"status"`
	CreatedAt     time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time         `db:"updated_at" json:"updated_at"`
	Attendees     []MeetingAttendee `db:"-" json:"attendees,omitempty"`
}

type MeetingAttendee struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	MeetingID  uuid.UUID  `db:"meeting_id" json:"meeting_id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	RSVPStatus RSVPStatus `db:"rsvp_status" json:"rsvp_status"`
	Note       string     `db:"note" json:"note,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

type MeetingSummary struct {
	ID              uuid.UUID `db:"id" json:"id"`
	MeetingID       uuid.UUID `db:"meeting_id" json:"meeting_id"`
	Summary         string    `db:"summary" json:"summary"`
	KeyDecisions    string    `db:"key_decisions" json:"key_decisions"`
	ActionItems     string    `db:"action_items" json:"action_items"`
	EfficiencyScore string    `db:"efficiency_score" json:"efficiency_score,omitempty"`
	Sentiment       string    `db:"sentiment" json:"sentiment,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
