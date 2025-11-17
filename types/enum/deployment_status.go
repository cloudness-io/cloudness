package enum

type DeploymentStatus string

const (
	DeploymentStatusPending   = "pending"
	DeploymentStatusRunning   = "running"
	DeploymentStatusSuccess   = "success"
	DeploymentStatusSkipped   = "skipped"
	DeploymentStatusFailed    = "failed"
	DeploymentStatusCancelled = "cancelled"
)

func (s DeploymentStatus) ToApplicationDeploymentStatus() ApplicationDeploymentStatus {
	switch s {
	case DeploymentStatusRunning:
		return ApplicationDeploymentStatusDeploying
	case DeploymentStatusFailed:
		return ApplicationDeploymentStatusFailed
	case DeploymentStatusSuccess:
		return ApplicationDeploymentStatusSuccess
	}
	return ApplicationDeploymentStatusNeedsDeployment
}
