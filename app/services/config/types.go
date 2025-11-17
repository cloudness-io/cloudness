package config

type PipelineConfig struct {
	PushRegistryURL string
	PullRegistryURL string
	//Proxy registry
	MirrorRegistryEnabled bool
	MirrorRegistryURL     string
}
