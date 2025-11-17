package logstream

import (
	"context"
	"sync"

	"github.com/cloudness-io/cloudness/types"
)

const bufferSize = 5000

type stream struct {
	sync.Mutex

	hist []*types.LogLine
	list map[*subscriber]struct{}
}

func newStream() *stream {
	return &stream{
		list: map[*subscriber]struct{}{},
	}
}

func (s *stream) write(line *types.LogLine) error {
	s.Lock()
	s.hist = append(s.hist, line)
	for l := range s.list {
		if !l.closed {
			l.publish(line)
		}
	}

	if size := len(s.hist); size > bufferSize {
		s.hist = s.hist[size-bufferSize:]
	}
	s.Unlock()
	return nil
}

func (s *stream) subscribe(ctx context.Context) (<-chan *types.LogLine, <-chan error) {
	sub := &subscriber{
		handler: make(chan *types.LogLine, bufferSize),
		closec:  make(chan struct{}),
	}

	err := make(chan error)

	s.Lock()
	for _, line := range s.hist {
		sub.publish(line)
	}
	s.list[sub] = struct{}{}
	s.Unlock()

	go func() {
		defer close(err)
		select {
		case <-sub.closec:
		case <-ctx.Done():
			sub.close()
		}
	}()

	return sub.handler, err
}

func (s *stream) close() error {
	s.Lock()
	defer s.Unlock()
	for sub := range s.list {
		delete(s.list, sub)
		sub.close()
	}
	return nil
}
