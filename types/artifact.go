package types

type Artifact struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type ArtifactLogLine struct {
	ArtifactUID string `json:"artifact_uid"`
	Log         string `json:"log"`
}
