package service

import (
	"cal-salary/core/ai"
	"cal-salary/core/errors"
	"cal-salary/core/google"
	"cal-salary/core/logger"
	"cal-salary/core/messaging"
	"cal-salary/core/notification"
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
	UpdateMeeting(ctx context.Context, hostID, meetingID uuid.UUID, req *dto.CreateMeetingRequest) (*dto.MeetingResponse, *errors.AppError)
	CompleteMeeting(ctx context.Context, meetingID uuid.UUID, audioData []byte, filename string) (*entity.MeetingSummary, *errors.AppError)
	SaveMeetingSummary(ctx context.Context, summary *entity.MeetingSummary) *errors.AppError
	GetMeetingSummaryByMeetingID(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingSummary, *errors.AppError)
}

type MeetingService struct {
	repo       repository.MeetingRepositoryInterface
	gcal       *google.CalendarService
	natsClient *messaging.NatsClient
	aiClient   *ai.MeetingAIClient
}

func NewMeetingService(repo repository.MeetingRepositoryInterface, gcal *google.CalendarService, natsClient *messaging.NatsClient) *MeetingService {
	return &MeetingService{
		repo:       repo,
		gcal:       gcal,
		natsClient: natsClient,
		aiClient:   ai.NewMeetingAIClient("http://127.0.0.1:5050"),
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

	resolvedAttendeeIDs := make([]uuid.UUID, len(req.AttendeeIDs))
	attendees := make([]entity.MeetingAttendee, len(req.AttendeeIDs))
	for i, attID := range req.AttendeeIDs {
		resolvedID, _ := s.repo.GetUserIDFromProfileOrUser(ctx, attID)
		resolvedAttendeeIDs[i] = resolvedID
		attendees[i] = entity.MeetingAttendee{
			ID:         uuid.New(),
			MeetingID:  meetingID,
			UserID:     resolvedID,
			RSVPStatus: entity.RSVPStatusPending,
		}
	}
	req.AttendeeIDs = resolvedAttendeeIDs

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

	// Publish real-time notification via NATS if available
	event := notification.NotificationEvent{
		Type:        "meeting.created",
		MeetingID:   meetingID.String(),
		Title:       req.Title,
		AttendeeIDs: req.AttendeeIDs,
		HostID:      hostID,
	}

	if s.natsClient != nil {
		if err := s.natsClient.PublishEvent("meeting.notification", event); err != nil {
			logger.Warn("NATS: failed to publish meeting.created notification", "error", err)
		}
	}

	// Always publish to local SSE GlobalHub for instant NotificationBell updates using resolved UserID
	for _, att := range attendees {
		if att.UserID != uuid.Nil && att.UserID != hostID {
			notification.GlobalHub.Publish(att.UserID, notification.Notification{
				ID:        uuid.New().String(),
				Text:      fmt.Sprintf("Bạn được mời tham gia cuộc họp: %s", req.Title),
				Time:      "Vừa xong",
				Read:      false,
				MeetingID: meetingID.String(),
			})
		}
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

	// Trigger real-time notifications
	meeting, err := s.repo.GetMeetingByID(ctx, meetingID)
	if err == nil && meeting != nil {
		userName, errName := s.repo.GetUserProfileNameByUserID(ctx, userID)
		if errName != nil || userName == "" {
			userName = "Nhân sự"
		}

		actionText := "từ chối"
		if status == entity.RSVPStatusAccepted {
			actionText = "đồng ý"
		}

		text := fmt.Sprintf("%s đã %s tham gia cuộc họp: %s", userName, actionText, meeting.Title)

		event := notification.NotificationEvent{
			Type:      "meeting.rsvp",
			MeetingID: meetingID.String(),
			Title:     meeting.Title,
			HostID:    meeting.HostID,
			UserID:    userID,
			UserName:  userName,
			Status:    string(status),
		}

		natsPublished := false
		if s.natsClient != nil {
			if err := s.natsClient.PublishEvent("meeting.notification", event); err == nil {
				natsPublished = true
			} else {
				logger.Warn("NATS: failed to publish meeting.rsvp notification, falling back to local hub", "error", err)
			}
		}

		if !natsPublished {
			notification.GlobalHub.Publish(meeting.HostID, notification.Notification{
				ID:        uuid.New().String(),
				Text:      text,
				Time:      "Vừa xong",
				Read:      false,
				MeetingID: meetingID.String(),
			})
		}
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

func (s *MeetingService) UpdateMeeting(ctx context.Context, hostID, meetingID uuid.UUID, req *dto.CreateMeetingRequest) (*dto.MeetingResponse, *errors.AppError) {
	existing, err := s.repo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "meeting not found", err)
	}

	if existing.HostID != hostID {
		return nil, errors.NewAppError(errors.ErrForbidden, "only host can update meeting", nil)
	}

	attendees := make([]entity.MeetingAttendee, len(req.AttendeeIDs))
	for i, attID := range req.AttendeeIDs {
		targetID, _ := s.repo.GetUserIDFromProfileOrUser(ctx, attID)
		attendees[i] = entity.MeetingAttendee{
			UserID:     targetID,
			RSVPStatus: entity.RSVPStatusPending,
		}
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.StartTime = req.StartTime
	existing.EndTime = req.EndTime
	existing.Attendees = attendees

	if err := s.repo.UpdateMeeting(ctx, existing); err != nil {
		return nil, errors.NewAppError(errors.ErrUpdateFailed, "failed to update meeting record", err)
	}

	for _, att := range attendees {
		if att.UserID != uuid.Nil && att.UserID != hostID {
			notification.GlobalHub.Publish(att.UserID, notification.Notification{
				ID:        uuid.New().String(),
				Text:      fmt.Sprintf("Cuộc họp '%s' đã được cập nhật!", req.Title),
				Time:      "Vừa xong",
				Read:      false,
				MeetingID: meetingID.String(),
			})
		}
	}

	return toMeetingResponse(existing), nil
}

func (s *MeetingService) SaveMeetingSummary(ctx context.Context, summary *entity.MeetingSummary) *errors.AppError {
	err := s.repo.SaveMeetingSummary(ctx, summary)
	if err != nil {
		return errors.NewAppError(errors.ErrInternalServer, "Failed to save meeting summary", err)
	}
	return nil
}

func (s *MeetingService) GetMeetingSummaryByMeetingID(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingSummary, *errors.AppError) {
	summary, err := s.repo.GetMeetingSummaryByMeetingID(ctx, meetingID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Meeting summary not found", err)
	}
	return summary, nil
}

// CompleteMeeting kết thúc cuộc họp và tự động tóm tắt nội dung bằng AI.
//
// Flow:
//  1. Cập nhật trạng thái meeting → COMPLETED
//  2. Gọi Python AI Service (/analyze-meeting) với file audio
//  3. Lưu MeetingSummary vào DB
//  4. Gửi notification cho host & attendees
func (s *MeetingService) CompleteMeeting(ctx context.Context, meetingID uuid.UUID, audioData []byte, filename string) (*entity.MeetingSummary, *errors.AppError) {
	// Bước 1: Lấy thông tin meeting
	meeting, err := s.repo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Không tìm thấy cuộc họp", err)
	}

	// Bước 2: Cập nhật status → COMPLETED
	if dbErr := s.repo.UpdateMeetingStatus(ctx, meetingID, entity.MeetingStatusCompleted); dbErr != nil {
		return nil, errors.NewAppError(errors.ErrUpdateFailed, "Không thể cập nhật trạng thái cuộc họp", dbErr)
	}

	// Bước 3: Gọi Python AI Service
	aiResult, aiErr := s.aiClient.AnalyzeMeeting(audioData, filename)
	if aiErr != nil {
		logger.Warn("CompleteMeeting: AI service thất bại, lưu summary rỗng", "error", aiErr)
		// Không block flow — vẫn đánh dấu meeting COMPLETED, summary sẽ rỗng
		aiResult = &ai.MeetingAIResponse{
			Summary:         "(Không thể phân tích âm thanh)",
			KeyDecisions:    "",
			ActionItems:     "",
			Sentiment:       "NEUTRAL",
			EfficiencyScore: "0/100",
		}
	}

	// Bước 4: Lưu MeetingSummary vào DB
	summary := &entity.MeetingSummary{
		ID:              uuid.New(),
		MeetingID:       meetingID,
		Summary:         aiResult.Summary,
		KeyDecisions:    aiResult.KeyDecisions,
		ActionItems:     aiResult.ActionItems,
		Sentiment:       aiResult.Sentiment,
		EfficiencyScore: aiResult.EfficiencyScore,
	}

	if dbErr := s.repo.SaveMeetingSummary(ctx, summary); dbErr != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "Không thể lưu tóm tắt cuộc họp", dbErr)
	}

	// Bước 5: Thông báo cho host và tất cả attendees
	notifText := fmt.Sprintf("Cuộc họp '%s' đã kết thúc. Xem tóm tắt ngay!", meeting.Title)

	notification.GlobalHub.Publish(meeting.HostID, notification.Notification{
		ID:        uuid.New().String(),
		Text:      notifText,
		Time:      "Vừa xong",
		Read:      false,
		MeetingID: meetingID.String(),
		Type:      "meeting.completed",
	})

	for _, att := range meeting.Attendees {
		if att.UserID != uuid.Nil && att.UserID != meeting.HostID {
			notification.GlobalHub.Publish(att.UserID, notification.Notification{
				ID:        uuid.New().String(),
				Text:      notifText,
				Time:      "Vừa xong",
				Read:      false,
				MeetingID: meetingID.String(),
				Type:      "meeting.completed",
			})
		}
	}

	return summary, nil
}
