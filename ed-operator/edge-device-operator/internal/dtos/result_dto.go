package dtos

type ResultDto struct {
	DeviceId    string      `json:"deviceId"`
	Correlation string      `json:"correlation"`
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
	Message     string      `json:"message"`
	Timestamp   int64       `json:"timestamp"`
}

type AckDto struct {
	DeviceId    string `json:"deviceId"`
	Correlation string `json:"correlation"`
	Timestamp   int64  `json:"timestamp"`
	Ack         bool   `json:"ack"`
}
