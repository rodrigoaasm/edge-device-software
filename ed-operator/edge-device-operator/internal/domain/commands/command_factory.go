package domain_commands

import (
	"context"
	"ed-operator/internal/domain/dto"
	"ed-operator/internal/domain/interfaces"
	"errors"
)

type CommandFactory struct {
	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewCommandFactory(ctx context.Context, drc interfaces.IReconcilerClient) *CommandFactory {
	return &CommandFactory{ctx: ctx, reconcilerClient: drc}
}

func (c *CommandFactory) _mapper(command dto.CommandDTO) (ICommand, error) {
	if command.Command == "deploy" {
		return NewDeployCommand(command.Args.Name, command.Args.Image, command.Args.Env), nil
	}
	if command.Command == "undeploy" {
		return NewUndeployCommand(command.Args.Name), nil
	}
	if command.Command == "update" {
		return NewUpdateDeployCommand(command.Args.Name, command.Args.Image, command.Args.Env), nil
	}

	return nil, errors.New("Unknown command: " + command.Command)
}

func (c *CommandFactory) Make(commandDTO dto.CommandDTO) (ICommand, error) {
	command, err := c._mapper(commandDTO)
	if err != nil {
		return nil, err
	}

	command.SetContext(c.ctx)
	command.SetReconcilerClient(c.reconcilerClient)
	return command, nil
}
