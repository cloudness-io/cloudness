package types

import "fmt"

func namespace(environmentUID int64) string {
	return fmt.Sprintf("ns-%d", environmentUID)
}

func (e *Environment) Namespace() string {
	return namespace(e.UID)
}

func (v *Volume) Namespace() string {
	return namespace(v.EnvironmentUID)
}

func (a *Application) Namespace() string {
	return namespace(a.EnvironmentUID)
}
