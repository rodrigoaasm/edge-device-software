package list_nodes

import (
	"fmt"

	"github.com/go-logr/logr"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type ListNodesCommand struct {
	CorrelationId string

	log              logr.Logger
	deviceRepository domain_interfaces.IDeviceRepository
}

func New(correlationId string, log logr.Logger, deviceRepository domain_interfaces.IDeviceRepository) *ListNodesCommand {
	return &ListNodesCommand{
		CorrelationId:    correlationId,
		log:              log,
		deviceRepository: deviceRepository,
	}
}

func (c *ListNodesCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (c *ListNodesCommand) Execute() (interface{}, error) {
	c.log.Info("Searching nodes from database..")
	devices, err := c.deviceRepository.List()
	if err != nil {
		return nil, err
	}

	c.log.Info(fmt.Sprintf("Nodes Found: %d", len(devices)))
	return devices, nil
}
