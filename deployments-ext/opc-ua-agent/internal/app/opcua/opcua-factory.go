package opcua

import (
	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type OPCUAClientFactory struct {
	log logr.Logger
}

func NewOPCUAClientFactory(log logr.Logger) *OPCUAClientFactory {
	return &OPCUAClientFactory{
		log: log,
	}
}

func (c *OPCUAClientFactory) Make(device entities.Device) domain_interfaces.IOutputDriver {
	return NewOPCUAClient(device, c.log)
}
