package server

import "github.com/cloudness-io/cloudness/http"

// Server is the http server for cloudness.
type Server struct {
	*http.Server
}
