package router

import (
	"net/http"
	"strings"

	"github.com/cloudness-io/cloudness/app/request"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	APIMount = "/api"
)

type Router struct {
	// api APIHandler
	web WebHandler
}

// NewRouter returns a new http.Handler that routes traffic
// to the appropriate handlers.
func NewRouter(
	// api APIHandler,
	web WebHandler,
) *Router {
	return &Router{
		// api: api,
		web: web,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// var err error
	// setup logger for request
	log := log.Logger.With().Logger()
	ctx := log.WithContext(req.Context())
	// add logger to logr interface for usage in 3rd party libs
	ctx = logr.NewContext(ctx, zerologr.New(&log))
	req = req.WithContext(ctx)
	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.
			Str("http.original_url", req.URL.String())
	})

	/*
	 * 1. REST API
	 *
	 * All Rest API calls start with "/api/", and thus can be uniquely identified.
	 */
	// if r.isAPITraffic(req) {
	// 	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
	// 		return c.Str("http.handler", "api")
	// 	})

	// 	// remove matched prefix to simplify API handlers
	// 	if err = stripPrefix(APIMount, req); err != nil {
	// 		log.Err(err).Msgf("Failed striping of prefix for api request.")
	// 		render.InternalError(ctx, w)
	// 		return
	// 	}

	// 	r.api.ServeHTTP(w, req)
	// 	return
	// }

	/*
	 * 2. WEB
	 *
	 * Everything else will be routed to web (or return 404)
	 */
	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("http.handler", "web")
	})

	r.web.ServeHTTP(w, req)
}

// stripPrefix removes the prefix from the request path (or noop if it's not there).
func stripPrefix(prefix string, req *http.Request) error {
	if !strings.HasPrefix(req.URL.Path, prefix) {
		return nil
	}
	return request.ReplacePrefix(req, prefix, "")
}

// isAPITraffic returns true iff the request is identified as part of our rest API.
func (r *Router) isAPITraffic(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, APIMount+"/")
}
