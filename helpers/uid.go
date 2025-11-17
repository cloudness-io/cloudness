package helpers

import "math/rand/v2"

const uidMin int64 = 10000000
const uidMax int64 = 99999999
const uidRange int64 = uidMax - uidMin

func GenerateUID() int64 {
	return rand.Int64N(uidRange) + uidMin
}
