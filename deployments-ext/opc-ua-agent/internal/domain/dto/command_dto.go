package dto

type ArgsDto struct {
	NodeId string `json:"nodeId"`
	Url    string `json:"url"`
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
