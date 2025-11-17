package spec

import (
	"strconv"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func DefaultVolumes() []*types.VolumeMounts {
	return []*types.VolumeMounts{
		{
			VolumeName: defaultMountName,
			VolumeSize: defaultMountSize,
			MountPath:  defaultMountPath,
		},
	}
}

func DefaultVolumeInput() *types.VolumeCreateInput {
	return &types.VolumeCreateInput{
		Name:      defaultMountName,
		MountPath: defaultMountPath,
		Size:      defaultMountSize,
	}
}

func DefaultNetworkInput() *types.NetworkInput {
	return &types.NetworkInput{
		ContainerPorts: []string{strconv.Itoa(defaultHTTPPort)},
	}
}

func DefaultRegistryIn() *types.ApplicationInput {
	return &types.ApplicationInput{
		RegistryInput: &types.RegistryInput{
			Image: "",
		},
		NetworkInput: DefaultNetworkInput(),
	}
}

func (s *Service) DefaultGitIn() *types.ApplicationInput {
	return &types.ApplicationInput{
		GitInput: &types.GitInput{
			Builder:      enum.BuilderTypeDockerfile,
			BasePath:     "/",
			IsStaticSite: true,
			PublishPath:  "/dist",
			Dockerfile:   "Dockerfile",
		},
		NetworkInput: DefaultNetworkInput(),
	}
}
