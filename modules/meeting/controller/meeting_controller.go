package controller

import (
	"cal-salary/modules/meeting/dto"
	"cal-salary/modules/meeting/service"
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

func (c *MeetingController) CreateMeeting(ctx echo.Context) error {
	userIDVal := ctx.Get("user_id")
	hostID, ok := userIDVal.(uuid.UUID)
	if !ok {
		if str, isStr := userIDVal.(string); isStr {
			hostID, _ = uuid.Parse(str)
		} else {
			hostID = uuid.Nil
		}
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

func (c *MeetingController) GetMeetings(ctx echo.Context) error {
	userIDVal := ctx.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		if str, isStr := userIDVal.(string); isStr {
			userID, _ = uuid.Parse(str)
		}
	}

	resp, appErr := c.svc.GetMeetings(ctx.Request().Context(), userID)
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

	userIDVal := ctx.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		if str, isStr := userIDVal.(string); isStr {
			userID, _ = uuid.Parse(str)
		}
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
