package enum

type ProjectRole string

const (
	ProjectRoleViewer      ProjectRole = "reader"
	ProjectRoleContributor ProjectRole = "contributor"
	ProjectRoleOwner       ProjectRole = "owner"
)

var ProjectRoles = sortEnum([]ProjectRole{
	ProjectRoleViewer,
	ProjectRoleContributor,
	ProjectRoleOwner,
})

var ProjectRolesStr = []string{
	string(ProjectRoleViewer),
	string(ProjectRoleContributor),
	string(ProjectRoleOwner),
}

func ProjectRoleFromString(s string) ProjectRole {
	switch s {
	case string(ProjectRoleViewer):
		return ProjectRoleViewer
	case string(ProjectRoleContributor):
		return ProjectRoleContributor
	case string(ProjectRoleOwner):
		return ProjectRoleOwner
	default:
		return ""
	}
}
