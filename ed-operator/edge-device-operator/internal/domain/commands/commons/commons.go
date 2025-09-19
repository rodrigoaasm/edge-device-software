package command_commons

import (
	"ed-operator/internal/domain/entities"
	"strings"

	corekubev1 "k8s.io/api/core/v1"
)

func EnvMapToEnvVar(env entities.Env) []corekubev1.EnvVar {
	var envVars []corekubev1.EnvVar
	for key, value := range env {
		envVars = append(envVars, corekubev1.EnvVar{Name: key, Value: value})
	}

	return envVars
}

func EnvVarToEnvMap(envs []corekubev1.EnvVar) entities.Env {
	envMap := make(entities.Env)
	for _, env := range envs {
		envMap[env.Name] = env.Value
	}

	return envMap
}

func GetServiceName(name string, t corekubev1.ServiceType) string {
	return name + "-" + strings.ToLower(string(t))
}
