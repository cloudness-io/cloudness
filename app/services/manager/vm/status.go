package vm

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

// dockerInspectStatus holds the fields we extract from `docker inspect`.
type dockerInspectStatus struct {
	ID             string
	Name           string
	State          string // running, exited, paused, restarting, dead, created
	Health         string // healthy, unhealthy, starting, none, ""
	AppUID         int64
	ProjectID      int64
	DeploymentTime int64
}

func (m *VmManager) ListApplicationStatuses(ctx context.Context, server *types.Server) ([]*types.AppStatus, error) {
	filter := fmt.Sprintf("--filter label=%s", managedByFilter())

	// Get container status info using docker inspect format.
	cmd := fmt.Sprintf(
		"docker ps -a --format '{{.ID}}|{{.Names}}|{{.State}}|{{.Label \"%s\"}}|{{.Label \"%s\"}}|{{.Label \"%s\"}}|{{.Label \"%s\"}}' %s",
		labelAppUID,
		labelProjectID,
		labelDeploymentTime,
		labelApp,
		filter,
	)

	out, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list container statuses: %w", err)
	}

	now := time.Now().UTC()
	nowMilli := now.UnixMilli()
	statusByApp := make(map[int64]*types.AppStatus)

	for _, line := range strings.Split(strings.TrimSpace(out.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		info := parseContainerStatusLine(line)
		if info.AppUID == 0 {
			continue
		}

		// Skip recently deployed containers (within 10 seconds).
		if info.DeploymentTime != 0 && info.DeploymentTime > (nowMilli-10_000) {
			continue
		}

		status, reason := evaluateContainerStatus(info)

		candidate := &types.AppStatus{
			Timestamp:      now,
			ApplicationUID: info.AppUID,
			ProjectID:      info.ProjectID,
			InstanceName:   info.Name,
			Status:         status,
			Reason:         reason,
		}

		mergeAppStatus(statusByApp, candidate)
	}

	statuses := make([]*types.AppStatus, 0, len(statusByApp))
	for _, s := range statusByApp {
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func parseContainerStatusLine(line string) dockerInspectStatus {
	// Format: ID|Name|State|AppUID|ProjectID|DeploymentTime|AppName
	parts := strings.SplitN(line, "|", 7)
	if len(parts) < 7 {
		return dockerInspectStatus{}
	}

	appUID, _ := strconv.ParseInt(parts[3], 10, 64)
	projectID, _ := strconv.ParseInt(parts[4], 10, 64)
	deployTime, _ := strconv.ParseInt(parts[5], 10, 64)

	name := parts[1]
	if parts[6] != "" {
		name = parts[6]
	}

	return dockerInspectStatus{
		ID:             parts[0],
		Name:           name,
		State:          parts[2],
		AppUID:         appUID,
		ProjectID:      projectID,
		DeploymentTime: deployTime,
	}
}

func evaluateContainerStatus(info dockerInspectStatus) (enum.ApplicationStatus, string) {
	switch info.State {
	case "running":
		return enum.ApplicationStatusRunning, ""
	case "paused":
		return enum.ApplicationStatusPaused, "container paused"
	case "restarting":
		return enum.ApplicationStatusError, "container restarting"
	case "exited", "dead":
		return enum.ApplicationStatusError, fmt.Sprintf("container %s", info.State)
	case "created":
		return enum.ApplicationStatusPaused, "container created but not started"
	default:
		return enum.ApplicationStatusError, fmt.Sprintf("unknown state: %s", info.State)
	}
}

func mergeAppStatus(statusByApp map[int64]*types.AppStatus, candidate *types.AppStatus) {
	current, found := statusByApp[candidate.ApplicationUID]
	if !found || isHigherPriorityStatus(candidate.Status, current.Status) {
		statusByApp[candidate.ApplicationUID] = candidate
		return
	}

	if current.Status == candidate.Status && current.Reason == "" && candidate.Reason != "" {
		current.Reason = candidate.Reason
	}
}

func isHigherPriorityStatus(candidate, current enum.ApplicationStatus) bool {
	priority := map[enum.ApplicationStatus]int{
		enum.ApplicationStatusError:    3,
		enum.ApplicationStatusPaused:   2,
		enum.ApplicationStatusRunning:  1,
		enum.ApplicationStatusSleeping: 0,
	}
	return priority[candidate] > priority[current]
}
