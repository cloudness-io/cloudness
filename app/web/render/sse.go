package render

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

const (
	pingInterval = 5 * time.Second
	tailMaxTime  = 1 * time.Hour
)

func StreamSSE[T *sse.Event | *types.LogLine | *types.ArtifactLogLine](
	ctx context.Context,
	w http.ResponseWriter,
	chStop <-chan struct{},
	chEvents <-chan T,
	chErr <-chan error,
	renderInit func(ctx context.Context, w io.Writer) error,
	renderFunc func(ctx context.Context, w io.Writer, event T) error,
) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Ctx(ctx).Error().Err(usererror.ErrResponseNotFlushable).Msg("failed to build SSE stream")
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")
	h.Set("Access-Control-Allow-Origin", "*")

	stream := sseStream[T]{
		writer:     w,
		flusher:    flusher,
		renderFunc: renderFunc,
		renderInit: renderInit,
	}

	ctx, ctxCancel := context.WithTimeout(ctx, tailMaxTime)
	defer ctxCancel()

	if err := stream.ping(); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send ping")
		return
	}

	defer func() {
		if err := stream.close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to close stream")
		}
	}()

	pingTimer := time.NewTimer(pingInterval)
	defer pingTimer.Stop()

	isFirstRender := true
	for {
		select {
		case <-ctx.Done():
			log.Ctx(ctx).Debug().Err(ctx.Err()).Msg("stream request context cancelled")
			return

		case <-chStop:
			log.Ctx(ctx).Debug().Msg("stream: app stopeed")
			return

		case err := <-chErr:
			log.Ctx(ctx).Debug().Err(err).Msg("stream: received error from SSE stream")
			return

		case <-pingTimer.C:
			if err := stream.ping(); err != nil {
				log.Ctx(ctx).Err(err).Msg("stream: failed to send ping")
				return
			}

		case event, canProduce := <-chEvents:
			if !canProduce {
				log.Ctx(ctx).Debug().Msg("stream: events channel drained and closed")
				return
			}
			if err := stream.event(ctx, event, isFirstRender); err != nil {
				log.Ctx(ctx).Err(err).Msg("stream: failed to send event")
				return
			}
			isFirstRender = false
		}

		pingTimer.Stop()

		select {
		case <-pingTimer.C:
		default:
		}

		pingTimer.Reset(pingInterval)
	}
}

type sseStream[T any] struct {
	writer     io.Writer
	flusher    http.Flusher
	renderFunc func(ctx context.Context, w io.Writer, event T) error
	renderInit func(ctx context.Context, w io.Writer) error
}

func (r sseStream[T]) event(ctx context.Context, event T, isFirstRender bool) error {

	if isFirstRender {
		if err := r.renderInitInternal(ctx); err != nil {
			return err
		}
	}

	if err := r.renderHeader(); err != nil {
		return err
	}

	if err := r.renderFunc(ctx, r.writer, event); err != nil {
		return fmt.Errorf("failed to render event data: %w", err)
	}

	if err := r.renderFooter(); err != nil {
		return err
	}

	r.flusher.Flush()

	return nil
}

func (r sseStream[T]) renderInitInternal(ctx context.Context) error {
	if err := r.renderHeader(); err != nil {
		return err
	}

	if err := r.renderInit(ctx, r.writer); err != nil {
		return fmt.Errorf("failed to render init data: %w", err)
	}

	if err := r.renderFooter(); err != nil {
		return err
	}

	r.flusher.Flush()
	return nil
}

func (r sseStream[T]) renderHeader() error {
	if _, err := io.WriteString(r.writer, "event: message\n"); err != nil {
		return fmt.Errorf("failed to send event header: %w", err)
	}

	if _, err := io.WriteString(r.writer, "data: "); err != nil {
		return fmt.Errorf("failed to render data header: %w", err)
	}
	return nil
}

func (r sseStream[T]) renderFooter() error {
	if _, err := io.WriteString(r.writer, "\n\n"); err != nil {
		return fmt.Errorf("failed to send event footer: %w", err)
	}
	return nil
}

func (r sseStream[T]) close() error {
	_, err := io.WriteString(r.writer, "event: close\ndata: eof\n\n")
	if err != nil {
		return fmt.Errorf("failed to send EOF: %w", err)
	}
	r.flusher.Flush()
	return nil
}

func (r sseStream[T]) ping() error {
	_, err := io.WriteString(r.writer, "event: ping\n\n")
	if err != nil {
		return fmt.Errorf("failed to send ping: %w", err)
	}
	r.flusher.Flush()
	return nil
}
