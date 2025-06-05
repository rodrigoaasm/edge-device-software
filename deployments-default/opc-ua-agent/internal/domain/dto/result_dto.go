package dto

type ResultDto struct {
	CorrelationId string      `json:"correlationId"`
	Success       bool        `json:"success"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message"`
}
