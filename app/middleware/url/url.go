package url

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
)

func PopulateCurrentUrl() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			domainUrl := r.URL.Scheme + "://" + r.Host
			fullUrl := domainUrl + r.RequestURI

			ctx = request.WithHostDomainUrl(ctx, domainUrl)
			ctx = request.WithCurrentFullUrl(ctx, fullUrl)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
