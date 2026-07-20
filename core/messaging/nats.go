package messaging

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsClient struct {
	Conn *nats.Conn
	JS   jetstream.JetStream
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &NatsClient{Conn: nc, JS: js}, nil
}

// Close drains the connection, flushing in-flight publishes/acks before closing.
func (c *NatsClient) Close() {
	if c.Conn != nil {
		_ = c.Conn.Drain()
	}
}

// EnsureStream creates the stream if missing, or updates it to match cfg.
func (c *NatsClient) EnsureStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	return c.JS.CreateOrUpdateStream(ctx, cfg)
}

// EnsureConsumer creates the durable consumer if missing, or updates it to match cfg.
func (c *NatsClient) EnsureConsumer(ctx context.Context, streamName string, cfg jetstream.ConsumerConfig) (jetstream.Consumer, error) {
	return c.JS.CreateOrUpdateConsumer(ctx, streamName, cfg)
}
