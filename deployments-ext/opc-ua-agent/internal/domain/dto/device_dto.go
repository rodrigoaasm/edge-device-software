package dto

import "time"

type Device struct {
	DeviceId  string    `json:"deviceId"`
	Url       string    `json:"url"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
