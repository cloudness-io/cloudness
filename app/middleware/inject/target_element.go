package inject

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	TargetQueryParam = "target"
)

func PopulateTargetElemet() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			targetElement := r.URL.Query().Get(TargetQueryParam)

			next.ServeHTTP(w, r.WithContext(
				request.WithTargetElement(ctx, targetElement),
			))
		})
	}
}
