package dto

type SourceCategory string

const (
	SourceCategoryGit      = "git"
	SourceCategoryRegistry = "registry"
	SourceCategoryDatabase = "database"
	SourceCategoryOneclick = "oneclick"
)

type SourceOptions struct {
	Identifier  string
	Name        string
	Description string
	Icon        string
	Disabled    bool
	Category    SourceCategory
}
