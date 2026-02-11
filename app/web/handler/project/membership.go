package project

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"

	"github.com/rs/zerolog/log"
)

func HandleListMembers(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderMembershipPage(w, r, projectCtrl)
	}
}

func HandleAddMember(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(project.ProjectMembershipAddModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)
		err := projectCtrl.AddMember(ctx, session, tenant.ID, project.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error adding project membership")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		err = renderMembershipPage(w, r, projectCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member added successfully")
		}
	}
}

func HandleUpdateMember(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(project.ProjectMembershipUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		err := projectCtrl.UpdateMember(ctx, tenant.ID, project.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating project membership")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderMembershipPage(w, r, projectCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member updated successfully")
		}
	}
}

func HandleDeleteMember(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := &project.ProjectMembershipRemoveModel{
			Email: r.URL.Query().Get("email"),
		}

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		err := projectCtrl.RemoveMember(ctx, tenant.ID, project.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting project membership")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderMembershipPage(w, r, projectCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member deleted successfully")
		}
	}
}

func HandleListAllMembers(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		members, err := tenantCtrl.ListAllMembers(ctx, tenant.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing members")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vproject.ListMembers(members))
	}
}

func renderMembershipPage(w http.ResponseWriter, r *http.Request, projectCtrl *project.Controller) error {
	ctx := r.Context()
	session, _ := request.AuthSessionFrom(ctx)
	tenant, _ := request.TenantFrom(ctx)
	project, _ := request.ProjectFrom(ctx)

	members, err := projectCtrl.ListMembers(ctx, tenant.ID, project.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error listing memberships")
		render.ToastError(ctx, w, err)
		return err
	}

	render.Page(ctx, w, vproject.Members(session, project, members))
	return nil
}
