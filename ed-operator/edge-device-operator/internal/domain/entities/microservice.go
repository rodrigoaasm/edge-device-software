package entities

import "ed-operator/internal/utils"

type Env map[string]string
type Microservice struct {
	Name            string
	Image           string
	Env             Env
	PriorityProfile string
	Port            uint16
	InternalPort    uint16
	ExternalPort    uint16
	RequestMemory   uint16
	LimitMemory     uint16
	RequestCPU      uint16
	LimitCPU        uint16
}

func NewMicroservice(
	name string,
	image string,
	env Env,
	priorityProfile string,
	port uint16,
	internalPort uint16,
	externalPort uint16,
	requestMemory uint16,
	limitMemory uint16,
	requestCPU uint16,
	limitCPU uint16,
) *Microservice {
	return &Microservice{
		Name:            name,
		Image:           image,
		Env:             env,
		PriorityProfile: utils.GetValueOrDefault(priorityProfile, PRIORITY_PROFILE_SIMPLE_SERVICE),
		Port:            port,
		InternalPort:    internalPort,
		ExternalPort:    externalPort,
		RequestMemory:   utils.GetValueOrDefault(requestMemory, 128),
		LimitMemory:     utils.GetValueOrDefault(limitMemory, 256),
		RequestCPU:      utils.GetValueOrDefault(requestCPU, 100),
		LimitCPU:        utils.GetValueOrDefault(limitCPU, 500),
	}
}

type MicroserviceSimpleStatus struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Healthy bool   `json:"healthy"`
}
