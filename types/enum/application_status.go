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
	ApplicationDeploymentStatusNeedsDeployment = "pending"
	ApplicationDeploymentStatusDeploying       = "deploying"
	ApplicationDeploymentStatusFailed          = "failed"
	ApplicationDeploymentStatusSuccess         = "success"
)
