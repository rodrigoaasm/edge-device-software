package domain_commands

import (
	"context"
	"ed-operator/internal/domain/commands/deploy"
	"ed-operator/internal/domain/commands/undeploy"
	"ed-operator/internal/domain/commands/update_deploy"
	"ed-operator/internal/domain/dto"
	"ed-operator/internal/domain/entities"
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
	if command.Command == "undeploy" {
		return undeploy.NewUndeployCommand(command.CorrelationId, command.Args.Name), nil
	} else {
		microservice := entities.NewMicroservice(
			command.Args.Name,
			command.Args.Image,
			command.Args.Env,
			command.Args.PriorityProfile,
			command.Args.Port,
			command.Args.InternalPort,
			command.Args.ExternalPort,
			command.Args.RequestMemory,
			command.Args.LimitMemory,
			command.Args.RequestCPU,
			command.Args.LimitCPU,
		)

		if command.Command == "deploy" {
			return deploy.NewDeployCommand(
				command.CorrelationId,
				microservice,
			), nil
		}
		if command.Command == "update" {
			return update_deploy.NewUpdateDeployCommand(
				command.CorrelationId,
				microservice,
			), nil
		}
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
