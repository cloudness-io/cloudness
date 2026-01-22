package spec

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

const (
	defaultBasePath           = "/"
	defaultPublishPath        = "/dist"
	defaultDockerfile         = "Dockerfile"
	defaultHTTPPort           = 8080
	defaultTCPPort            = 8000
	defaultMountName          = "Storage"
	defaultMountPath          = "/data"
	defaultMountSize          = 1
	defaultReplicas           = 1
	defaultMemory             = 0.5
	defaultCPU                = 1
	defaultRestartPolicy      = enum.RestartPolicyTypeOnFailure
	defaultRestartMaxRetries  = 5
	defaultHealthCheckPath    = "/"
	defaultHealthCheckTimeout = 300
)

// Application input mapper example
func (s *Service) ToApplication(ctx context.Context, in *types.ApplicationInput, application *types.Application) (*types.Application, error) {
	log.Debug().Any("Input", in).Msg("Application input")
	now := time.Now().UTC().UnixMilli()
	var spec *types.ApplicationSpec

	if application == nil {
		spec = &types.ApplicationSpec{
			Build: &types.BuildConfiguration{},
			Deploy: &types.DeployConfiguration{
				MaxReplicas: defaultReplicas,
				Memory:      defaultMemory,
				CPU:         defaultCPU,
			},
		}

		application = &types.Application{
			UID:     helpers.GenerateUID(),
			Spec:    spec,
			Created: now,
			Updated: now,
		}
	} else {
		spec = application.Spec
	}

	if build, err := s.ToBuildConfigration(ctx, in, application); err != nil {
		return nil, err
	} else {
		spec.Build = build
	}

	if deploy, err := s.ToDeployConfigration(in, application); err != nil {
		return nil, err
	} else {
		spec.Deploy = deploy
	}

	//Networks
	if networks, err := s.ToNetworks(in, application); err != nil {
		return nil, err
	} else {
		spec.Networking = networks
	}

	//Volumes
	if volumes, err := s.ToVolumes(in, application); err != nil {
		return nil, err
	} else {
		spec.Volumes = volumes
	}

	s.ToGeneralSettings(in, application, spec)
	s.UpdateDefaults(in, spec)

	application.Name = spec.Name
	application.Description = spec.Description
	application.Spec = spec
	application.Updated = now
	application.Type = getAppType(spec)

	log.Debug().Any("Spec", spec).Msg("Generated Spec")
	return application, nil
}

func (s *Service) ToGeneralSettings(in *types.ApplicationInput, application *types.Application, spec *types.ApplicationSpec) {
	var name, description, icon string
	//Name
	if in.Name != "" {
		name = in.Name
	} else if application.Name != "" {
		name = application.Name
	} else {
		name = DetectName(spec)
	}

	//Description
	if in.Description != "" {
		description = strings.TrimSpace(in.Description)
	} else if application.Description != "" {
		description = application.Description
	} else {
		description = spec.Description
	}

	//Icon
	if in.Icon != "" {
		icon = in.Icon
	} else if application.Spec.Icon != "" {
		icon = application.Spec.Icon
	} else {
		icon = spec.Icon
	}

	spec.Name = name
	spec.Description = description
	spec.Icon = icon
}

func (s *Service) ToBuildConfigration(ctx context.Context, in *types.ApplicationInput, application *types.Application) (*types.BuildConfiguration, error) {
	switch {
	case in.GitInput != nil:
		git := &types.GitSource{
			RepoURL:  in.RepoURL,
			Branch:   in.Branch,
			Commit:   in.Commit,
			Builder:  in.Builder,
			BasePath: in.BasePath,
		}

		//repo url
		if application.GetGithubAppID() > 0 {
			if ghApp, err := s.ghAppSvc.Find(ctx, application.TenantID, application.ProjectID, *application.GithubAppID); err != nil {
				return nil, err
			} else {
				git.RepoURL = ghApp.GetHttpUrl(in.Repo)
			}
		}

		if in.BasePath == "" {
			git.BasePath = defaultBasePath
		}

		switch in.Builder {
		case enum.BuilderTypeDockerfile:
			if in.Dockerfile != "" {
				git.Dockerfile = in.Dockerfile
			} else {
				git.Dockerfile = defaultDockerfile
			}
		case enum.BuilderTypeNixpacks:
			git.BuildCommand = in.BuildCommand
			// git.StartCommand = in.StartCommand
			if in.IsStaticSite {
				git.IsStaticSite = in.IsStaticSite
				if in.PublishPath != "" {
					git.PublishPath = in.PublishPath
				} else {
					git.PublishPath = defaultPublishPath
				}
			}
		}

		return &types.BuildConfiguration{
			Source: &types.Source{Git: git},
		}, nil
	case in.RegistryInput != nil:
		reg := &types.RegistrySource{
			Image: in.Image,
		}

		return &types.BuildConfiguration{
			Source: &types.Source{Registry: reg},
		}, nil
	}

	// no git input found
	if application.Spec.IsGit() {
		in.GitInput = application.Spec.ToGitInput()
	}
	return application.Spec.Build, nil
}

