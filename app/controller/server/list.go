package server

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ListCertificates(ctx context.Context) ([]*types.Certificate, error) {
	server, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	manageer, err := c.factory.GetServerManager(server)
	if err != nil {
		return nil, err
	}

	return manageer.ListCertificates(ctx, server)
}
