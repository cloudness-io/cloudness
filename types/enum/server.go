package enum

type ServerType string

const (
	ServerTypeK8s ServerType = "kube"
	ServerTypeVM  ServerType = "vm"
)
