package dto

type ArgsDto struct {
	NodeId string `json:"nodeId"`
	Url    string `json:"url"`
}

type CommandDTO struct {
	CorrelationId string  `json:"correlationId"`
	DeviceId      string  `json:"deviceId"`
	Command       string  `json:"command"`
	Args          ArgsDto `json:"args"`
}
