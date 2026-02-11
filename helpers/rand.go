package helpers

import (
	"math/rand/v2"

	"github.com/dchest/uniuri"
)

func RandomLower(length int) string {
	return uniuri.NewLenChars(length, []byte("abcdefghijklmnopqrstuvwxyz"))
}

func Random(length int, characters string) string {
	return uniuri.NewLenChars(length, []byte(characters))
}

func RandomNum(min int64, max int64) int64 {
	return rand.Int64N(max-min) + min
}
