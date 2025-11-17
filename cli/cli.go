package cli

import (
	"os"
)

func GetArguments() []string {
	args := os.Args[1:]

	return args
}