func (s *Service) ToDeployConfigration(in *types.ApplicationInput, application *types.Application) (*types.DeployConfiguration, error) {
	if in.DeployInput != nil {
		config := &types.DeployConfiguration{
			StartCommand:     in.StartCommand,
			SleepApplication: in.SleepApplication,
			MaxReplicas:      in.MaxReplicas,
			CPU:              in.CPU,
			Memory:           in.Memory,
			HealthcheckPath:  in.HealthcheckPath,
		}

		return config, nil
	}
	return application.Spec.Deploy, nil
}

func (s *Service) ToNetworks(in *types.ApplicationInput, application *types.Application) (*types.NetworkConfiguration, error) {
	if in.NetworkInput != nil {
		ports := make([]int, 0)
		for _, port := range in.ContainerPorts {
			p, err := strconv.Atoi(port)
			if err != nil {
				return nil, err
			}
			ports = append(ports, p)
		}

		var tcpProxies *types.TCPProxiesConfiguration
		if in.TCPProxiesInput != nil {
			tcpProxies = &types.TCPProxiesConfiguration{
				TCPPort: in.TCPProxiesInput.Port,
			}
		}
		return &types.NetworkConfiguration{
			ContainerPorts: ports,
			ServiceDomain:  s.ToServiceDomain(in.ServiceDomainInput),
			TCPProxies:     tcpProxies,
		}, nil
	} else if application.Spec.Networking != nil {
		return application.Spec.Networking, nil
	} else if in.Builder == enum.BuilderTypeStatic {
		return new(types.NetworkConfiguration), nil
	}
	return new(types.NetworkConfiguration), nil
}

func (s *Service) ToServiceDomain(in *types.ServiceDomainInput) *types.ServiceDomainConfiguraiton {
	if in != nil {
		return &types.ServiceDomainConfiguraiton{
			Domain: in.FQDN,
			Port:   in.Port,
			Path:   in.Path,
		}
	}
	return nil
}

func (s *Service) ToVolumes(in *types.ApplicationInput, application *types.Application) ([]*types.VolumeMounts, error) {
	if in.VolumeInput != nil {
		if in.Volumes == nil {
			if err := in.DecodeVolumes(); err != nil {
				return nil, errors.BadRequest("Error decoding volumes")
			}
		}
		return in.Volumes, nil
	} else if application.Spec.Volumes != nil {
		return application.Spec.Volumes, nil
	}
	return []*types.VolumeMounts{}, nil
}

func (s *Service) UpdateDefaults(in *types.ApplicationInput, spec *types.ApplicationSpec) {
	if spec.Deploy.HealthcheckTimeout <= 0 {
		spec.Deploy.HealthcheckTimeout = defaultHealthCheckTimeout
	}

	//update defaults
	if spec.Deploy.MaxReplicas == 0 {
		spec.Deploy.MaxReplicas = defaultReplicas
	}
	if spec.Deploy.Memory == 0 {
		spec.Deploy.Memory = defaultMemory
	}
	if spec.Deploy.CPU == 0 {
		spec.Deploy.CPU = defaultCPU
	}
	if spec.Deploy.RestartPolicyType == "" {
		spec.Deploy.RestartPolicyType = defaultRestartPolicy
	}
	if spec.Deploy.RestartPolicyMaxRetries == 0 {
		spec.Deploy.RestartPolicyMaxRetries = defaultRestartMaxRetries
	}
}
