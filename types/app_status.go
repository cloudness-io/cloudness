package types

import (
	"time"

	"github.com/cloudness-io/cloudness/types/enum"
)

type AppStatus struct {
	Timestamp      time.Time              `json:"timestamp"`
	ApplicationUID int64                  `json:"application_uid"`
	ProjectID      int64                  `json:"project_id"`
	InstanceName   string                 `json:"instance_name"`
	Status         enum.ApplicationStatus `json:"status"`
	Reason         string                 `json:"reason"`
}

func (s *AppStatus) ToEvent() *AppStatusEvent {
	return &AppStatusEvent{
		ApplicationUID: s.ApplicationUID,
		ProjectID:      s.ProjectID,
		Status:         s.Status,
		Reason:         s.Reason,
	}
}

type AppStatusEvent struct {
	ApplicationUID int64                  `json:"application_uid"`
	ProjectID      int64                  `json:"project_id"`
	Status         enum.ApplicationStatus `json:"status"`
	Reason         string                 `json:"reason"`
}
