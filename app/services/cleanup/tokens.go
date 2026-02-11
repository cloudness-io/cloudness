package cleanup

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/job"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

const (
	jobTypeTokens        = "cloudness:cleanup:tokens"
	jobCronTokens        = "10 */1 * * * *" // At minute 10 past every hour
	jobMaxDurationTokens = 1 * time.Minute

	// tokenRetentionTime specify the time after which the tokens are key even after expiry
	// can be used for tracing incase
	tokenRetentionTime = 24 * time.Hour // 1 day
)

type tokenCleanupJob struct {
	tokenStore store.TokenStore
}

func newTokenCleanupJob(tokenStore store.TokenStore) *tokenCleanupJob {
	return &tokenCleanupJob{tokenStore}
}

// Pruge expired tokens
func (j *tokenCleanupJob) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {

	expiredBefore := time.Now().Add(-tokenRetentionTime)

	log.Ctx(ctx).Info().Msgf("cleanup: starting to purge tokens expired before %s: %s", expiredBefore.Format(time.RFC3339Nano), expiredBefore.String())

	n, err := j.tokenStore.DeleteExpiredBefore(ctx, expiredBefore, []enum.TokenType{enum.TokenTypeSession})
	if err != nil {
		return "", fmt.Errorf("failed to purge expired tokens: %w", err)
	}

	result := "no expired tokens found"
	if n > 0 {
		result = fmt.Sprintf("deleted %d expired tokens", n)
	}
	log.Ctx(ctx).Info().Msg(result)

	return result, nil
}
