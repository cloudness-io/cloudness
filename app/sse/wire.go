package sse

import (
	"github.com/cloudness-io/cloudness/pubsub"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideEventStreamer,
)

func ProvideEventStreamer(pubsub pubsub.PubSub) Streamer {
	return NewStreamer(pubsub)
}
