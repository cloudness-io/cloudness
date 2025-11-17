package enum

type VolumeFormAction string

const (
	VolumeFormActionCreate VolumeFormAction = "create"
	VolumeFormActionUpdate VolumeFormAction = "update"
	VolumeFormActionDelete VolumeFormAction = "delete"
	VolumeFormActionAttach VolumeFormAction = "attach"
)
