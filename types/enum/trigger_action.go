package enum

type TriggerAction string

const (
	TriggerActionCreate TriggerAction = "Created"
	TriggerActionManual TriggerAction = "Manual"
	TriggerActionHook   TriggerAction = "Hook"
)
