package domain_commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	actuate_opc "opc.ua.agent/internal/domain/commands/actuate"
	list_nodes "opc.ua.agent/internal/domain/commands/list_nodes"
	register_node "opc.ua.agent/internal/domain/commands/register_node"
	remove_node "opc.ua.agent/internal/domain/commands/remove_node"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
	"opc.ua.agent/internal/domain/services"
	"opc.ua.agent/internal/utils"
)

type CommandFactory struct {
	ctx                         context.Context
	deviceRepository            domain_interfaces.IDeviceRepository
	outputDriverFactory         domain_interfaces.IOutputDriverFactory
	outputClientsManagerService *services.OutputClientsManagerService
}

func NewCommandFactory(
	ctx context.Context,
	outputDriverFactory domain_interfaces.IOutputDriverFactory,
	deviceRepository domain_interfaces.IDeviceRepository,
	outputClientsManagerService *services.OutputClientsManagerService,
) *CommandFactory {
	return &CommandFactory{
		ctx:                         ctx,
		deviceRepository:            deviceRepository,
		outputDriverFactory:         outputDriverFactory,
		outputClientsManagerService: outputClientsManagerService,
	}
}

func (c *CommandFactory) _mapper(command map[string]interface{}) (ICommand, error) {
	log := logr.FromContextOrDiscard(c.ctx)

	fmt.Printf("Command: %v\n", command)

	commandField, ok := command["command"].(string)
	if !ok {
		return nil, utils.EmitError(log, "The command is null.")
	}

	correlationId, ok := command["correlationId"].(string)
	if !ok {
		return nil, utils.EmitError(log, "The command not have correlation id.")
	}

	if command["args"] == nil {
		return nil, utils.EmitError(log, "The command not have args.")
	}

	if commandField == "register_opc" {
		device, err := entities.NewDevice(
			command["args"].(map[string]interface{})["nodeId"].(string),
			command["args"].(map[string]interface{})["url"].(string),
			int(command["args"].(map[string]interface{})["intervalSeconds"].(float64)),
			command["args"].(map[string]interface{})["path"].(string),
			int(command["args"].(map[string]interface{})["ns"].(float64)),
		)
		if err != nil {
			return nil, err
		}
		return register_node.New(
			correlationId,
			*device,
			log,
			c.outputDriverFactory,
			c.deviceRepository,
			c.outputClientsManagerService,
		), nil
	}

	if commandField == "list_opc" {
		return list_nodes.New(correlationId, log, c.deviceRepository, c.outputClientsManagerService), nil
	}

	if commandField == "remove_opc" {
		return remove_node.New(
			correlationId,
			command["args"].(map[string]interface{})["nodeId"].(string),
			log,
			c.deviceRepository,
			c.outputClientsManagerService,
		), nil
	}

	return nil, errors.New("Unknown command: " + commandField)
}

func (c *CommandFactory) Make(commandDTO map[string]interface{}) (ICommand, error) {
	command, err := c._mapper(commandDTO)
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (c *CommandFactory) MakeActuateCmd(timestamp int64, deviceId string, data map[string]interface{}) (ICommand, error) {
	log := logr.FromContextOrDiscard(c.ctx)

	return actuate_opc.New(
		timestamp,
		data,
		deviceId,
		log,
		c.outputClientsManagerService,
	), nil
}
