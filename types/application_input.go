package types

import (
	"encoding/json"

	"github.com/cloudness-io/cloudness/types/enum"
)

type ApplicationInput struct {
	*GeneralInput
	*RegistryInput
	*GitInput
	*DeployInput
	*NetworkInput
	*VolumeInput
}

type GitInput struct {
	RepoURL      string           `json:"repoURL"`
	Repo         string           `json:"repo"`
	Branch       string           `json:"branch"`
	Commit       string           `json:"commit"`
	CommitMsg    string           `json:"commitMsg"`
	Builder      enum.BuilderType `json:"builder"`
	BasePath     string           `json:"basePath"`
	IsStaticSite bool             `json:"isStaticSite,string"`
	PublishPath  string           `json:"publishPath"`
	Dockerfile   string           `json:"dockerfile"`
	BuildCommand string           `json:"buildCommand,omitempty"`
}

type RegistryInput struct {
	Image string `json:"image"`
	*NetworkInput
}

type DeployInput struct {
	StartCommand            string                 `json:"startCommand,omitempty"`
	SleepApplication        bool                   `json:"sleepApplication,string"`
	MaxReplicas             int64                  `json:"maxReplicas,string"`
	CPU                     int64                  `json:"cpu,string"`
	Memory                  float64                `json:"memory,string"`
	HealthcheckPath         string                 `json:"healthcheckPath,omitempty"`
	HealthcheckTimeout      int                    `json:"healthcheckTimeout,string,omitempty"`
	RestartPolicyType       enum.RestartPolicyType `json:"restartPolicyType"`
	RestartPolicyMaxRetries int                    `json:"restartPolicyMaxRetries"`
}

type NetworkInput struct {
	ContainerPorts []string `json:"containerPorts"`
	*ServiceDomainInput
	*TCPProxiesInput
}

type ServiceDomainInput struct {
	Enabled   bool   `json:"serviceDomainEnabled"`
	FQDN      string `json:"serviceDomainFQDN"`
	Scheme    string `json:"serviceDomainScheme"`
	Subdomain string `json:"serviceDomainSubdomain"`
	Domain    string `json:"serviceDomainDomain"`
	Port      int    `json:"serviceDomainPort,string"`
	Path      string `json:"serviceDomainPath"`
}

type TCPProxiesInput struct {
	Enabled bool `json:"tcpProxiesEnabled"`
	Port    int  `json:"tcpProxiesPort,string"`
}

type VolumeInput struct {
	VolumesRaw string          `json:"volumesJson"`
	Volumes    []*VolumeMounts `json:"volumes"`
}

type GeneralInput struct {
	Icon        string               `json:"icon"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        enum.ApplicationType `json:"type"`
}

func (v *VolumeInput) DecodeVolumes() error {
	if err := json.Unmarshal([]byte(v.VolumesRaw), &v.Volumes); err != nil {
		return err
	}
	return nil
}
