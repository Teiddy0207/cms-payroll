package google

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	srv *calendar.Service
}

func NewCalendarService(ctx context.Context) (*CalendarService, error) {
	apiKey := os.Getenv("GOOGLE_CALENDAR_API_KEY")
	if apiKey == "" {
		// Return service without error; calls will fallback to local calendar if no API key
		return &CalendarService{srv: nil}, nil
	}

	srv, err := calendar.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return &CalendarService{srv: srv}, nil
}

type CreateEventInput struct {
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Attendees   []string
}

func (s *CalendarService) CreateMeetingEvent(ctx context.Context, input CreateEventInput) (string, string, error) {
	if s.srv == nil {
		// Mock event ID and Room URL if Google Calendar API Key is not set in env
		mockEventID := fmt.Sprintf("gcal_%d", time.Now().UnixNano())
		roomSlug := fmt.Sprintf("cms-payroll-room-%d", time.Now().Unix())
		mockRoomURL := fmt.Sprintf("/meetings/%s", roomSlug)
		return mockEventID, mockRoomURL, nil
	}

	attendees := make([]*calendar.EventAttendee, len(input.Attendees))
	for i, email := range input.Attendees {
		attendees[i] = &calendar.EventAttendee{Email: email}
	}

	event := &calendar.Event{
		Summary:     input.Summary,
		Description: input.Description,
		Start: &calendar.EventDateTime{
			DateTime: input.StartTime.Format(time.RFC3339),
			TimeZone: "Asia/Ho_Chi_Minh",
		},
		End: &calendar.EventDateTime{
			DateTime: input.EndTime.Format(time.RFC3339),
			TimeZone: "Asia/Ho_Chi_Minh",
		},
		Attendees: attendees,
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: fmt.Sprintf("req_%d", time.Now().UnixNano()),
			},
		},
	}

	createdEvent, err := s.srv.Events.Insert("primary", event).ConferenceDataVersion(1).Context(ctx).Do()
	if err != nil {
		return "", "", fmt.Errorf("failed to insert calendar event: %w", err)
	}

	meetURL := ""
	if createdEvent.HangoutLink != "" {
		meetURL = createdEvent.HangoutLink
	} else {
		meetURL = fmt.Sprintf("/meetings/room-%s", createdEvent.Id)
	}

	return createdEvent.Id, meetURL, nil
}
