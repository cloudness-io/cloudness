package convert

import (
	"fmt"

	"github.com/cloudness-io/cloudness/version"
)

const (
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
		return "cloudnessio/helper:latest"
	}
	// Use semver String() which includes prerelease (e.g., "0.1.0-alpha.1")
	return fmt.Sprintf("cloudnessio/helper:v%s", ver.String())
}
