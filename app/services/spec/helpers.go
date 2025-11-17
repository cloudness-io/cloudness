package spec

import (
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func GetImage(app *types.Application, deployment *types.Deployment, config *config.PipelineConfig) (buildImage, pullImage, cacheImage string) {
	switch true {
	case app.Spec.IsGit():
		buildImage = fmt.Sprintf("%s/%d:%d", config.PushRegistryURL, app.UID, deployment.UID)
		pullImage = fmt.Sprintf("%s/%d:%d", config.PullRegistryURL, app.UID, deployment.UID)
		cacheImage = fmt.Sprintf("%s/%d-cache", config.PushRegistryURL, app.UID)
	case app.Spec.IsRegistry():
		pullImage = app.Spec.Build.Source.Registry.Image
	}
	return
}

func IsStateless(spec *types.ApplicationSpec) bool {
	return spec.Deploy.MaxReplicas <= 1
}

func IsStateful(spec *types.ApplicationSpec) bool {
	return len(spec.Volumes) > 0
}

// Pipeline step controller
func NeedsBuild(spec *types.ApplicationSpec) bool {
	return spec.IsGit()
}

func NeedsInit(spec *types.ApplicationSpec) bool {
	return spec.IsGit()
}

func GetBuilder(spec *types.ApplicationSpec) enum.BuilderType {
	if spec.IsGit() {
		return spec.Build.Source.Git.Builder
	}
	return enum.BuilderTypeStatic
}

func GetGitRepoUrl(spec *types.ApplicationSpec) string {
	if spec.IsGit() {
		return spec.Build.Source.Git.RepoURL
	}
	return ""
}

func DetectName(spec *types.ApplicationSpec) string {
	//git spec
	if spec.IsGit() {
		_, repo, err := helpers.SplitGitRepoUrl(spec.Build.Source.Git.RepoURL)
		if err != nil {
			return getName(spec.Build.Source.Git.RepoURL)
		}
		return getName(repo)
	}

	//registry spec
	if spec.IsRegistry() {
		return getNameFromImage(spec.Build.Source.Registry.Image)
	}
	return ""
}

func getName(identifer string) string {
	base := identifer
	splits := strings.Split(base, "/")
	if len(splits) > 0 {
		base = splits[len(splits)-1]
	}

	base = strings.ReplaceAll(base, "-", " ")

	return helpers.CapitalizeSentance(base)
}

func getNameFromImage(identifier string) string {
	// Remove the tag if present
	imageWithoutTag := strings.Split(identifier, ":")[0]
	// Split by "/" and get the last part
	parts := strings.Split(imageWithoutTag, "/")
	return helpers.CapitalizeSentance(parts[len(parts)-1])
}

func getAppType(spec *types.ApplicationSpec) enum.ApplicationType {
	switch true {
	case IsStateful(spec):
		return enum.ApplicationTypeStateful
	default:
		return enum.ApplicationTypeStateless
	}
}

// Web switchers
func GetSourceIcon(app *types.Application) string {
	switch true {
	case app.GetGithubAppID() > 0:
		return icons.SourceGithubIcon
	case app.Spec.IsGit():
		return icons.SourceGitIcon
	case app.Spec.IsRegistry():
		return icons.SourceRegistryIcon
	}
	return icons.SourceGitIcon
}

func GetSourceText(app *types.Application) string {
	switch true {
	case app.GetGithubAppID() > 0:
		return "GitHub Source"
	case app.Spec.IsGit():
		return "Git Source"
	case app.Spec.IsRegistry():
		return "Registry Source"
	}
	return "Source"
}
