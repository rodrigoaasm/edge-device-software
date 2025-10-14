package dto

import "time"

type Device struct {
	DeviceId  string    `json:"deviceId"`
	Url       string    `json:"url"`
	Active    bool      `json:"active"`
	Interval  int       `json:"interval"`
	Path      string    `json:"path"`
	Ns        int       `json:"ns"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
