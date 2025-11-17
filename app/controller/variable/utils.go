package variable

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

var varRefRegex = regexp.MustCompile(`\${{\s*(.*?)\s*}}`)
var varSecretRegex = regexp.MustCompile(`\$\{\{secret\((\d+)\)\}\}`)

func newVariable(envID int64, appID int64, key, value string, varType enum.VariableType) *types.Variable {
	now := time.Now().UTC().UnixMilli()
	return &types.Variable{
		UID:           helpers.GenerateUID(),
		EnvironmentID: envID,
		ApplicationID: appID,
		Key:           key,
		Value:         value,
		TextValue:     value,
		Type:          varType,
		Created:       now,
		Updated:       now,
	}
}

func isSecret(v string) bool {
	return varSecretRegex.MatchString(v)
}

func getRef(value string) []string {
	dst := make([]string, 0)
	matches := varRefRegex.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		dst = append(dst, match[1])
	}
	return dst
}

func replaceWithMap(v string, refMap map[string]string) string {
	return varRefRegex.ReplaceAllStringFunc(v, func(match string) string {
		// Extract the key from inside the match
		key := strings.TrimSpace(varRefRegex.FindStringSubmatch(match)[1])
		if val, ok := refMap[key]; ok {
			return val
		}
		return match
	})
}

func getNameRefKey(v *types.Variable, appID int64) string {
	if v.ApplicationID == appID {
		return v.Key
	}
	return fmt.Sprintf("%s.%s", v.ApplicationName, v.Key)
}

func getIDRefKey(v *types.Variable) string {
	return fmt.Sprintf("${{%d}}", v.UID)
}

func getIDToNameMap(vars []*types.Variable, appID int64) map[string]string {
	varMap := make(map[string]string)
	for _, v := range vars {
		varMap[strconv.FormatInt(v.UID, 10)] = fmt.Sprintf("${{%s}}", getNameRefKey(v, appID)) // adding map[123123] = `${{postgres.URL}}`
	}
	return varMap
}

func getNameAndIDRefMap(vars []*types.Variable, appId int64) (map[string]string, map[string]string) {
	nameToIDMap := make(map[string]string)
	idToValueMap := make(map[string]string)
	for _, v := range vars {
		nameRef := getNameRefKey(v, appId)                       // returns `postgres.URL`
		idRef := getIDRefKey(v)                                  // returns `${{123123}}` where 123123 is UID of postgres.URL
		nameToIDMap[nameRef] = idRef                             // adding map[postgres.URL] = `${{123123}}`
		idToValueMap[strconv.FormatInt(v.UID, 10)] = v.TextValue // adding map[123123] = postgres://postgres:password@localhost:5432
	}

	return nameToIDMap, idToValueMap
}

type placeholder struct {
	AppName string
	Key     string
}

// extract placeholders
// ${{Postgres.URL}} will return { Postgres, URL}
// ${{Cloudness.URL}} will return { Cloudness, URL}
// ${{URL}} will return { "", URL}

func extractServiceVariablePlaceholders(input string) []placeholder {
	// Regex to match both ${{KEY}} and ${{DOMAIN.KEY}}
	re := regexp.MustCompile(`\${{\s*(?:(\w+)\.)?(\w+)\s*}}`)

	matches := re.FindAllStringSubmatch(input, -1)
	var results []placeholder

	for _, match := range matches {
		appName := match[1] // May be empty
		key := match[2]
		results = append(results, placeholder{
			AppName: appName,
			Key:     key,
		})
	}

	return results
}

func extractSecretValue(input string) (int, bool) {
	matches := varSecretRegex.FindStringSubmatch(input)
	if len(matches) == 2 {
		num, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Error().Err(err).Any("input", input).Msg("Error decoding secret value")
			return 0, false
		}
		return num, true
	}
	return 0, false
}
