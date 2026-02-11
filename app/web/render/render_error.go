package render

import (
	"context"
	"net/http"

	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/app/web/views"
	"github.com/cloudness-io/cloudness/app/web/views/components/toast"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types/check"

	"github.com/rs/zerolog/log"
)

func renderAalidationErrors(ctx context.Context, w http.ResponseWriter, model any, errs map[string]string) {
	var fields = helpers.ListJsonName(model)
	log.Ctx(ctx).Debug().Any("fields", fields).Any("errs", errs).Msg("Fields")
	for _, field := range fields {
		if _, ok := errs[field]; !ok {
			errs[field] = ""
		}
	}
	HTML(ctx, w, views.ValidationErrors(errs))
}

func ToastErrorMsg(ctx context.Context, w http.ResponseWriter, msg string) {
	toast.ToastError(msg).Render(ctx, w)
}

func ToastError(ctx context.Context, w http.ResponseWriter, err error) {
	var validationErr *check.ValidationErrors
	if errors.As(err, &validationErr) {
		renderAalidationErrors(ctx, w, model{}, validationErr.Errors())
		return
	}
	toast.ToastError(usererror.TranslateErrMsg(ctx, err)).Render(ctx, w)
}

func ToastErrorWithValidation(ctx context.Context, w http.ResponseWriter, model any, err error) {
	var validationErr *check.ValidationErrors
	if errors.As(err, &validationErr) {
		renderAalidationErrors(ctx, w, model, validationErr.Errors())
		return
	}
	toast.ToastError(usererror.TranslateErrMsg(ctx, err)).Render(ctx, w)
}
