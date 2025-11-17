package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cloudness-io/cloudness/pubsub"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Event struct {
	Type enum.SSEType    `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Streamer interface {
	// Publish publishes a new event for the given project.
	Publish(ctx context.Context, projectID int64, eventType enum.SSEType, data any) error

	Stream(ctx context.Context, projectID int64) (<-chan *Event, <-chan error, func(context.Context) error)
}

type streamer struct {
	pubsub pubsub.PubSub
}

func NewStreamer(pubsub pubsub.PubSub) Streamer {
	return &streamer{pubsub: pubsub}
}

func (s *streamer) Publish(ctx context.Context, projectID int64, eventType enum.SSEType, data any) error {
	dataSerialized, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}

	event := Event{
		Type: eventType,
		Data: dataSerialized,
	}
	serializedEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	namespaceOption := pubsub.WithPublishNamespace("sse")
	topic := getProjectTopic(projectID)
	err = s.pubsub.Publish(ctx, topic, serializedEvent, namespaceOption)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (s *streamer) Stream(
	ctx context.Context,
	projectID int64,
) (<-chan *Event, <-chan error, func(context.Context) error) {
	chEvent := make(chan *Event)
	chError := make(chan error)
	g := func(payload []byte) error {
		event := &Event{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			return err
		}

		select {
		case chEvent <- event:
		default:
		}
		return nil
	}

	namespaceOption := pubsub.WithChannelNamespace("sse")
	topic := getProjectTopic(projectID)
	consumer := s.pubsub.Subscribe(ctx, topic, g, namespaceOption)
	cleanupFN := func(_ context.Context) error {
		return consumer.Close()
	}

	return chEvent, chError, cleanupFN
}

func getProjectTopic(projectID int64) string {
	return "projects:" + strconv.Itoa(int(projectID))
}
