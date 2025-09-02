package domain_interfaces

import "opc.ua.agent/internal/domain/entities"

type IOutputDriver interface {
	Connect() error
	Disconnect() error
	GetClientId() string
	PublishData(key string, message interface{}) error
}

type IOutputDriverFactory interface {
	Make(entities.Device) IOutputDriver
}
