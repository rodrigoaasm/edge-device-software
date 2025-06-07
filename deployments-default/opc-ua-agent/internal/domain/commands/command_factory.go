package domain_commands

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	domain_command_list_nodes "opc.ua.agent/internal/domain/commands/list_nodes"
	domain_commands_register_node "opc.ua.agent/internal/domain/commands/register_node"
	domain_commands_remove_node "opc.ua.agent/internal/domain/commands/remove_node"
	"opc.ua.agent/internal/domain/dto"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type CommandFactory struct {
	ctx                 context.Context
	deviceRepository    domain_interfaces.IDeviceRepository
	outputDriverFactory domain_interfaces.IOutputDriverFactory
}

func NewCommandFactory(
	ctx context.Context,
	outputDriverFactory domain_interfaces.IOutputDriverFactory,
	deviceRepository domain_interfaces.IDeviceRepository,
) *CommandFactory {
	return &CommandFactory{
		ctx:                 ctx,
		deviceRepository:    deviceRepository,
		outputDriverFactory: outputDriverFactory,
	}
}

func (c *CommandFactory) _mapper(command dto.CommandDTO) (ICommand, error) {
	log := logr.FromContextOrDiscard(c.ctx)

	if command.Command == "register_opc" {
		device, err := entities.NewDevice(command.Args.NodeId, command.Args.Url)
		if err != nil {
			return nil, err
		}
		return domain_commands_register_node.New(command.CorrelationId, *device, log, c.outputDriverFactory, c.deviceRepository), nil
	}

	if command.Command == "list_opc" {
		return domain_command_list_nodes.New(command.CorrelationId, log, c.deviceRepository), nil
	}

	if command.Command == "remove_opc" {
		return domain_commands_remove_node.New(command.CorrelationId, command.Args.NodeId, log, c.deviceRepository), nil
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
