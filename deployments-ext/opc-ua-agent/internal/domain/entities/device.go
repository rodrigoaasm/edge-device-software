package entities

import (
	"errors"
	"regexp"
	"time"
)

type Device struct {
	DeviceId        string
	Url             string
	IntervalSeconds int
	Path            string
	Ns              int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewDevice(deviceId string, url string, intervalSeconds int, path string, ns int) (*Device, error) {
	hexRegex := regexp.MustCompile(`^[A-Fa-f0-9]{4,6}$`)

	if !hexRegex.MatchString(deviceId) {
		return nil, errors.New("invalid device id")
	}

	if intervalSeconds <= 0 {
		return nil, errors.New("invalid interval seconds")
	}

	return &Device{
		DeviceId:        deviceId,
		Url:             url,
		Path:            path,
		Ns:              ns,
		IntervalSeconds: intervalSeconds,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}
