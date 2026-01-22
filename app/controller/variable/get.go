package variable

import (
	"context"
	"net/url"
	"strconv"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func (c *Controller) UpdateDefaultVariables(
	ctx context.Context,
	server *types.Server,
	tenant *types.Tenant,
	project *types.Project,
	env *types.Environment,
	app *types.Application,
) error {
	varsToDelete := make([]string, 0)
	vars := []*types.Variable{
		newVariable(env.ID, app.ID, SystemVarTeamID, strconv.FormatInt(tenant.UID, 10), enum.VariableTypeBuildAndRun),
		newVariable(env.ID, app.ID, SystemVarProjectID, strconv.FormatInt(project.UID, 10), enum.VariableTypeBuildAndRun),
		newVariable(env.ID, app.ID, SystemVarEnvironmentID, strconv.FormatInt(env.UID, 10), enum.VariableTypeBuildAndRun),
		newVariable(env.ID, app.ID, SystemVarAppID, strconv.FormatInt(app.UID, 10), enum.VariableTypeBuildAndRun),
	}

	if app.PrivateDomain != "" {
		if server.Type == enum.ServerTypeK8s {
			vars = append(vars, newVariable(env.ID, app.ID, SystemVarAppPrivateDomain, app.PrivateDomain, enum.VariableTypeBuildAndRun))
		}
		vars = append(vars, newVariable(env.ID, app.ID, SystemVarServiceName, app.PrivateDomain, enum.VariableTypeBuildAndRun))
	}

	if app.Spec.Networking.TCPProxies != nil {
		vars = append(vars, newVariable(env.ID, app.ID, SystemVarAppTCPPort, strconv.Itoa(app.Spec.Networking.TCPProxies.TCPPort), enum.VariableTypeBuildAndRun))
	} else {
		varsToDelete = append(varsToDelete, SystemVarAppTCPPort)
	}

	if app.Domain != "" {
		vars = append(vars, newVariable(env.ID, app.ID, SystemVarAppPublicDomain, app.Domain, enum.VariableTypeBuildAndRun))
		u, err := url.Parse(app.Domain)
		if err == nil { //NOTE: Ignore error
			vars = append(vars, newVariable(env.ID, app.ID, SystemVarAppPublicURL, u.Hostname(), enum.VariableTypeBuildAndRun))
		}
	} else {
		varsToDelete = append(varsToDelete, SystemVarAppPublicDomain)
		varsToDelete = append(varsToDelete, SystemVarAppPublicURL)
	}

	if err := c.variableStore.UpsertMany(ctx, vars); err != nil {
		return err
	}

	if len(varsToDelete) > 0 {
		if err := c.variableStore.DeleteByKeys(ctx, app.ID, varsToDelete); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting old variables")
			return err
		}
	}
	return nil
}
