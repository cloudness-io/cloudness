package schema

import "strings"

func getPropertyFromPath(path string) string {
	slice := strings.Split(path, "/")
	return slice[len(slice)-1]
}
