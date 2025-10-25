package opcua

import (
	"github.com/go-logr/logr"
	"opc.ua.agent/config"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type OPCUAClientFactory struct {
	log          logr.Logger
	OutputDriver domain_interfaces.IOutputDriver
	config       *config.ContainerConfig
}

func NewOPCUAClientFactory(log logr.Logger, config *config.ContainerConfig) *OPCUAClientFactory {
	return &OPCUAClientFactory{
		log:    log,
		config: config,
	}
}

func (c *OPCUAClientFactory) Make(device entities.Device) domain_interfaces.IOutputDriver {
	return NewOPCUAClient(device, c.log, c.OutputDriver, c.config)
}
