package enum

type SSEType string

const (
	SSETypeApplicationDeploymentUpdated SSEType = "application_deployment_updated"
	SSETypeApplicationStatusUpdated     SSEType = "application_status_updated"
	SSETypeDeploymentUpdated            SSEType = "deployment_updated"
)
