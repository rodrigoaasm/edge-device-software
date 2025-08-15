package entities

import (
	"errors"
	"regexp"
	"time"
)

type Device struct {
	DeviceId  string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDevice(deviceId string, url string) (*Device, error) {
	// ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)
	hexRegex := regexp.MustCompile(`^[A-Fa-f0-9]{4}$`)

	// if !ipRegex.MatchString(ip) {
	// 	return nil, errors.New("invalid IP address")
	// }
	if !hexRegex.MatchString(deviceId) {
		return nil, errors.New("invalid device id")
	}

	return &Device{
		DeviceId:  deviceId,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
