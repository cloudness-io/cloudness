package views

import (
	"fmt"

	"github.com/cloudness-io/cloudness/version"
)

var assetVersion = version.Version.String()

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s?v=%s", name, assetVersion)
}
