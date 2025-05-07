package dto

type ArgsDto struct {
	Name  string            `json:"name"`
	Image string            `json:"image"`
	Env   map[string]string `json:"env"`
}

type CommandDTO struct {
	Command string  `json:"command"`
	Args    ArgsDto `json:"args"`
}

/*
{"name":"telegraf","image":"influxdb:2.7.11-alpine","env":{}}
*/
