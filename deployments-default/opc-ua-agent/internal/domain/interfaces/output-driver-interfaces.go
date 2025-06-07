package domain_interfaces

import "opc.ua.agent/internal/domain/entities"

type IOutputDriver interface {
	Connect() error
	Publish(topic string, message string) error
}

type IOutputDriverFactory interface {
	Make(entities.Device) IOutputDriver
}
