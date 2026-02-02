package domain_commands

import (
	"context"
	"ed-operator/internal/domain/commands/deploy"
	"ed-operator/internal/domain/commands/healthcheck"
	"ed-operator/internal/domain/commands/list_deploy"
	"ed-operator/internal/domain/commands/undeploy"
	"ed-operator/internal/domain/commands/update_deploy"
	"ed-operator/internal/domain/dto"
	"ed-operator/internal/domain/entities"
	"ed-operator/internal/domain/interfaces"
	"ed-operator/internal/utils"
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
	var microservice *entities.Microservice
	if command.Command == "deploy" {
		microservice = entities.NewMicroservice(
			command.Args.Name,
			command.Args.Image,
			command.Args.Env,
			utils.GetValueOrDefault(command.Args.PriorityProfile, entities.PRIORITY_PROFILE_SIMPLE_SERVICE),
			command.Args.Port,
			command.Args.InternalPort,
			command.Args.ExternalPort,
			utils.GetValueOrDefault(command.Args.RequestMemory, 128),
			utils.GetValueOrDefault(command.Args.LimitMemory, 256),
			utils.GetValueOrDefault(command.Args.RequestCPU, 100),
			utils.GetValueOrDefault(command.Args.LimitCPU, 500),
		)
	}

	if command.Command == "update" {
		microservice = entities.NewMicroservice(
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
	}

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
	if command.Command == "healthcheck" {
		return healthcheck.NewHealthCheckCommand(
			command.CorrelationId,
		), nil
	}
	if command.Command == "list_deploy" {
		return list_deploy.NewListDeployCommand(
			command.CorrelationId,
		), nil
	}
	if command.Command == "undeploy" {
		return undeploy.NewUndeployCommand(command.CorrelationId, command.Args.Name), nil
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
