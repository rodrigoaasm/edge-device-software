package entities

import "time"

type Device struct {
	DeviceId  string    `json:"deviceId"`
	Ip        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
