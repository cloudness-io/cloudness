package convert

const (
	// images
	busyBoxImage = "busybox:1.37.0"

	//build
	baseImage = "cloudnessio/builder:sha-35ec1d1"

	// workspace volume
	wsBuildVolumePath = "/cloudness/workspace/build"

	//deployment workspace volume
	wsDeployVolumePath = "/cloudness/workspace/deploy"
)
