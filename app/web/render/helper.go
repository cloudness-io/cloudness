package render

import (
	"fmt"
	"io"
)

func WriteMesageHeader(w io.Writer) error {
	_, err := io.WriteString(w, "event: message\n")
	if err != nil {
		return fmt.Errorf("failed to send event header: %w", err)
	}

	_, err = io.WriteString(w, "data: ")
	if err != nil {
		return fmt.Errorf("failed to render data header: %w", err)
	}

	return nil
}
