package render

import (
	"context"
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/common"
	"github.com/cloudness-io/cloudness/app/web/views/components/toast"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
	"github.com/cloudness-io/cloudness/app/web/views/layouts"
	"github.com/cloudness-io/cloudness/app/web/views/pages"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/a-h/templ"
	"github.com/rs/zerolog/log"
)

func getBaseLayoutDTO(ctx context.Context) *dto.BaseLayoutOption {
	dto := &dto.BaseLayoutOption{
		Nav: &dto.NavBarOption{},
	}

	if session, ok := request.AuthSessionFrom(ctx); ok {
		dto.Nav.DisplayName = session.Principal.DisplayName
		dto.Nav.Email = session.Principal.Email
	}
	return dto
}

// HTML get at templ component and renders to response writer
func HTML(ctx context.Context, w http.ResponseWriter, c templ.Component) {
	if request.HxIndicatorFrom(ctx) {
		err := c.Render(ctx, w)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
		}
	} else {
		err := layouts.BaseLayoutWithNav(getBaseLayoutDTO(ctx), c).Render(ctx, w)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
		}
	}
}

func HTMLWithBreadCrumb(ctx context.Context, w http.ResponseWriter, c templ.Component) {
	shared.BreadCrumb().Render(ctx, w)
	HTML(ctx, w, c)
}

func Page(ctx context.Context, w http.ResponseWriter, c templ.Component) {
	if request.HxIndicatorFrom(ctx) {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.Header().Set("HX-Retarget", "#main")
		shared.Title().Render(ctx, w)
		shared.BreadCrumb().Render(ctx, w)
		shared.BaseUrl().Render(ctx, w)
		if _, ok := request.ProjectFrom(ctx); ok {
			common.RenderSSEStreamer().Render(ctx, w)
		}
		err := c.Render(ctx, w)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
		}
	} else {
		err := layouts.BaseLayoutWithNav(getBaseLayoutDTO(ctx), c).Render(ctx, w)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
		}
	}
}

func RootWithoutNav(ctx context.Context, w http.ResponseWriter, c templ.Component, url string) {
	if url != "" {
		w.Header().Set("HX-Push-Url", url)
	}
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Retarget", "#root")

	var err error
	if request.HxIndicatorFrom(ctx) {
		err = layouts.LayoutWithoutNav(c).Render(ctx, w)
	} else {
		err = layouts.BaseLayoutWithoutNav(c).Render(ctx, w)
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
	}
}

func RootWithNav(ctx context.Context, w http.ResponseWriter, c templ.Component, url string) {
	if url != "" {
		w.Header().Set("HX-Push-Url", url)
	}
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Retarget", "#root")

	var err error
	if request.HxIndicatorFrom(ctx) {
		err = layouts.LayoutWithNav(getBaseLayoutDTO(ctx), c).Render(ctx, w)
	} else {
		err = layouts.BaseLayoutWithNav(getBaseLayoutDTO(ctx), c).Render(ctx, w)
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to render data from reader")
	}
}

// HTMLWithToastError function renders the template component with the error toast
func HTMLWithToastError(ctx context.Context, w http.ResponseWriter, c templ.Component, toastMsg string) {
	ToastErrorMsg(ctx, w, toastMsg)
	HTML(ctx, w, c)
}

// Redirect makes a client side redirection
func Redirect(w http.ResponseWriter, redirectTo string) {
	w.Header().Set("Hx-Location", redirectTo)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// RedirectWithRefresh redirects the client with a full page refresh, used mostly on error pages
func RedirectWithRefresh(w http.ResponseWriter, redirectTo string) {
	w.Header().Set("Location", redirectTo)
	w.Header().Set("Hx-Redirect", redirectTo)
	w.WriteHeader(http.StatusSeeOther)
}

// RedirectExternal makes a redirection to an external url
func RedirectExternal(w http.ResponseWriter, r *http.Request, redirectTo string) {
	w.Header().Set("Location", redirectTo)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prevURL, _ := request.HxCallerUrlFrom(ctx)
	RootWithoutNav(ctx, w, pages.Error404(prevURL), routes.NotFound)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if request.HxIndicatorFrom(ctx) {
		ToastErrorMsg(ctx, w, "Forbidden")
	} else {
		currURL, _ := request.HxCallerUrlFrom(ctx)
		_, ok := request.PrincipalFrom(ctx)
		if ok {
			RootWithNav(ctx, w, pages.Forbiden(&pages.ForbiddenProps{}), currURL)
		} else {
			RootWithoutNav(ctx, w, pages.Forbiden(&pages.ForbiddenProps{}), currURL)
		}
	}
}

func Error500(w http.ResponseWriter, r *http.Request) {
	RootWithoutNav(r.Context(), w, pages.Error500(), routes.Error)
	w.WriteHeader(http.StatusInternalServerError)
}

func CleanToast(ctx context.Context, w http.ResponseWriter) {
	toast.CleanToast().Render(ctx, w)
}

func ToastSuccess(ctx context.Context, w http.ResponseWriter, msg string) {
	toast.ToastSuccess(msg).Render(ctx, w)
}
