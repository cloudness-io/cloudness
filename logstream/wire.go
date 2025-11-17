package logstream

import "github.com/google/wire"

var WireSet = wire.NewSet(
	ProvideLogStream,
)

func ProvideLogStream() LogStream {
	return NewInMemory()
}
