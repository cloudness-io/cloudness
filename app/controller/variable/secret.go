package variable

import (
	"fmt"

	"github.com/cloudness-io/cloudness/helpers"
)

const secretCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (c *Controller) generateSecret(length int) (string, string) {
	if length < 0 {
		length = 15
	}

	return helpers.Random(length, secretCharacters), fmt.Sprintf("${{secret(%d)}}", length)
}
