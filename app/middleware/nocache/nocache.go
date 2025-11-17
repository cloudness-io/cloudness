package nocache

import (
	"net/http"
	"time"
)

// Ported from Chi's middleware, source:
// https://github.com/go-chi/chi/blob/v5.0.12/middleware/nocache.go

// Modified the middleware to retain ETags.

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

// NoCache is same as chi's default NoCache middleware except it doesn't remove etag headers.
func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
