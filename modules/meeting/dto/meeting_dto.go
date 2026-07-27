package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMeetingRequest struct {
	Title         string      `json:"title" validate:"required"`
	Description   string      `json:"description"`
	StartTime     time.Time   `json:"start_time" validate:"required"`
	EndTime       time.Time   `json:"end_time" validate:"required"`
	AttendeeIDs   []uuid.UUID `json:"attendee_ids" validate:"required,gt=0"`
	SyncGoogle    bool        `json:"sync_google"`
	EnableWebRTC  bool        `json:"enable_webrtc"`
}

type UpdateRSVPRequest struct {
	Status string `json:"status" validate:"required,oneof=ACCEPTED DECLINED"`
	Note   string `json:"note"`
}

type MeetingResponse struct {
	ID            uuid.UUID              `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	HostID        uuid.UUID              `json:"host_id"`
	HostName      string                 `json:"host_name,omitempty"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	GoogleEventID string                 `json:"google_event_id,omitempty"`
	RoomURL       string                 `json:"room_url,omitempty"`
	Status        string                 `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	Attendees     []AttendeeResponse `json:"attendees"`
}

type AttendeeResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"`
	RSVPStatus string    `json:"rsvp_status"`
	Note       string    `json:"note,omitempty"`
}

type WebRTCSignalRequest struct {
	Type       string      `json:"type"`
	MeetingID  string      `json:"meeting_id"`
	FromUserID string      `json:"from_user_id"`
	ToUserID   string      `json:"to_user_id,omitempty"`
	Offer      interface{} `json:"offer,omitempty"`
	Answer     interface{} `json:"answer,omitempty"`
	Candidate  interface{} `json:"candidate,omitempty"`
}
