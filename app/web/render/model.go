package render

import (
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/types"
)

type model struct {
	types.ApplicationInput
	auth.RegisterUserInput
	tenant.TenantMembershipModel
	project.ProjectMembershipAddModel
	project.ProjectMembershipUpdateModel
	variable.AddVariableInput
	auth.DemoUserSettings
}
