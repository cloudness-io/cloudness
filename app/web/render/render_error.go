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
)

var fields = helpers.ListJsonName(model{})

func ValidationError(ctx context.Context, w http.ResponseWriter, key string, msg string) {
	ValidationErrors(ctx, w, map[string]string{key: msg})
}

func ValidationErrors(ctx context.Context, w http.ResponseWriter, errs map[string]string) {
	for _, field := range fields {
		if _, ok := errs[field]; !ok {
			errs[field] = ""
		}
	}
	HTML(ctx, w, views.ValidationErrors(errs))
}

func ResetValidationErrors(ctx context.Context, w http.ResponseWriter) {
	ValidationErrors(ctx, w, map[string]string{})
}

func ToastErrorMsg(ctx context.Context, w http.ResponseWriter, msg string) {
	toast.ToastError(msg).Render(ctx, w)
}

func ToastError(ctx context.Context, w http.ResponseWriter, err error) {
	var validationErr *check.ValidationErrors
	if errors.As(err, &validationErr) {
		ValidationErrors(ctx, w, validationErr.Errors())
		return
	}
	toast.ToastError(usererror.TranslateErrMsg(ctx, err)).Render(ctx, w)
}
