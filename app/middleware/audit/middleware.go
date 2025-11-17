package audit

import (
	"context"
	"net"
	"net/http"
	"strings"
)

var (
	trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

// Middleware process request headers to fill internal info data.
func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if rip := realIP(r); rip != "" {
				ctx = context.WithValue(ctx, realIPKey, rip)
			}

			ctx = context.WithValue(ctx, requestMethod, r.Method)
			ctx = context.WithValue(ctx, requestID, w.Header().Get("X-Request-Id"))

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func realIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
