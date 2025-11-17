package pubsub

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ErrClosed = errors.New("pubsub: subscriber is closed")
)

type InMemory struct {
	config   Config
	mutex    sync.Mutex
	registry []*inMemorySubscriber
}

func newInMemory(options ...Option) *InMemory {
	config := Config{
		App:            "app",
		Namespace:      "default",
		HealthInterval: 3 * time.Second,
		SendTimeout:    60 * time.Second,
		ChannelSize:    100,
	}

	for _, f := range options {
		f.Apply(&config)
	}

	return &InMemory{
		config:   config,
		registry: make([]*inMemorySubscriber, 0),
	}
}

func (m *InMemory) Subscribe(
	ctx context.Context,
	topic string,
	handler func(payload []byte) error,
	options ...SubscribeOption,
) Consumer {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config := SubscribeConfig{
		topics:      make([]string, 0),
		app:         m.config.App,
		namespace:   m.config.Namespace,
		sendTimeout: 60,
		channelSize: 100,
	}

	for _, f := range options {
		f.Apply(&config)
	}

	subscriber := &inMemorySubscriber{
		config:  &config,
		handler: handler,
	}

	config.topics = append(config.topics, topic)
	subscriber.topics = subscriber.formatTopics(config.topics...)

	// start subscriber
	go subscriber.start(ctx)

	m.registry = append(m.registry, subscriber)

	return subscriber
}

func (m *InMemory) Publish(ctx context.Context, topic string, payload []byte, options ...PublishOption) error {
	if len(m.registry) == 0 {
		log.Ctx(ctx).Warn().Msg("pubsub: no subscribers")
		return nil
	}

	pubConfig := PublishConfig{
		app:       m.config.App,
		namespace: m.config.Namespace,
	}
	for _, f := range options {
		f.Apply(&pubConfig)
	}

	topic = formatTopic(pubConfig.app, pubConfig.namespace, topic)
	wg := sync.WaitGroup{}
	for _, sub := range m.registry {
		if slices.Contains(sub.topics, topic) && !sub.isClosed() {
			wg.Add(1)
			go func(subscriber *inMemorySubscriber) {
				defer wg.Done()
				//timer based on subscriber config
				t := time.NewTimer(subscriber.config.sendTimeout)
				defer t.Stop()
				select {
				case <-ctx.Done():
					return
				case subscriber.channel <- payload:
					time.Sleep(1 * time.Second) //for some reason in memory subscriber is not always available
					// log.Ctx(ctx).Trace().Msgf("pubsub: message %v send to topic %s", string(payload), topic)
				}

			}(sub)
		}
	}

	wg.Wait()

	return nil
}

type inMemorySubscriber struct {
	config  *SubscribeConfig
	handler func([]byte) error
	channel chan []byte
	once    sync.Once
	mutext  sync.RWMutex
	topics  []string
	closed  bool
}

func (s *inMemorySubscriber) start(ctx context.Context) {
	s.channel = make(chan []byte, s.config.channelSize)
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-s.channel:
			if !ok {
				return
			}
			if err := s.handler(msg); err != nil {
				// TODO: bump
				// err to
				// caller
				log.Ctx(ctx).Err(err).Msgf("in pubsub start: error while running handler for topic")
			}
		}
	}
}

func (s *inMemorySubscriber) Close() error {
	s.mutext.Lock()
	defer s.mutext.Unlock()

	if s.closed {
		return ErrClosed
	}

	s.closed = true
	s.once.Do(func() {
		close(s.channel)
	})
	return nil
}

func (s *inMemorySubscriber) isClosed() bool {
	s.mutext.RLock()
	defer s.mutext.RUnlock()
	return s.closed
}

func (s *inMemorySubscriber) formatTopics(topics ...string) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = formatTopic(s.config.app, s.config.namespace, topic)
	}
	return result
}
