package domain_commands

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/dto"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type CommandFactory struct {
	ctx              context.Context
	deviceRepository domain_interfaces.IDeviceRepository
}

func NewCommandFactory(ctx context.Context, deviceRepository domain_interfaces.IDeviceRepository) *CommandFactory {
	return &CommandFactory{
		ctx:              ctx,
		deviceRepository: deviceRepository,
	}
}

func (c *CommandFactory) _mapper(command dto.CommandDTO) (ICommand, error) {
	log := logr.FromContextOrDiscard(c.ctx)

	if command.Command == "register_opc" {
		device := entities.Device{DeviceId: command.DeviceId, Ip: command.Ip}
		return NewRegisterNodeCommand(device, log, c.deviceRepository), nil
	}

	if command.Command == "list_opc" {
		return NewListNodesCommand(log, c.deviceRepository), nil
	}

	return nil, errors.New("Unknown command: " + command.Command)
}

func (c *CommandFactory) Make(commandDTO dto.CommandDTO) (ICommand, error) {
	command, err := c._mapper(commandDTO)
	if err != nil {
		return nil, err
	}
	return command, nil
}
