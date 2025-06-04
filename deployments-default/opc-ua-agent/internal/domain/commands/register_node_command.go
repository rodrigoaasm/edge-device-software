package domain_commands

import (
	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type RegisterNodeCommand struct {
	Device entities.Device

	log              logr.Logger
	deviceRepository domain_interfaces.IDeviceRepository
}

func NewRegisterNodeCommand(device entities.Device, log logr.Logger, deviceRepository domain_interfaces.IDeviceRepository) *RegisterNodeCommand {
	return &RegisterNodeCommand{
		Device:           device,
		log:              log,
		deviceRepository: deviceRepository,
	}
}

func (c *RegisterNodeCommand) Execute() (interface{}, error) {
	c.log.Info("Registering node in database..")
	transManager, err := c.deviceRepository.Create(c.Device)
	if err != nil {
		return nil, err
	}

	c.log.Info("Transaction committed")
	return nil, transManager.Commit()
}
