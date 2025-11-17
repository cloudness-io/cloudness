package profiler

import (
	"fmt"
	"strings"
)

type Profiler interface {
	StartProfiling(serviceName, serviceVersion string)
}

type Type string

const (
	TypeGCP Type = "gcp"
)

func ParseType(profilerType string) (Type, bool) {
	switch strings.ToLower(strings.TrimSpace(profilerType)) {
	case string(TypeGCP):
		return TypeGCP, true
	default:
		return "", false
	}
}

func New(profiler Type) (Profiler, error) {
	switch profiler {
	case TypeGCP:
		return &GCPProfiler{}, nil
	default:
		return &NoopProfiler{}, fmt.Errorf("profiler '%s' not supported", profiler)
	}
}
