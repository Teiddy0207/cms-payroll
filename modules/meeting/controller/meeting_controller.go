package controller

import (
	"cal-salary/core/constants"
	"cal-salary/core/notification"
	"cal-salary/core/utils"
	"cal-salary/modules/meeting/dto"
	"cal-salary/modules/meeting/entity"
	"cal-salary/modules/meeting/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MeetingController struct {
	svc service.MeetingServiceInterface
}

func NewMeetingController(svc service.MeetingServiceInterface) *MeetingController {
	return &MeetingController{svc: svc}
}

func getUserID(ctx echo.Context) uuid.UUID {
	userData := ctx.Get(constants.ContextTokenData)
	if userData == nil {
		return uuid.Nil
	}
	claims, ok := userData.(*utils.TokenClaims)
	if !ok {
		return uuid.Nil
	}
	return claims.UserID
}

func (c *MeetingController) CreateMeeting(ctx echo.Context) error {
	hostID := getUserID(ctx)
	if hostID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req dto.CreateMeetingRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	resp, appErr := c.svc.CreateMeeting(ctx.Request().Context(), hostID, &req)
	if appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (c *MeetingController) UpdateMeeting(ctx echo.Context) error {
	hostID := getUserID(ctx)
	if hostID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	meetingID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid meeting id"})
	}

	var req dto.CreateMeetingRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	resp, appErr := c.svc.UpdateMeeting(ctx.Request().Context(), hostID, meetingID, &req)
	if appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *MeetingController) GetMeetings(ctx echo.Context) error {
	userID := getUserID(ctx)
	if userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	resp, appErr := c.svc.GetMeetings(ctx.Request().Context(), userID)
	if appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *MeetingController) GetMeetingByID(ctx echo.Context) error {
	idStr := ctx.Param("id")
	meetingID, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid meeting id"})
	}

	userID := getUserID(ctx)
	if userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	resp, appErr := c.svc.GetMeetingByID(ctx.Request().Context(), meetingID)
	if appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *MeetingController) UpdateRSVP(ctx echo.Context) error {
	idStr := ctx.Param("id")
	meetingID, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid meeting id"})
	}

	userID := getUserID(ctx)
	if userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req dto.UpdateRSVPRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if appErr := c.svc.UpdateRSVP(ctx.Request().Context(), meetingID, userID, &req); appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "RSVP updated successfully"})
}

func (c *MeetingController) StreamNotifications(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	ctx.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	ctx.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")

	userID := getUserID(ctx)
	if userID == uuid.Nil {
		tokenStr := ctx.QueryParam("token")
		if tokenStr != "" {
			claims, err := utils.ValidateAndParseToken(tokenStr)
			if err == nil && claims != nil {
				userID = claims.UserID
			}
		}
	}

	if userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	ch := notification.GlobalHub.Register(userID)
	defer notification.GlobalHub.Unregister(userID, ch)

	ctx.Response().Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(ctx.Response().Writer, "event: ping\ndata: {}\n\n")
	ctx.Response().Flush()

	for {
		select {
		case notif, ok := <-ch:
			if !ok {
				return nil
			}
			data, errMarshal := json.Marshal(notif)
			if errMarshal != nil {
				continue
			}
			fmt.Fprintf(ctx.Response().Writer, "data: %s\n\n", string(data))
			ctx.Response().Flush()
		case <-ctx.Request().Context().Done():
			return nil
		}
	}
}

func (c *MeetingController) SendWebRTCSignal(ctx echo.Context) error {
	userID := getUserID(ctx)
	if userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req dto.WebRTCSignalRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid signal payload"})
	}

	req.FromUserID = userID.String()

	if req.ToUserID != "" {
		targetUUID, err := uuid.Parse(req.ToUserID)
		if err == nil && targetUUID != uuid.Nil {
			notification.GlobalHub.Publish(targetUUID, notification.Notification{
				ID:        uuid.New().String(),
				Text:      "WebRTC Signaling Message",
				MeetingID: req.MeetingID,
				Type:      "webrtc_signal",
				Payload:   req,
			})
		}
	} else if req.MeetingID != "" {
		meetingUUID, err := uuid.Parse(req.MeetingID)
		if err == nil {
			meeting, appErr := c.svc.GetMeetingByID(ctx.Request().Context(), meetingUUID)
			if appErr == nil && meeting != nil {
				if meeting.HostID != userID {
					notification.GlobalHub.Publish(meeting.HostID, notification.Notification{
						ID:        uuid.New().String(),
						Text:      "WebRTC Signaling Message",
						MeetingID: req.MeetingID,
						Type:      "webrtc_signal",
						Payload:   req,
					})
				}
				for _, att := range meeting.Attendees {
					if att.UserID != userID {
						notification.GlobalHub.Publish(att.UserID, notification.Notification{
							ID:        uuid.New().String(),
							Text:      "WebRTC Signaling Message",
							MeetingID: req.MeetingID,
							Type:      "webrtc_signal",
							Payload:   req,
						})
					}
				}
			}
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (c *MeetingController) SaveSummary(ctx echo.Context) error {
	idStr := ctx.Param("id")
	meetingID, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid meeting id"})
	}

	var req entity.MeetingSummary
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid summary payload"})
	}

	req.MeetingID = meetingID
	appErr := c.svc.SaveMeetingSummary(ctx.Request().Context(), &req)
	if appErr != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": req})
}

func (c *MeetingController) GetSummary(ctx echo.Context) error {
	idStr := ctx.Param("id")
	meetingID, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid meeting id"})
	}

	summary, appErr := c.svc.GetMeetingSummaryByMeetingID(ctx.Request().Context(), meetingID)
	if appErr != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": appErr.Message})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": summary})
}
