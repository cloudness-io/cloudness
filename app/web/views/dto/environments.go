package dto

type ListEnvironments struct {
	TargetDiv string
	Envs      []*Environment
}

type Environment struct {
	Name        string
	Description string
	ActionUrl   string
	SettingsUrl string
}
