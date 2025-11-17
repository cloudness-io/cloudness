package common

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/web/render"
)

func Handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.NotFound(w, r)
	}
}
