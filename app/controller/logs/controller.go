package logs

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/logstream"
)

type Controller struct {
	logStore  store.LogStore
	logStream logstream.LogStream
}

func NewController(
	logStore store.LogStore,
	logStream logstream.LogStream,
) *Controller {
	return &Controller{
		logStore:  logStore,
		logStream: logStream,
	}
}
