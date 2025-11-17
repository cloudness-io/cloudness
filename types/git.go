package types

type GitRepo struct {
	Name     string
	FullName string
}

type GitBranch struct {
	Name string
}

type GitCommit struct {
	Sha     string
	Message string
}
