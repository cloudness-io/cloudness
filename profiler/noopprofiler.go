package profiler

import "github.com/rs/zerolog/log"

type NoopProfiler struct {
}

func (noopProfiler *NoopProfiler) StartProfiling(serviceName, serviceVersion string) {
	log.Info().Msgf("Not starting profiler for service '%s' with version '%s'", serviceName, serviceVersion)
}
