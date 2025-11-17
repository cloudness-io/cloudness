package hx

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	HXRequestHeader = "Hx-Request"
)

func PopulateHxIndidcator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			hxRequest := r.Header.Get(HXRequestHeader)

			next.ServeHTTP(w, r.WithContext(
				request.WithHxIndicator(ctx, hxRequest == "true"),
			))
		})
	}
}

func PopulateHxCallerUrl() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			hxPrevUrl := r.Header.Get("Hx-Current-Url")

			next.ServeHTTP(w, r.WithContext(
				request.WithHxCallerUrl(ctx, hxPrevUrl),
			))
		})
	}
}
