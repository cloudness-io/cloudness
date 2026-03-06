package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

// dockerStats is the JSON structure returned by `docker stats --no-stream --format json`.
type dockerStats struct {
	Container string `json:"Container"`
	Name      string `json:"Name"`
	CPUPerc   string `json:"CPUPerc"`  // e.g. "0.50%"
	MemUsage  string `json:"MemUsage"` // e.g. "123.4MiB / 1GiB"
}

func (m *VmManager) ListMetrics(ctx context.Context, server *types.Server) ([]*types.AppMetrics, error) {
	// Get stats for all cloudness-managed containers.
	filter := fmt.Sprintf("--filter label=%s", managedByFilter())
	cmd := fmt.Sprintf(
		"docker ps -q %s | xargs -r docker stats --no-stream --format '{{json .}}'",
		filter,
	)

	out, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker stats: %w", err)
	}

	// Also get container labels to map containers to application UIDs.
	labelCmd := fmt.Sprintf(
		"docker ps --format '{{.ID}}|{{.Label \"%s\"}}' %s",
		labelAppUID, filter,
	)
	labelOut, err := m.runCmd(ctx, server, labelCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get container labels: %w", err)
	}

	// Build container ID -> app UID map.
	appUIDs := make(map[string]int64)
	for _, line := range strings.Split(strings.TrimSpace(labelOut.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			uid, _ := strconv.ParseInt(parts[1], 10, 64)
			appUIDs[parts[0]] = uid
		}
	}

	now := time.Now().UTC()
	var metrics []*types.AppMetrics

	for _, line := range strings.Split(strings.TrimSpace(out.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var stats dockerStats
		if err := json.Unmarshal([]byte(line), &stats); err != nil {
			continue
		}

		appUID := appUIDs[stats.Container]
		if appUID == 0 {
			continue
		}

		cpu := parseCPUPercent(stats.CPUPerc)
		mem := parseMemoryBytes(stats.MemUsage)

		metrics = append(metrics, &types.AppMetrics{
			Timestamp:      now,
			ApplicationUID: appUID,
			InstanceName:   stats.Name,
			CPU:            cpu,
			Memory:         mem,
		})
	}

	return metrics, nil
}

// parseCPUPercent converts "1.50%" to millicores (1500 -> 15 millicores at 1 CPU).
// Docker reports CPU as percentage of total host CPU, we convert to millicores.
func parseCPUPercent(s string) int64 {
	s = strings.TrimSuffix(strings.TrimSpace(s), "%")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	// Convert percentage to millicores (1% = 10 millicores).
	return int64(val * 10)
}

// parseMemoryBytes extracts the used memory from "123.4MiB / 1GiB" format.
func parseMemoryBytes(s string) int64 {
	parts := strings.Split(s, "/")
	if len(parts) == 0 {
		return 0
	}
	used := strings.TrimSpace(parts[0])
	return parseByteSize(used)
}

func parseByteSize(s string) int64 {
	s = strings.TrimSpace(s)
	multiplier := int64(1)
	switch {
	case strings.HasSuffix(s, "GiB"):
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GiB")
	case strings.HasSuffix(s, "MiB"):
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "MiB")
	case strings.HasSuffix(s, "KiB"):
		multiplier = 1024
		s = strings.TrimSuffix(s, "KiB")
	case strings.HasSuffix(s, "B"):
		s = strings.TrimSuffix(s, "B")
	}

	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return int64(val * float64(multiplier))
}
