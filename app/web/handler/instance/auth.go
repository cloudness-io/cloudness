package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleGetAuth(instanceCtrl *instance.Controller, authCtrl *auth.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instance, err := instanceCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance")
			render.ToastError(ctx, w, err)
			return
		}

		auths, err := authCtrl.ListAllMap(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing auths")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vinstance.Auth(instance, auths))
	}
}

func renderAuthPage(ctx context.Context, w http.ResponseWriter, instanceCtrl *instance.Controller, authCtrl *auth.Controller) error {
	instance, err := instanceCtrl.Get(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error getting instance")
		render.ToastError(ctx, w, err)
		return err
	}

	auths, err := authCtrl.ListAllMap(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error listing auths")
		render.ToastError(ctx, w, err)
		return err
	}

	render.Page(ctx, w, vinstance.Auth(instance, auths))
	return nil
}

func HandlePatchPassword(instanceCtrl *instance.Controller, authCtrl *auth.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(auth.ChangePasswordSettings)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		err := authCtrl.UpdatePasswordSettings(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating password settings")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderAuthPage(ctx, w, instanceCtrl, authCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Password authentication setting updated successfully")
		}
	}
}

func HandlePatchDemoUser(instanceCtrl *instance.Controller, authCtrl *auth.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(auth.DemoUserSettings)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		err := authCtrl.UpdateDemoUserAuth(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating demo user settings")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderAuthPage(ctx, w, instanceCtrl, authCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Demo user authentication setting updated successfully")
		}
	}
}

func HandlePatchOauthProvider(instanceCtrl *instance.Controller, authCtrl *auth.Controller, provider enum.AuthProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(types.AuthSetting)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		err := authCtrl.UpdateOauthSettings(ctx, in, provider)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any("Provider", provider).Msg("Error updating oauth settings")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderAuthPage(ctx, w, instanceCtrl, authCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, fmt.Sprintf("%s setting updated successfully", provider))
		}
	}
}
