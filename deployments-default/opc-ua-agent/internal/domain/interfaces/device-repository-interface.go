package domain_interfaces

import "opc.ua.agent/internal/domain/entities"

type IDeviceRepository interface {
	Create(entities.Device) (ITransactionManager, error)
	List() ([]entities.Device, error)
}
