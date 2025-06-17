package command_commons

import (
	"strings"

	corekubev1 "k8s.io/api/core/v1"
)

func GetServiceName(name string, t corekubev1.ServiceType) string {
	return name + "-" + strings.ToLower(string(t))
}
