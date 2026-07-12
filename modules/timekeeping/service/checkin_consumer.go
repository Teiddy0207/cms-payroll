package service

import (
	"cal-salary/core/logger"
	"cal-salary/modules/timekeeping/dto"
	"cal-salary/modules/timekeeping/entity"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// EnsureStreamAndConsumer creates (or updates) the JetStream stream and durable
// pull consumer backing the check-in pipeline. Safe to call on every startup.
func (s *TimekeepingServiceImpl) EnsureStreamAndConsumer(ctx context.Context) error {
	if s.natsClient == nil {
		logger.Warn("EnsureStreamAndConsumer: NATS client is nil, bypassing setup")
		return nil
	}
	_, err := s.natsClient.EnsureStream(ctx, jetstream.StreamConfig{
		Name:       s.streamName,
		Subjects:   []string{s.checkinSubject},
		Retention:  jetstream.LimitsPolicy,
		Storage:    jetstream.FileStorage,
		Replicas:   1,
		MaxAge:     72 * time.Hour,
		Duplicates: 10 * time.Minute,
	})
	if err != nil {
		return err
	}

	_, err = s.natsClient.EnsureConsumer(ctx, s.streamName, jetstream.ConsumerConfig{
		Durable:       s.durableConsumer,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		FilterSubject: s.checkinSubject,
		MaxDeliver:    10,
		AckWait:       30 * time.Second,
	})
	return err
}

// StartCheckinConsumer starts consuming check-in events published to JetStream
// and persists each one to attendance_logs. Message handling runs in the
// background via the underlying NATS client's goroutines.
func (s *TimekeepingServiceImpl) StartCheckinConsumer(ctx context.Context) error {
	if s.natsClient == nil {
		logger.Warn("StartCheckinConsumer: NATS client is nil, bypassing startup")
		return nil
	}
	cons, err := s.natsClient.JS.Consumer(ctx, s.streamName, s.durableConsumer)
	if err != nil {
		return err
	}

	consCtx, err := cons.Consume(func(msg jetstream.Msg) {
		s.handleCheckinMessage(msg)
	})
	if err != nil {
		return err
	}

	s.consumeCtx = consCtx
	return nil
}

// StopCheckinConsumer stops the running consumer, if any (used on graceful shutdown).
func (s *TimekeepingServiceImpl) StopCheckinConsumer() {
	if s.consumeCtx != nil {
		s.consumeCtx.Stop()
	}
}

func (s *TimekeepingServiceImpl) handleCheckinMessage(msg jetstream.Msg) {
	var evt dto.CheckinEvent
	if errUnmarshal := json.Unmarshal(msg.Data(), &evt); errUnmarshal != nil {
		logger.Error("checkin consumer: invalid payload, terminating message", "error", errUnmarshal)
		_ = msg.Term()
		return
	}

	log := &entity.AttendanceLog{
		ID:           uuid.New(),
		EventID:      evt.EventID,
		EmployeeCode: evt.EmployeeCode,
		Timestamp:    evt.Timestamp,
		LocationGPS:  evt.LocationGPS,
		DeviceId:     evt.DeviceId,
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if errCreate := s.repo.CreateAttendanceLog(dbCtx, log); errCreate != nil {
		logger.Error("checkin consumer: failed to persist attendance log, will redeliver", "error", errCreate, "event_id", evt.EventID)
		_ = msg.Nak()
		return
	}

	_ = msg.Ack()
}
