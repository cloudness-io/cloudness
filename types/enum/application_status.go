package enum

type ApplicationStatus string

const (
	ApplicationStatusRunning  = "running"
	ApplicationStatusPaused   = "paused"
	ApplicationStatusSleeping = "sleeping"
	ApplicationStatusError    = "error"
)

type ApplicationDeploymentStatus string

const (
	ApplicationDeploymentStatusNeedsDeployment = "needs deployment"
	ApplicationDeploymentStatusDeploying       = "deploying"
	ApplicationDeploymentStatusFailed          = "failed"
	ApplicationDeploymentStatusSuccess         = "success"
)
