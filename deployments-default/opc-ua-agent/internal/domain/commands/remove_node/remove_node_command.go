package remove_node

import (
	"github.com/go-logr/logr"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type RemoveNodeCommand struct {
	CorrelationId string
	DeviceId      string

	log              logr.Logger
	deviceRepository domain_interfaces.IDeviceRepository
}

func New(correlationId string, deviceId string, log logr.Logger, deviceRepository domain_interfaces.IDeviceRepository) *RemoveNodeCommand {
	return &RemoveNodeCommand{
		CorrelationId:    correlationId,
		DeviceId:         deviceId,
		log:              log,
		deviceRepository: deviceRepository,
	}
}

func (c *RemoveNodeCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (c *RemoveNodeCommand) Execute() (interface{}, error) {
	c.log.Info("Removing node from database..")
	transManager, err := c.deviceRepository.Remove(c.DeviceId)
	if err != nil {
		return nil, err
	}

	err = transManager.Commit()
	if err != nil {
		return nil, err
	}
	c.log.Info("Node removed from database.")

	return nil, nil
}
