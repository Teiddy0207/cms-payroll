package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHub_RegisterPublishUnregister(t *testing.T) {
	hub := &Hub{
		clients: make(map[uuid.UUID][]chan Notification),
	}

	userID := uuid.New()
	ch := hub.Register(userID)

	assert.NotNil(t, ch)

	notif := Notification{
		ID:        "test-id",
		Text:      "Test notification",
		Time:      "Vừa xong",
		Read:      false,
		MeetingID: "meeting-123",
	}

	hub.Publish(userID, notif)

	received := <-ch
	assert.Equal(t, "test-id", received.ID)
	assert.Equal(t, "Test notification", received.Text)

	hub.Unregister(userID, ch)

	// After unregistering, channel should be closed
	_, ok := <-ch
	assert.False(t, ok)
}
