package source

import (
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
)

const (
	GitPublic = "git-public"
	Github    = "github"
)

var githubSource = &dto.SourceOptions{
	Identifier:  Github,
	Name:        "Private Repository (Github App)",
	Description: "Deploy public and private repository through Github App.",
	Icon:        icons.SourceGithubIcon,
}

var gitSource = &dto.SourceOptions{
	Identifier:  GitPublic,
	Name:        "Public Repository",
	Description: "Deploy any public repository.",
	Icon:        icons.SourceGitIcon,
}

var allSourceOptions = []*dto.SourceOptions{
	gitSource,
	githubSource,
}

var git = []*dto.SourceOptions{
	gitSource,
	githubSource,
}

var database = []*dto.SourceOptions{}

var configurationSourceOptions = []*dto.SourceOptions{
	githubSource,
}

func GetSources(category dto.SourceCategory) []*dto.SourceOptions {
	switch category {
	case dto.SourceCategoryGit:
		return git
	case dto.SourceCategoryDatabase:
		return database
	}

	return []*dto.SourceOptions{}
}

func GetConfirableSources() []*dto.SourceOptions {
	return configurationSourceOptions
}
