package project

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/app/web/views/components/vdeployment"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleEvents(appCtx context.Context, projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project, _ := request.ProjectFrom(ctx)

		chEvents, chError, sseCancel, err := projectCtrl.Events(ctx, project.ID)
		if err != nil {
			return
		}

		defer func() {
			if err := sseCancel(ctx); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("HandleEvents: error canceling SSE")
			}
		}()

		render.StreamSSE(ctx, w, appCtx.Done(), chEvents, chError, renderInit, renderFunc)
	}
}

func renderInit(ctx context.Context, w io.Writer) error {
	return nil
}

func renderFunc(ctx context.Context, w io.Writer, event *sse.Event) error {
	switch event.Type {
	case enum.SSETypeDeploymentUpdated:
		d := new(types.Deployment)
		if err := json.Unmarshal(event.Data, d); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error unmarshalling deployment")
			return err
		}
		if err := vdeployment.DeploymentStatusStream(d).Render(ctx, w); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error rendering deployment full status")
			return err
		}
	case enum.SSETypeApplicationDeploymentUpdated:
		d := new(types.Application)
		if err := json.Unmarshal(event.Data, d); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error unmarshalling application")
			return err
		}
		if err := vapplication.AppDeploymentStatus(d, true).Render(ctx, w); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error rendering application status")
			return err
		}
	case enum.SSETypeApplicationStatusUpdated:
		s := new(types.AppStatusEvent)
		if err := json.Unmarshal(event.Data, s); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error unmarshalling application status")
			return err
		}
		if err := vapplication.AppStatus(s.ApplicationUID, s.Status, true).Render(ctx, w); err != nil {
			log.Ctx(ctx).Err(err).Msg("sse render: error rendering application status")
			return err
		}
	}
	return nil
}
