package job

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// State represents state of a background job.
type State string

// State enumeration.
const (
	JobStateScheduled State = "scheduled"
	JobStateRunning   State = "running"
	JobStateFinished  State = "finished"
	JobStateFailed    State = "failed"
	JobStateCanceled  State = "canceled"
)

var jobStates = sortEnum([]State{
	JobStateScheduled,
	JobStateRunning,
	JobStateFinished,
	JobStateFailed,
	JobStateCanceled,
})

func (State) Enum() []interface{} { return toInterfaceSlice(jobStates) }

func (s State) Sanitize() (State, bool) {
	return Sanitize(s, GetAllJobStates)
}

func GetAllJobStates() ([]State, State) {
	return jobStates, ""
}

// Priority represents priority of a background job.
type Priority int

// JobPriority enumeration.
const (
	JobPriorityNormal   Priority = 0
	JobPriorityElevated Priority = 1
)

func (s State) IsCompleted() bool {
	return s == JobStateFinished || s == JobStateFailed || s == JobStateCanceled
}

func sortEnum[T constraints.Ordered](slice []T) []T {
	slices.Sort(slice)
	return slice
}

func toInterfaceSlice[T interface{}](vals []T) []interface{} {
	res := make([]interface{}, len(vals))
	for i := range vals {
		res[i] = vals[i]
	}
	return res
}

func Sanitize[E constraints.Ordered](element E, all func() ([]E, E)) (E, bool) {
	allValues, defValue := all()
	var empty E
	if element == empty && defValue != empty {
		return defValue, true
	}
	idx, exists := slices.BinarySearch(allValues, element)
	if exists {
		return allValues[idx], true
	}
	return defValue, false
}
