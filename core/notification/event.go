package notification

import "github.com/google/uuid"

// NotificationEvent represents the standardized event payload published across NATS
type NotificationEvent struct {
	Type        string      `json:"type"`
	MeetingID   string      `json:"meeting_id,omitempty"`
	Title       string      `json:"title,omitempty"`
	AttendeeIDs []uuid.UUID `json:"attendee_ids,omitempty"`
	HostID      uuid.UUID   `json:"host_id,omitempty"`
	UserID      uuid.UUID   `json:"user_id,omitempty"`
	UserName    string      `json:"user_name,omitempty"`
	Status      string      `json:"status,omitempty"`
}
