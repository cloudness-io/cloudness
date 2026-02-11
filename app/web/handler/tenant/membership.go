package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleListMembers(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderMembershipPage(w, r, tenantCtrl)
	}
}

func HandleAddMember(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(tenant.TenantMembershipModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastError(ctx, w, err)
			return
		}

		tenant, _ := request.TenantFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)
		err := tenantCtrl.CreateTenantMembership(ctx, tenant, session, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error adding member")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderMembershipPage(w, r, tenantCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member added successfully")
		}
	}
}

func HandlePatchMember(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(tenant.TenantMembershipModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Msg("invalid request body")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		tenant, _ := request.TenantFrom(ctx)
		err := tenantCtrl.UpdateMembership(ctx, tenant.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating membership")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderMembershipPage(w, r, tenantCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member updated successfully")
		}
	}
}

func HandleDeleteMember(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := &tenant.TenantMembershipModel{
			Email: r.URL.Query().Get("email"),
			Role:  enum.TenantRole(r.URL.Query().Get("role")),
		}

		tenant, _ := request.TenantFrom(ctx)
		err := tenantCtrl.DeleteMembership(ctx, tenant.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting membership")
			render.ToastError(ctx, w, err)
			return
		}

		err = renderMembershipPage(w, r, tenantCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Member deleted successfully")
		}
	}
}

func renderMembershipPage(w http.ResponseWriter, r *http.Request, tenantCtrl *tenant.Controller) error {
	ctx := r.Context()
	session, _ := request.AuthSessionFrom(ctx)
	tenant, _ := request.TenantFrom(ctx)

	memberships, err := tenantCtrl.ListAllMembers(ctx, tenant.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error lising memberships of tenant")
		render.ToastError(ctx, w, err)
		return err
	}

	canEdit := canEdit(ctx, tenantCtrl, tenant)

	render.Page(ctx, w, vtenant.Members(tenant, session, memberships, canEdit))
	return nil
}
