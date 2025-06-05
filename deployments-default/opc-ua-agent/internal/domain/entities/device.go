package entities

import (
	"errors"
	"regexp"
	"time"
)

type Device struct {
	DeviceId  string    `json:"deviceId"`
	Ip        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewDevice(deviceId string, ip string) (*Device, error) {
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)
	hexRegex := regexp.MustCompile(`^[A-Fa-f0-9]{6}$`)

	if !ipRegex.MatchString(ip) {
		return nil, errors.New("invalid IP address")
	}
	if !hexRegex.MatchString(deviceId) {
		return nil, errors.New("invalid device id")
	}

	return &Device{
		DeviceId:  deviceId,
		Ip:        ip,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
