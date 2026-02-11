package request

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
	"github.com/cloudness-io/cloudness/types"
)

type key int

const (
	authSessionKey key = iota
	instanceKey
	hxRequestKey
	hxPreviousUrl
	hostDomainUrl
	currentFullUrl
	userKey
	requestIDKey
	tenantKey
	tenantMembershipKey
	projectKey
	projectMembershipKey
	environmentKey
	githubAppKey
	applicationKey
	deploymentKey
	volumeKey
	authSettingKey
	navItemsKey
)

// WithRequestID returns a copy of parent in which the request id value is set.
func WithRequestID(parent context.Context, v string) context.Context {
	return context.WithValue(parent, requestIDKey, v)
}

// WithHxIndicator function    returns a copy of parent in which the hx indicator
// value is set
func WithHxIndicator(parent context.Context, v bool) context.Context {
	return context.WithValue(parent, hxRequestKey, v)
}

// HxIndicatorFrom function    returns the value of hx indicator key on the context
func HxIndicatorFrom(ctx context.Context) bool {
	v, ok := ctx.Value(hxRequestKey).(bool)
	return ok && v
}

// WithHxCallerUrl function    returns a copy of parent in which the previous url
// value is set
func WithHxCallerUrl(parent context.Context, v string) context.Context {
	return context.WithValue(parent, hxPreviousUrl, v)
}

// HxCallerUrlFrom function    returns the value of hx prev url key on the context
func HxCallerUrlFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(hxPreviousUrl).(string)
	return v, ok && v != ""
}

// WithHostDomainUrl function    returns a copy of parent in which the host url
// value is set
func WithHostDomainUrl(parent context.Context, v string) context.Context {
	return context.WithValue(parent, hostDomainUrl, v)
}

// HostDomainUrlFrom function    returns the value of host url on the context
func HostDomainUrlFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(hostDomainUrl).(string)
	return v, ok && v != ""
}

// WithCurrentFullUrl function    returns a copy of parent in which the current url
// value is set
func WithCurrentFullUrl(parent context.Context, v string) context.Context {
	return context.WithValue(parent, currentFullUrl, v)
}

// CurrentFullUrlFrom function    returns the value of current request url on the context
func CurrentFullUrlFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(currentFullUrl).(string)
	return v, ok && v != ""
}

// WithInstanceSettings function    returns a copy of parent in which the instance
func WithInstanceSettings(parent context.Context, v *types.Instance) context.Context {
	return context.WithValue(parent, instanceKey, v)
}

// InstanceSettingsFrom function    returns the value of instance
func InstanceSettingsFrom(ctx context.Context) (*types.Instance, bool) {
	v, ok := ctx.Value(instanceKey).(*types.Instance)
	return v, ok && v != nil
}

// WithAuthSession returns a copy of parent in which the principal
// value is set.
func WithAuthSession(parent context.Context, v *auth.Session) context.Context {
	return context.WithValue(parent, authSessionKey, v)
}

// AuthSessionFrom returns the value of the principal key on the
// context.
func AuthSessionFrom(ctx context.Context) (*auth.Session, bool) {
	v, ok := ctx.Value(authSessionKey).(*auth.Session)
	return v, ok && v != nil
}

// PrincipalFrom returns the principal of the authsession.
func PrincipalFrom(ctx context.Context) (*types.Principal, bool) {
	v, ok := AuthSessionFrom(ctx)
	if !ok {
		return nil, false
	}

	return &v.Principal, true
}

// WithTenant function    returns a copy of parent in which the tenant
func WithTenant(parent context.Context, t *types.Tenant) context.Context {
	return context.WithValue(parent, tenantKey, t)
}

// TenantFrom function    returns the value of the tenant
func TenantFrom(ctx context.Context) (*types.Tenant, bool) {
	t, ok := ctx.Value(tenantKey).(*types.Tenant)
	return t, ok && t != nil
}

// WithTenantMembership function    returns a copy of parent in which the tenant membership
func WithTenantMembership(parent context.Context, t *types.TenantMembership) context.Context {
	return context.WithValue(parent, tenantMembershipKey, t)
}

// TenantMembershipFrom function    returns the value of the tenant membership
func TenantMembershipFrom(ctx context.Context) (*types.TenantMembership, bool) {
	t, ok := ctx.Value(tenantMembershipKey).(*types.TenantMembership)
	return t, ok && t != nil
}

