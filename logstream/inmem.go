package logstream

import (
	"context"
	"sync"

	"github.com/cloudness-io/cloudness/types"
)

type streamer struct {
	sync.Mutex

	streams map[string]*stream
}

func NewInMemory() LogStream {
	return &streamer{
		streams: make(map[string]*stream),
	}
}

func (s *streamer) Create(ctx context.Context, deploymentId int64) error {
	s.Lock()
	defer s.Unlock()
	s.streams[getStreamId(deploymentId)] = newStream()
	return nil
}

func (s *streamer) Delete(ctx context.Context, deploymentId int64) error {
	s.Lock()
	id := getStreamId(deploymentId)
	stream, ok := s.streams[id]
	if ok {
		delete(s.streams, id)
	}
	s.Unlock()
	return stream.close()
}

func (s *streamer) Write(ctx context.Context, deploymentId int64, line *types.LogLine) error {
	s.Lock()
	defer s.Unlock()
	stream, ok := s.streams[getStreamId(deploymentId)]
	if ok {
		return stream.write(line)
	}
	return nil
}

func (s *streamer) Tail(ctx context.Context, deploymentId int64) (<-chan *types.LogLine, <-chan error) {
	s.Lock()
	stream, ok := s.streams[getStreamId(deploymentId)]
	s.Unlock()
	if !ok {
		return nil, nil
	}
	return stream.subscribe(ctx)
}
