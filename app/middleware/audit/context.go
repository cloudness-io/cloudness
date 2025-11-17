package audit

import "context"

type key int

const (
	realIPKey key = iota
	requestID
	requestMethod
)

// GetRealIP returns IP address from context.
func GetRealIP(ctx context.Context) string {
	ip, ok := ctx.Value(realIPKey).(string)
	if !ok {
		return ""
	}

	return ip
}

// GetRequestID returns requestID from context.
func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(requestID).(string)
	if !ok {
		return ""
	}

	return id
}

// GetRequestMethod returns http method from context.
func GetRequestMethod(ctx context.Context) string {
	method, ok := ctx.Value(requestMethod).(string)
	if !ok {
		return ""
	}

	return method
}
