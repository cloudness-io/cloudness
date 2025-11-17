package enum

type VariableType string

const (
	// VariableTypeBuild represents a build variable
	VariableTypeBuild VariableType = "build"

	// VariableTypeRun represents a run variable
	VariableTypeRun VariableType = "run"

	// VariableTypeBuildAndRun represents a build and run variable
	VariableTypeBuildAndRun VariableType = "build_run"
)

var VariableTypesStr = []string{
	string(VariableTypeBuild),
	string(VariableTypeRun),
	string(VariableTypeBuildAndRun),
}
