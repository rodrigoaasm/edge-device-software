package list_nodes

import (
	"fmt"

	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/dto"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
	"opc.ua.agent/internal/domain/services"
)

type ListNodesCommand struct {
	CorrelationId string

	log                         logr.Logger
	deviceRepository            domain_interfaces.IDeviceRepository
	outputClientsManagerService *services.OutputClientsManagerService
}

func New(
	correlationId string,
	log logr.Logger,
	deviceRepository domain_interfaces.IDeviceRepository,
	outputClientsManagerService *services.OutputClientsManagerService,
) *ListNodesCommand {
	return &ListNodesCommand{
		CorrelationId:               correlationId,
		log:                         log,
		deviceRepository:            deviceRepository,
		outputClientsManagerService: outputClientsManagerService,
	}
}

func (c *ListNodesCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (c *ListNodesCommand) Execute() (interface{}, error) {
	var devicesOutput []dto.Device

	c.log.Info("Searching nodes from database..")
	devices, err := c.deviceRepository.List()
	if err != nil {
		return nil, err
	}

	// transform
	for _, deviceEnt := range devices {
		deviceDTO := dto.Device{
			DeviceId:  deviceEnt.DeviceId,
			Url:       deviceEnt.Url,
			Active:    false,
			CreatedAt: deviceEnt.CreatedAt,
			UpdatedAt: deviceEnt.UpdatedAt,
		}
		client := c.outputClientsManagerService.GetClient(deviceEnt.DeviceId)
		if client != nil {
			deviceDTO.Active = true
		}

		devicesOutput = append(devicesOutput, deviceDTO)
	}

	c.log.Info(fmt.Sprintf("Nodes Found: %d", len(devices)))
	return devicesOutput, nil
}
