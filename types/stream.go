package types

type Stream[T any] interface {
	Next() (T, error)
}
