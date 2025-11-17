package application

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/app/web/views/components/vlog"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGetLogs(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		render.Page(ctx, w, vapplication.Logs(application))
	}
}

func HandleTailLogs(appCtx context.Context, appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)

		logc, errc, err := appCtrl.TailLog(ctx, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error tailing logs")
			render.ToastError(ctx, w, err)
			return
		}

		render.StreamSSE(ctx, w, appCtx.Done(), logc, errc, renderInit, renderFunc)
	}
}

func renderInit(ctx context.Context, w io.Writer) error {
	if err := vlog.EmptyLogContainerSwap().Render(ctx, w); err != nil {
		log.Ctx(ctx).Err(err).Msg("sse render: error rendering log line")
		return err
	}
	return nil
}

func renderFunc(ctx context.Context, w io.Writer, line *types.ArtifactLogLine) error {
	if line.Log == "" {
		return nil
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)
	if err := vapplication.LogLine(line).Render(ctx, writer); err != nil {
		log.Ctx(ctx).Err(err).Msg("sse render: error rendering log line")
		return err
	}
	// hack so ansii converted html new line will not fragment sse data
	if _, err := w.Write([]byte(strings.ReplaceAll(buf.String(), "\n", " "))); err != nil {
		log.Ctx(ctx).Err(err).Msg("sse render: error writing log line")
		return err
	}

	return nil
}
