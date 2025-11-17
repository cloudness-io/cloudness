package user

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (c *Controller) Init(ctx context.Context) error {
	//demo user
	if _, err := c.CreateNoAuth(ctx, &CreateInput{
		Email:       demoUserEmail,
		DisplayName: "Demo User",
		Password:    demoUserPassword,
	}); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error creating demo user")
	}

	//agent pipeline user ???

	return nil
}
