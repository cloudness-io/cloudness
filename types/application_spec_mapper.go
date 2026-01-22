package types

import (
	"strconv"

	"github.com/cloudness-io/cloudness/helpers"
)

// Application spec to input mapper

func (s *ApplicationSpec) ToInput() *ApplicationInput {
	return &ApplicationInput{
		GeneralInput:  s.ToGeneralInput(),
		GitInput:      s.ToGitInput(),
		RegistryInput: s.ToRegistryInput(),
		DeployInput:   s.ToDeployInput(),
		NetworkInput:  s.ToNetorkInput(),
		VolumeInput:   s.ToVolumeInput(),
	}
}

func (s *ApplicationSpec) ToGeneralInput() *GeneralInput {
	return &GeneralInput{
		Icon:        s.Icon,
		Name:        s.Name,
		Description: s.Description,
	}
}

func (s *ApplicationSpec) ToRegistryInput() *RegistryInput {
	if s.IsRegistry() {
		return &RegistryInput{
			Image: s.Build.Source.Registry.Image,
		}
	}
	return nil
}

func (s *ApplicationSpec) ToGitInput() *GitInput {
	if s.IsGit() {
		gitInput := &GitInput{
			RepoURL:      s.Build.Source.Git.RepoURL,
			Branch:       s.Build.Source.Git.Branch,
			Commit:       s.Build.Source.Git.Commit,
			Builder:      s.Build.Source.Git.Builder,
			BasePath:     s.Build.Source.Git.BasePath,
			IsStaticSite: s.Build.Source.Git.IsStaticSite,
			Dockerfile:   s.Build.Source.Git.Dockerfile,
			BuildCommand: s.Build.Source.Git.BuildCommand,
		}

		owner, repo, err := helpers.SplitGitRepoUrl(s.Build.Source.Git.RepoURL)
		if err == nil {
			gitInput.Repo = owner + "/" + repo
		}
		return gitInput
	}
	return nil
}

func (s *ApplicationSpec) ToDeployInput() *DeployInput {
	return &DeployInput{
		StartCommand:            s.Deploy.StartCommand,
		SleepApplication:        s.Deploy.SleepApplication,
		MaxReplicas:             s.Deploy.MaxReplicas,
		CPU:                     s.Deploy.CPU,
		Memory:                  s.Deploy.Memory,
		HealthcheckPath:         s.Deploy.HealthcheckPath,
		HealthcheckTimeout:      s.Deploy.HealthcheckTimeout,
		RestartPolicyType:       s.Deploy.RestartPolicyType,
		RestartPolicyMaxRetries: s.Deploy.RestartPolicyMaxRetries,
	}
}

func (s *ApplicationSpec) ToNetorkInput() *NetworkInput {
	ports := make([]string, len(s.Networking.ContainerPorts))
	for i, p := range s.Networking.ContainerPorts {
		ports[i] = strconv.Itoa(p)
	}
	return &NetworkInput{
		ContainerPorts:     ports,
		ServiceDomainInput: s.ToServiceDomainInput(),
		TCPProxiesInput:    s.ToTCPPorxiesInput(),
	}
}

func (s *ApplicationSpec) ToServiceDomainInput() *ServiceDomainInput {
	if s.Networking.ServiceDomain == nil {
		return nil
	}

	d := s.Networking.ServiceDomain
	schema, subdomain, domain := helpers.ParseFQDN(d.Domain)
	return &ServiceDomainInput{
		Enabled:   true,
		FQDN:      d.Domain,
		Scheme:    schema,
		Domain:    domain,
		Subdomain: subdomain,
		Port:      d.Port,
	}
}

func (s *ApplicationSpec) ToTCPPorxiesInput() *TCPProxiesInput {
	if s.Networking.TCPProxies == nil {
		return nil
	}
	return &TCPProxiesInput{
		Enabled: true,
		Port:    s.Networking.TCPProxies.TCPPort,
	}
}

func (s *ApplicationSpec) ToVolumeInput() *VolumeInput {
	if len(s.Volumes) > 0 {
		return &VolumeInput{
			Volumes: s.Volumes,
		}
	}
	return nil
}