// WithProject function    returns a copy of parent in which the project
func WithProject(parent context.Context, p *types.Project) context.Context {
	return context.WithValue(parent, projectKey, p)
}

// ProjectFrom function    returns the value of the project
func ProjectFrom(ctx context.Context) (*types.Project, bool) {
	p, ok := ctx.Value(projectKey).(*types.Project)
	return p, ok && p != nil
}

// WithProjectMembership function    returns a copy of parent in which the project membership
func WithProjectMembership(parent context.Context, p *types.ProjectMembership) context.Context {
	return context.WithValue(parent, projectMembershipKey, p)
}

// ProjectMembershipFrom function    returns the value of the project membership
func ProjectMembershipFrom(ctx context.Context) (*types.ProjectMembership, bool) {
	p, ok := ctx.Value(projectMembershipKey).(*types.ProjectMembership)
	return p, ok && p != nil
}

// WithGithubApp function    returns a copy of parent in which the github app
func WithGithubApp(parent context.Context, ghApp *types.GithubApp) context.Context {
	return context.WithValue(parent, githubAppKey, ghApp)
}

// GithubAppFrom function    returns the value of the github app
func GithubAppFrom(ctx context.Context) (*types.GithubApp, bool) {
	g, ok := ctx.Value(githubAppKey).(*types.GithubApp)
	return g, ok && g != nil
}

// WithEnvironment function    returns a copy of parent in which the environment
func WithEnvironment(parent context.Context, e *types.Environment) context.Context {
	return context.WithValue(parent, environmentKey, e)
}

// EnvironmentFrom function    returns the value of the environment
func EnvironmentFrom(ctx context.Context) (*types.Environment, bool) {
	pe, ok := ctx.Value(environmentKey).(*types.Environment)
	return pe, ok && pe != nil
}

// WithApplication function    returns a copy of parent in which the application
func WithApplication(parent context.Context, c *types.Application) context.Context {
	return context.WithValue(parent, applicationKey, c)
}

// ApplicationFrom function    returns the value of the application
func ApplicationFrom(ctx context.Context) (*types.Application, bool) {
	c, ok := ctx.Value(applicationKey).(*types.Application)
	return c, ok && c != nil
}

// WithDeployment function    returns a copy of parent in which the deployment
func WithDeployment(parent context.Context, d *types.Deployment) context.Context {
	return context.WithValue(parent, deploymentKey, d)
}

// DeploymentFrom function    returns the value of the deployment
func DeploymentFrom(ctx context.Context) (*types.Deployment, bool) {
	d, ok := ctx.Value(deploymentKey).(*types.Deployment)
	return d, ok && d != nil
}

// WithVolume function    returns a copy of parent in which the volume
func WithVolume(parent context.Context, c *types.Volume) context.Context {
	return context.WithValue(parent, volumeKey, c)
}

// VolumeFrom function    returns the value of the volume
func VolumeFrom(ctx context.Context) (*types.Volume, bool) {
	c, ok := ctx.Value(volumeKey).(*types.Volume)
	return c, ok && c != nil
}

// WithAuthSetting function returns a copy of parent in which the authsetting
func WithAuthSetting(parent context.Context, c *types.AuthSetting) context.Context {
	return context.WithValue(parent, authSettingKey, c)
}

// AuthSettingFrom returns the value of the authsetting
func AuthSettingFrom(ctx context.Context) (*types.AuthSetting, bool) {
	c, ok := ctx.Value(authSettingKey).(*types.AuthSetting)
	return c, ok && c != nil
}

// WithNavItem function returns a copy of parent with navigation items
func WithNavItem(parent context.Context, nav *dto.NavItem) context.Context {
	navItems, _ := NavItemsFrom(parent)
	navItems = append(navItems, nav)
	return context.WithValue(parent, navItemsKey, navItems)
}

// NavItensFrom returns the value of the main navigation crumb
func NavItemsFrom(ctx context.Context) ([]*dto.NavItem, bool) {
	c, ok := ctx.Value(navItemsKey).([]*dto.NavItem)
	if c == nil {
		return []*dto.NavItem{}, true
	}
	return c, ok
}

// NavItemsReset resets the navigation items
func NavItemsReset(ctx context.Context) context.Context {
	return context.WithValue(ctx, navItemsKey, []*dto.NavItem{})
}
