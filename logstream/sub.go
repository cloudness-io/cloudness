package logstream

import (
	"sync"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type subscriber struct {
	sync.Mutex

	handler chan *types.LogLine
	closec  chan struct{}
	closed  bool
}

func (s *subscriber) publish(line *types.LogLine) {
	defer func() {
		r := recover()
		if r != nil {
			log.Debug().Msg("publish to closed subscriber")
		}
	}()

	s.Lock()
	defer s.Unlock()

	select {
	case <-s.closec:
	case s.handler <- line:
	default:
	}
}

func (s *subscriber) close() {
	s.Lock()
	if !s.closed {
		close(s.closec)
		close(s.handler)
	}
	s.Unlock()
}
