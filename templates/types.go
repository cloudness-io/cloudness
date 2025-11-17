package templates

import (
	"embed"
	"strings"
)

//go:embed *.json
var templates embed.FS

func List() ([][]byte, error) {
	files, err := templates.ReadDir(".")
	if err != nil {
		return nil, err
	}

	dst := make([][]byte, 0)

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".json") {
			content, err := templates.ReadFile(fileName)
			if err != nil {
				return nil, err
			}

			dst = append(dst, content)
		}
	}

	return dst, nil
}
