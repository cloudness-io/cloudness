package enum

// BuilderType represents the type of the Builder.
type BuilderType string

func GetAllBuilderTypeStr() []string { return builderTypesStr }

const (
	BuilderTypeDockerfile BuilderType = "Dockerfile"

	BuilderTypeNixpacks BuilderType = "Nixpacks"

	BuilderTypeStatic BuilderType = "Static"
)

var builderTypesStr = sortEnum([]string{
	string(BuilderTypeDockerfile),
	string(BuilderTypeNixpacks),
	string(BuilderTypeStatic),
})
