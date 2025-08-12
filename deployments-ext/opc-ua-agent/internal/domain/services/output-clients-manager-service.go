package services

import (
	"github.com/go-logr/logr"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type OutputClientsManagerService struct {
	clients map[string]domain_interfaces.IOutputDriver

	log logr.Logger
}

func NewOutputClientsManagerService(log logr.Logger) *OutputClientsManagerService {
	return &OutputClientsManagerService{
		clients: make(map[string]domain_interfaces.IOutputDriver),
		log:     log,
	}
}

func (ms *OutputClientsManagerService) AddClient(client domain_interfaces.IOutputDriver) {
	ms.clients[client.GetClientId()] = client
	ms.log.Info("Client added in output clients manager")
}

func (ms *OutputClientsManagerService) RemoveClient(clientId string) error {
	if ms.clients[clientId] != nil {
		ms.log.Info("Client found in output clients manager. Disconnecting...")
		if err := ms.clients[clientId].Disconnect(); err != nil {
			return err
		}
		ms.log.Info("Client removed from output clients manager")
		delete(ms.clients, clientId)
	}

	return nil
}
