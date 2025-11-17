package variable

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudness-io/cloudness/dag"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type variableIndex struct {
	index        map[string]*types.Variable
	nameToIDMap  map[string]string
	idToValueMap map[string]string
}

func newVariableIndex() *variableIndex {
	return &variableIndex{
		index:        make(map[string]*types.Variable),
		nameToIDMap:  make(map[string]string),
		idToValueMap: make(map[string]string),
	}
}

func (i *variableIndex) AddOrUpdate(serviceName string, key string, v *types.Variable) {
	nameRef := fmt.Sprintf("%s.%s", serviceName, v.Key)        // returns `postgres.URL`
	idRef := fmt.Sprintf("${{%d}}", v.UID)                     // returns `${{123123}}` where 123123 is UID of postgres.URL
	i.nameToIDMap[nameRef] = idRef                             // adding map[postgres.URL] = `${{123123}}`
	i.idToValueMap[strconv.FormatInt(v.UID, 10)] = v.TextValue // adding map[123123] = postgres://postgres:password@localhost:5432
	i.index[nameRef] = v
}

func (c *Controller) AddTemplateVariables(ctx context.Context, envID int64, serviceAppMap map[string]*types.Application, template *types.Template) error {

	// variables
	dag := dag.NewGraph[string]()
	vIndex := newVariableIndex()

	for _, service := range template.Spec.Services {
		app := serviceAppMap[service.Name]

		for _, v := range service.Variables {
			varKey := fmt.Sprintf("%s.%s", service.Name, v.Key)

			ph := extractServiceVariablePlaceholders(v.Value)
			if len(ph) == 0 {
				dag.AddVertex(varKey)
			} else {
				for _, p := range ph {
					if p.AppName == "" {
						p.AppName = service.Name //current service, so replace with service name
						v.Value = strings.ReplaceAll(v.Value,
							fmt.Sprintf("${{%s}}", p.Key),
							fmt.Sprintf("${{%s.%s}}", p.AppName, p.Key)) // replace {{POSTGRES_URL}} with {{postgres.POSTGRES_URL}}
					}
					dag.AddEdge(fmt.Sprintf("%s.%s", p.AppName, p.Key), varKey)
				}
			}
			vIndex.AddOrUpdate(service.Name, v.Key, newVariable(envID, app.ID, v.Key, v.Value, v.Type))
		}
	}

	for service, app := range serviceAppMap {
		vars, err := c.variableStore.List(ctx, envID, app.ID)
		if err != nil {
			return err
		}
		for _, v := range vars {
			vIndex.AddOrUpdate(service, v.Key, v)
		}
	}

	sorted, err := dag.TopoSort()
	if err != nil {
		return err
	}

	for _, s := range sorted {
		v, exists := vIndex.index[s]
		if !exists {
			log.Ctx(ctx).Warn().Any("Key", s).Msg("variable not found")
			continue
		}

		//check if its a secret value
		num, isSecret := extractSecretValue(v.Value)
		if isSecret {
			secretValue, secretRef := c.generateSecret(num)
			v.TextValue = secretValue
			v.Value = secretRef
		} else {
			ref := getRef(v.Value)
			if len(ref) > 0 {
				v.Value = replaceWithMap(v.Value, vIndex.nameToIDMap)
				v.TextValue = replaceWithMap(v.Value, vIndex.idToValueMap)
			} else {
				v.TextValue = v.Value
			}
		}
		sName, sKey := splitServiceVarKey(s)
		vIndex.AddOrUpdate(sName, sKey, v)
	}

	newVars := make([]*types.Variable, 0, len(vIndex.index))
	for _, v := range vIndex.index {
		newVars = append(newVars, v)
	}
	return c.UpsertMany(ctx, newVars)
}

func splitServiceVarKey(key string) (string, string) {
	split := strings.Split(key, ".")
	switch len(split) {
	case 2:
		return split[0], split[1]
	default:
		return "", split[0]
	}
}
