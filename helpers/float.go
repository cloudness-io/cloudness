package helpers

import "fmt"

func ToFloat64String(val float64) string {
	return fmt.Sprintf("%.2f", val)
}
