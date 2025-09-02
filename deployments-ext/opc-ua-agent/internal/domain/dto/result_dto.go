package dto

type ResultDto struct {
	CorrelationId string      `json:"correlationId"`
	Success       bool        `json:"success"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message"`
	Timestamp     int64       `json:"timestamp"`
}

type AckDto struct {
	CorrelationId string `json:"correlationId"`
	Timestamp     int64  `json:"timestamp"`
	Ack           bool   `json:"ack"`
}
