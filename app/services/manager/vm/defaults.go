package vm

const (
	traefikDynamicConfigDir = "/data/cloudness/proxy/dynamic"
	traefikACMEStoragePath  = "/data/cloudness/proxy/acme.json"

	defaultRegistryImage = "registry:2.8.3"
	defaultRegistryPort  = 5000

	// Docker labels used for managing containers.
	labelManagedBy      = "cloudness.managed-by"
	labelManagedByValue = "cloudness"
	labelApp            = "cloudness.app"
	labelAppUID         = "cloudness.app-uid"
	labelProjectID      = "cloudness.project-id"
	labelInstance       = "cloudness.instance"
	labelDeploymentTime = "cloudness.deployment-time"
)

// managedByFilter returns a Docker label filter for containers managed by cloudness.
func managedByFilter() string {
	return labelManagedBy + "=" + labelManagedByValue
}

// appFilter returns a Docker label filter for a specific app identifier.
func appFilter(identifier string) string {
	return labelApp + "=" + identifier
}
