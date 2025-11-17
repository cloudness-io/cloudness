package types

import (
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types/enum"
)

type Deployment struct {
	ID            int64                 `db:"deployment_id"              json:"-"`
	UID           int64                 `db:"deployment_uid"             json:"uid"`
	ApplicationID int64                 `db:"deployment_application_id"  json:"application_id"`
	SpecJson      string                `db:"deployment_spec"            json:"spec"`
	Spec          *ApplicationSpec      `db:"-"                          json:"-"`
	NeedsBuild    bool                  `db:"deployment_needs_build"     json:"needs_build"`
	Triggerer     string                `db:"deployment_triggerer"       json:"triggerer"`
	Title         string                `db:"deployment_title"           json:"title"`
	Action        enum.TriggerAction    `db:"deployment_action"          json:"action"`
	Status        enum.DeploymentStatus `db:"deployment_status"          json:"status"`
	Error         string                `db:"deployment_error"           json:"error"`
	Version       int64                 `db:"deployment_version"         json:"version"`
	Machine       string                `db:"deployment_machine"         json:"machine,omitempty"`
	Started       int64                 `db:"deployment_started"         json:"started"`
	Stopped       int64                 `db:"deployment_stopped"         json:"stopped"`
	Created       int64                 `db:"deployment_created"         json:"created"`
	Updated       int64                 `db:"deployment_updated"         json:"updated"`
}

type DeploymentInfo struct {
	Title       string
	Description string
}

// helpers
func (d *Deployment) GetInfo() *DeploymentInfo {
	info := &DeploymentInfo{}
	splits := strings.Split(d.Title, "\n")
	switch len(splits) {
	case 0, 1:
		info.Title = d.Title
	default:
		info.Title = splits[0]
		info.Description = strings.Join(splits[1:], "\n")
	}
	return info
}

func (d *Deployment) Start() {
	d.Status = enum.DeploymentStatusRunning
	d.Started = time.Now().UTC().UnixMilli()
	d.Stopped = 0
	d.Error = ""
}

func (d *Deployment) Fail(err error) {
	switch d.Status {
	case enum.DeploymentStatusPending, enum.DeploymentStatusRunning:
		d.Status = enum.DeploymentStatusFailed
		d.Error = err.Error()
		d.Stopped = time.Now().UTC().UnixMilli()
		if d.Started == 0 {
			d.Started = d.Stopped
		}
	}
}

func (d *Deployment) Finish() {
	d.Status = enum.DeploymentStatusSuccess
	d.Stopped = time.Now().UTC().UnixMilli()
}

func (d *Deployment) IsDone() bool {
	return isDoneStatus(d.Status)
}

func isDoneStatus(status enum.DeploymentStatus) bool {
	switch status {
	case enum.DeploymentStatusRunning, enum.DeploymentStatusPending:
		return false
	default:
		return true
	}
}
