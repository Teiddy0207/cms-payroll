package service

import (
	"cal-salary/core/errors"
	"cal-salary/core/google"
	"cal-salary/core/messaging"
	"cal-salary/modules/meeting/dto"
	"cal-salary/modules/meeting/entity"
	"cal-salary/modules/meeting/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MeetingServiceInterface interface {
	CreateMeeting(ctx context.Context, hostID uuid.UUID, req *dto.CreateMeetingRequest) (*dto.MeetingResponse, *errors.AppError)
	GetMeetings(ctx context.Context, userID uuid.UUID) ([]dto.MeetingResponse, *errors.AppError)
	GetMeetingByID(ctx context.Context, meetingID uuid.UUID) (*dto.MeetingResponse, *errors.AppError)
	UpdateRSVP(ctx context.Context, meetingID, userID uuid.UUID, req *dto.UpdateRSVPRequest) *errors.AppError
}

type MeetingService struct {
	repo       repository.MeetingRepositoryInterface
	gcal       *google.CalendarService
	natsClient *messaging.NatsClient
}

func NewMeetingService(repo repository.MeetingRepositoryInterface, gcal *google.CalendarService, natsClient *messaging.NatsClient) *MeetingService {
	return &MeetingService{
		repo:       repo,
		gcal:       gcal,
		natsClient: natsClient,
	}
}

func (s *MeetingService) CreateMeeting(ctx context.Context, hostID uuid.UUID, req *dto.CreateMeetingRequest) (*dto.MeetingResponse, *errors.AppError) {
	meetingID := uuid.New()
	
	roomURL := fmt.Sprintf("/meetings/room-%s", meetingID.String())

	var googleEventID string
	if req.SyncGoogle && s.gcal != nil {
		gID, meetURL, err := s.gcal.CreateMeetingEvent(ctx, google.CreateEventInput{
			Summary:     req.Title,
			Description: req.Description,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
		})
		if err == nil {
			googleEventID = gID
			if meetURL != "" {
				roomURL = meetURL
			}
		}
	}

	attendees := make([]entity.MeetingAttendee, len(req.AttendeeIDs))
	for i, attID := range req.AttendeeIDs {
		attendees[i] = entity.MeetingAttendee{
			ID:         uuid.New(),
			MeetingID:  meetingID,
			UserID:     attID,
			RSVPStatus: entity.RSVPStatusPending,
		}
	}

	meeting := &entity.Meeting{
		ID:            meetingID,
		Title:         req.Title,
		Description:   req.Description,
		HostID:        hostID,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		GoogleEventID: googleEventID,
		RoomURL:       roomURL,
		Status:        entity.MeetingStatusScheduled,
		Attendees:     attendees,
	}

	if err := s.repo.CreateMeeting(ctx, meeting); err != nil {
		return nil, errors.NewAppError(errors.ErrCreateFailed, "failed to create meeting record", err)
	}

	if s.natsClient != nil {
		_ = s.natsClient.PublishEvent("meeting.created", map[string]interface{}{
			"meeting_id": meetingID.String(),
			"title":      req.Title,
			"host_id":    hostID.String(),
			"start_time": req.StartTime,
		})
	}

	return toMeetingResponse(meeting), nil
}

func (s *MeetingService) GetMeetings(ctx context.Context, userID uuid.UUID) ([]dto.MeetingResponse, *errors.AppError) {
	meetings, err := s.repo.GetMeetingsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrGetFailed, "failed to fetch meetings", err)
	}

	responses := make([]dto.MeetingResponse, len(meetings))
	for i, m := range meetings {
		responses[i] = *toMeetingResponse(&m)
	}
	return responses, nil
}

func (s *MeetingService) GetMeetingByID(ctx context.Context, meetingID uuid.UUID) (*dto.MeetingResponse, *errors.AppError) {
	meeting, err := s.repo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "meeting not found", err)
	}
	return toMeetingResponse(meeting), nil
}

func (s *MeetingService) UpdateRSVP(ctx context.Context, meetingID, userID uuid.UUID, req *dto.UpdateRSVPRequest) *errors.AppError {
	status := entity.RSVPStatus(req.Status)
	err := s.repo.UpdateRSVPStatus(ctx, meetingID, userID, status, req.Note)
	if err != nil {
		return errors.NewAppError(errors.ErrUpdateFailed, "failed to update rsvp status", err)
	}
	return nil
}

func toMeetingResponse(m *entity.Meeting) *dto.MeetingResponse {
	attendeesResp := make([]dto.AttendeeResponse, len(m.Attendees))
	for i, a := range m.Attendees {
		attendeesResp[i] = dto.AttendeeResponse{
			UserID:     a.UserID,
			RSVPStatus: string(a.RSVPStatus),
			Note:       a.Note,
		}
	}

	return &dto.MeetingResponse{
		ID:            m.ID,
		Title:         m.Title,
		Description:   m.Description,
		HostID:        m.HostID,
		StartTime:     m.StartTime,
		EndTime:       m.EndTime,
		GoogleEventID: m.GoogleEventID,
		RoomURL:       m.RoomURL,
		Status:        string(m.Status),
		CreatedAt:     m.CreatedAt,
		Attendees:     attendeesResp,
	}
}
