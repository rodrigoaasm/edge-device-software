package register_node

import (
	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
	"opc.ua.agent/internal/domain/services"
)

type RegisterNodeCommand struct {
	correlationId string
	device        entities.Device

	log                         logr.Logger
	deviceRepository            domain_interfaces.IDeviceRepository
	outputDriverFactory         domain_interfaces.IOutputDriverFactory
	outputClientsManagerService *services.OutputClientsManagerService
}

func New(
	CorrelationId string,
	device entities.Device,
	log logr.Logger,
	outputDriverFactory domain_interfaces.IOutputDriverFactory,
	deviceRepository domain_interfaces.IDeviceRepository,
	outputClientsManagerService *services.OutputClientsManagerService,
) *RegisterNodeCommand {
	return &RegisterNodeCommand{
		correlationId:               CorrelationId,
		device:                      device,
		log:                         log,
		deviceRepository:            deviceRepository,
		outputDriverFactory:         outputDriverFactory,
		outputClientsManagerService: outputClientsManagerService,
	}
}

func (c *RegisterNodeCommand) GetCorrelationId() string {
	return c.correlationId
}

func (c *RegisterNodeCommand) Execute() (interface{}, error) {
	c.log.Info("Registering node in database..")
	transManager, err := c.deviceRepository.Create(c.device)
	if err != nil {
		return nil, err
	}

	c.log.Info("Connecting to opcua server..")
	opcDriver := c.outputDriverFactory.Make(c.device)
	if err = opcDriver.Connect(); err != nil {
		c.log.Error(err, "Connection failed to opcua server.")
		return nil, err
	}
	c.outputClientsManagerService.AddClient(opcDriver)
	c.log.Info("Connected to opcua server.")

	c.log.Info("Transaction committed")
	return nil, transManager.Commit()
}
