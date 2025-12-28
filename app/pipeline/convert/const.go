package convert

import (
	"fmt"

	"github.com/cloudness-io/cloudness/version"
)

const (
	// images
	busyBoxImage = "busybox:1.37.0"

	// workspace volume
	wsBuildVolumePath = "/cloudness/workspace/build"

	//deployment workspace volume
	wsDeployVolumePath = "/cloudness/workspace/deploy"
)

var (
	//build
	baseImage = getBuilderImage()
)

func getBuilderImage() string {
	ver := version.Version
	if ver.Major == 0 && ver.Minor == 0 && ver.Patch == 0 {
		return "cloudnessio/builder:latest"
	}
	return fmt.Sprintf("cloudnessio/builder:v%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
}
