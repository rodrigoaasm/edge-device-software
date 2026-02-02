package entities

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
		PriorityProfile: priorityProfile,
		Port:            port,
		InternalPort:    internalPort,
		ExternalPort:    externalPort,
		RequestMemory:   requestMemory,
		LimitMemory:     limitMemory,
		RequestCPU:      requestCPU,
		LimitCPU:        limitCPU,
	}
}

type MicroserviceSimpleStatus struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Healthy bool   `json:"healthy"`
}
