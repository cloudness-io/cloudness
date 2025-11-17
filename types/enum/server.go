package enum

type ServerType string

const (
	ServerTypeK8s    ServerType = "k8s"
	ServerTypeDocker ServerType = "docker"
)
