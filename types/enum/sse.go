package enum

type SSEType string

const (
	SSETypeApplicationUpdated SSEType = "application_updated"
	SSETypeDeploymentUpdated  SSEType = "deployment_updated"
)
