package dto

type ArgsDto struct {
	Name  string            `json:"name"`
	Image string            `json:"image"`
	Env   map[string]string `json:"env"`
}

type CommandDTO struct {
	CorrelationId string  `json:"correlationId"`
	Command       string  `json:"command"`
	Args          ArgsDto `json:"args"`
}

type MessageDTO struct {
	DeviceId string       `json:"deviceId"`
	Commands []CommandDTO `json:"commands"`
}
