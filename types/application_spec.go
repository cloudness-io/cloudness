package types

import (
	"github.com/cloudness-io/cloudness/types/enum"
)

type ApplicationSpec struct {
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	Build      *BuildConfiguration   `json:"build,omitempty" yaml:"build,omitempty" mapstructure:"build"`
	Deploy     *DeployConfiguration  `json:"deploy,omitempty" yaml:"deploy,omitempty" mapstructure:"deploy"`
	Networking *NetworkConfiguration `json:"networking,omitempty" yaml:"networking,omitempty" mapstructure:"networking"`

	Volumes []*VolumeMounts `json:"volumes,omitempty" yaml:"volumes,omitempty" mapstructure:"volumes"`
}

type BuildConfiguration struct {
	Source *Source `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source"`
}

type Source struct {
	Git      *GitSource      `json:"git,omitempty" yaml:"git,omitempty" mapstructure:"git"`
	Registry *RegistrySource `json:"registry,omitempty" yaml:"registry,omitempty" mapstructure:"registry"`
}

type GitSource struct {
	RepoURL      string           `json:"repoURL" yaml:"repoURL" mapstructure:"repoURL"`
	Branch       string           `json:"branch" yaml:"branch" mapstructure:"branch"`
	Commit       string           `json:"commit" yaml:"commit" mapstructure:"commit"`
	Builder      enum.BuilderType `json:"builder" yaml:"builder" mapstructure:"builder"`
	BasePath     string           `json:"basePath" yaml:"basePath" mapstructure:"basePath"`
	IsStaticSite bool             `json:"isStaticSite,string" yaml:"isStaticSite" mapstructure:"isStaticSite"`
	PublishPath  string           `json:"publishPath" yaml:"publishPath" mapstructure:"publishPath"`
	Dockerfile   string           `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty" mapstructure:"dockerfile"`
	BuildCommand string           `json:"buildCommand,omitempty" yaml:"buildCommand,omitempty" mapstructure:"buildCommand"`
}

type RegistrySource struct {
	Image     string `json:"image" yaml:"image" mapstructure:"image"`
	IsPrivate bool   `json:"isPrivate,string" yaml:"isPrivate" mapstructure:"isPrivate"`
}

// Deploy
type DeployConfiguration struct {
	StartCommand            string                 `json:"startCommand,omitempty" yaml:"startCommand,omitempty" mapstructure:"startCommand"`
	SleepApplication        bool                   `json:"sleepApplication" yaml:"sleepApplication" mapstructure:"sleepApplication"`
	MaxReplicas             int64                  `json:"maxReplicas" yaml:"maxReplicas" mapstructure:"maxReplicas"`
	CPU                     int64                  `json:"cpu" yaml:"cpu" mapstructure:"cpu"`
	Memory                  float64                `json:"memory" yaml:"memory" mapstructure:"memory"`
	HealthcheckPath         string                 `json:"healthcheckPath,omitempty" yaml:"healthcheckPath" mapstructure:"healthcheckPath"`
	HealthcheckTimeout      int                    `json:"healthcheckTimeout" yaml:"healthcheckTimeout" mapstructure:"healthcheckTimeout"`
	RestartPolicyType       enum.RestartPolicyType `json:"restartPolicyType" yaml:"restartPolicyType" mapstructure:"restartPolicyType"`
	RestartPolicyMaxRetries int                    `json:"restartPolicyMaxRetries" yaml:"restartPolicyMaxRetries" mapstructure:"restartPolicyMaxRetries"`
}

type NetworkConfiguration struct {
	ContainerPorts []int                       `json:"containerPorts" yaml:"containerPorts" mapstructure:"containerPorts"`
	ServiceDomain  *ServiceDomainConfiguraiton `json:"serviceDomain" yaml:"serviceDomain" mapstructure:"serviceDomain"`
	TCPProxies     *TCPProxiesConfiguration    `json:"tcpProxies" yaml:"tcpProxies" mapstructure:"tcpProxies"`
}

type ServiceDomainConfiguraiton struct {
	Port   int    `json:"port" yaml:"port" mapstructure:"port"`
	Domain string `json:"domain" yaml:"domain" mapstructure:"domain"`
	Path   string `json:"path" yaml:"path" mapstructure:"path"`
}

type TCPProxiesConfiguration struct {
	TCPPort int `json:"tcpPort" yaml:"tcpPort" mapstructure:"tcpPort"`
}

// Volumemounts
type VolumeMounts struct {
	VolumeName string `json:"volumeName" yaml:"volumeName" mapstructure:"volumeName"`
	VolumeSize int64  `json:"volumeSize" yaml:"volumeSize" mapstructure:"volumeSize"`
	MountPath  string `json:"mountPath" yaml:"mountPath" mapstructure:"mountPath"`
}

// Helper utilities
func (s *ApplicationSpec) IsGit() bool {
	return s.Build != nil && s.Build.Source != nil && s.Build.Source.Git != nil
}

func (s *ApplicationSpec) IsRegistry() bool {
	return s.Build != nil && s.Build.Source != nil && s.Build.Source.Registry != nil
}

func (s *ApplicationSpec) HasBuildCommand() bool {
	return s.IsGit() && s.Build.Source.Git.BuildCommand != ""
}

func (s *ApplicationSpec) HasStartCommand() bool {
	return s.Deploy.StartCommand != ""
}

func (s *ApplicationSpec) ShouldAddStartCommand() bool {
	return s.Deploy.StartCommand != "" && ((s.IsGit() && s.Build.Source.Git.Builder != enum.BuilderTypeNixpacks) || !s.IsGit())
}

func (b *BuildConfiguration) GetSource() string {
	if b.Source.Git != nil {
		return b.Source.Git.RepoURL
	} else if b.Source.Registry != nil {
		return b.Source.Registry.Image
	}
	return ""
}

func (s *ApplicationSpec) IsServiceDomainEnabled() bool {
	serviceDomain := s.ToServiceDomainInput()
	return serviceDomain != nil && serviceDomain.Enabled
}

func (s *ApplicationSpec) IsTCPProxiesEnabled() bool {
	tcpProxies := s.ToTCPPorxiesInput()
	return tcpProxies != nil && tcpProxies.Enabled
}
