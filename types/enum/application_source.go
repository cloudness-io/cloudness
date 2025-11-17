package enum

type ApplicationSource string

const (
	ApplicationSourceGit      ApplicationSource = "Git"
	ApplicationSourceGithub   ApplicationSource = "GitHub"
	ApplicationSourceTemplate ApplicationSource = "Template"
)
