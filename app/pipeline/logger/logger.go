package logger

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/types"
)

type Logger struct {
	sync.Mutex

	client       client.RunnerClient
	deploymentID int64
	now          time.Time
	lineNum      int

	interval time.Duration
	pending  []*types.LogLine
	history  []*types.LogLine

	close    chan struct{}
	ready    chan struct{}
	replacer *strings.Replacer
}

func New(client client.RunnerClient, deploymentID int64, secret []string) *Logger {
	l := &Logger{
		client:       client,
		deploymentID: deploymentID,
		now:          time.Now().UTC(),
		lineNum:      0,

		interval: time.Second,
		close:    make(chan struct{}),
		ready:    make(chan struct{}),
		replacer: newSecretReplacer(secret),
	}
	go l.start()
	return l
}

func (l *Logger) Write(p []byte) error {
	for _, part := range split(p) {
		line := &types.LogLine{
			Number:    l.lineNum,
			Message:   part,
			Timestamp: int64(time.Since(l.now).Seconds()),
		}

		l.lineNum++

		if l.replacer != nil {
			line.Message = l.replacer.Replace(line.Message)
		}

		l.Lock()
		l.pending = append(l.pending, line)
		l.history = append(l.history, line)
		l.Unlock()
	}

	select {
	case l.ready <- struct{}{}:
	default:
	}
	return nil
}

func (l *Logger) Close() error {
	l.flush()
	close(l.close)
	return l.upload()
}

func (l *Logger) upload() error {
	return l.client.Upload(context.Background(), l.deploymentID, l.history)
}

func (l *Logger) flush() error {
	l.Lock()
	lines := slices.Clone(l.pending)
	l.pending = l.pending[:0]
	l.Unlock()
	if len(lines) == 0 {
		return nil
	}
	return l.client.Stream(context.Background(), l.deploymentID, lines)
}

func (l *Logger) start() error {
	for {
		select {
		case <-l.close:
			return nil
		case <-l.ready:
			select {
			case <-l.close:
				return nil
			case <-time.After(l.interval):
				l.flush()
			}
		}
	}
}

func split(p []byte) []string {
	s := string(p)
	v := []string{s}
	if strings.Contains(strings.TrimSuffix(s, "\n"), "\n") {
		v = strings.SplitAfter(s, "\n")
	}
	return v
}
