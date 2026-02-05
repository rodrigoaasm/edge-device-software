package dtos

type ArgsDto struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	Env             map[string]string `json:"env"`
	PriorityProfile string            `json:"priorityProfile"`
	Port            uint16            `json:"port"`
	InternalPort    uint16            `json:"internalPort"`
	ExternalPort    uint16            `json:"externalPort"`
	RequestMemory   uint16            `json:"requestMemory"`
	LimitMemory     uint16            `json:"limitMemory"`
	RequestCPU      uint16            `json:"requestCPU"`
	LimitCPU        uint16            `json:"limitCPU"`
}

type CommandDTO struct {
	CorrelationId string  `json:"correlationId"`
	Command       string  `json:"command"`
	Args          ArgsDto `json:"args"`
}

type MessageDTO struct {
	DeviceId      string       `json:"deviceId"`
	CorrelationId string       `json:"correlationId"`
	Commands      []CommandDTO `json:"commands"`
}
