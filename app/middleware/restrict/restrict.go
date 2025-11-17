package restrict

import (
	"context"
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
)

func ToSuperAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !request.IsSuperAdmin(ctx) {
				render.Forbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ToTeamAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !request.IsTeamAdmin(ctx) {
				render.Forbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ToProjectOwner() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !request.IsProjectOwner(ctx) {
				render.Forbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ToProjectContributor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !request.IsProjectContributor(ctx) {
				render.Forbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ModificationToProjectOwner makes sure that only the owner of the project can modify project and environment level
// settings
func ModificationToProjectOwner() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if r.Method == http.MethodPatch && !request.IsProjectOwner(ctx) {
				render.Forbidden(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ToProjectRole() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if !request.IsTeamAdmin(ctx) {
				if !checkMethodByProjectRole(ctx, r.Method) {
					render.Forbidden(w, r)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkMethodByProjectRole(ctx context.Context, method string) bool {
	switch method {
	case http.MethodPost, http.MethodDelete:
		return request.IsProjectOwner(ctx)
	case http.MethodPatch:
		return request.IsProjectContributor(ctx)
	}

	return true
}
