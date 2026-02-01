package breadcrumb

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
)

func PopulateNavItem(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(
				request.WithNavItem(ctx, &dto.NavItem{Title: key}),
			))
		})
	}
}
