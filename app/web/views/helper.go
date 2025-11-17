package views

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/request"
)

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

func HasEnvironment(ctx context.Context) bool {
	_, ok := request.EnvironmentFrom(ctx)

	return ok
}
