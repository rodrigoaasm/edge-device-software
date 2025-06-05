package dto

type ArgsDto struct {
	NodeId string `json:"nodeId"`
	Ip     string `json:"ip"`
}

type CommandDTO struct {
	CorrelationId string  `json:"correlationId"`
	DeviceId      string  `json:"deviceId"`
	Command       string  `json:"command"`
	Args          ArgsDto `json:"args"`
}
