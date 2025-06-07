package domain_commands_register_node

import (
	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type RegisterNodeCommand struct {
	CorrelationId string
	Device        entities.Device

	log                 logr.Logger
	deviceRepository    domain_interfaces.IDeviceRepository
	outputDriverFactory domain_interfaces.IOutputDriverFactory
}

func New(
	CorrelationId string,
	device entities.Device,
	log logr.Logger,
	outputDriverFactory domain_interfaces.IOutputDriverFactory,
	deviceRepository domain_interfaces.IDeviceRepository,
) *RegisterNodeCommand {
	return &RegisterNodeCommand{
		CorrelationId:       CorrelationId,
		Device:              device,
		log:                 log,
		deviceRepository:    deviceRepository,
		outputDriverFactory: outputDriverFactory,
	}
}

func (c *RegisterNodeCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (c *RegisterNodeCommand) Execute() (interface{}, error) {
	c.log.Info("Registering node in database..")
	transManager, err := c.deviceRepository.Create(c.Device)
	if err != nil {
		return nil, err
	}

	c.log.Info("Connecting to opcua server..")
	opcDriver := c.outputDriverFactory.Make(c.Device)
	if err = opcDriver.Connect(); err != nil {
		c.log.Error(err, "Connection failed to opcua server.")
		return nil, err
	}
	c.log.Info("Connected to opcua server.")

	c.log.Info("Transaction committed")
	return nil, transManager.Commit()
}
