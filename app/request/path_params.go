package request

import (
	"net/http"
	"strconv"
)

const (
	PathParamTenantUID      = "tenant_uid"
	PathParamProjectUID     = "project_uid"
	PathParamEnvironmentUID = "environment_uid"
	PathParamSourceUID      = "source_uid"
	PathParamApplicationUID = "application_uid"
	PathParamDeploymentUID  = "deployment_uid"
	PathParamDeploymentStep = "deployment_step"
	PathParamVolumeUID      = "volume_uid"
	PathParamAuthProvider   = "auth_provider"
	PathParamVariableUID    = "variable_uid"
	PathParamTemplateID     = "template_id"
)

func GetTenantUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamTenantUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetProjectUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamProjectUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetEnvironmentUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamEnvironmentUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetSourceUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamSourceUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetApplicationUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamApplicationUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetDeploymentUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamDeploymentUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetDeploymentStepFromPath(r *http.Request) (string, error) {
	return PathParamOrError(r, PathParamDeploymentStep)
}

func GetVolumeUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamVolumeUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetVariableUIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamVariableUID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetTemplateIDFromPath(r *http.Request) (int64, error) {
	id, err := PathParamOrError(r, PathParamTemplateID)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(id, 10, 64)
}

func GetAuthProviderFromPath(r *http.Request) (string, error) {
	return PathParamOrError(r, PathParamAuthProvider)
}
