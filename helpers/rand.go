package helpers

import (
	"github.com/dchest/uniuri"
)

func RandomLower(length int) string {
	return uniuri.NewLenChars(length, []byte("abcdefghijklmnopqrstuvwxyz"))
}

func Random(length int, characters string) string {
	return uniuri.NewLenChars(length, []byte(characters))
}
