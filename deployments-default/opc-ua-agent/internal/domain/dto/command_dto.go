package dto

type CommandDTO struct {
	DeviceId string `json:"deviceId"`
	Command  string `json:"command"`
	Ip       string `json:"ip"`
}

/*
{"name":"telegraf","image":"influxdb:2.7.11-alpine","env":{}}
*/
