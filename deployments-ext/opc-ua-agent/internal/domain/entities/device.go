package entities

import (
	"errors"
	"regexp"
	"time"
)

type Device struct {
	DeviceId  string
	Url       string
	Interval  int
	Path      string
	Ns        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDevice(deviceId string, url string, interval int, path string, ns int) (*Device, error) {
	hexRegex := regexp.MustCompile(`^[A-Fa-f0-9]{4,6}$`)

	if !hexRegex.MatchString(deviceId) {
		return nil, errors.New("invalid device id")
	}

	if interval <= 0 {
		return nil, errors.New("invalid interval seconds")
	}

	return &Device{
		DeviceId:  deviceId,
		Url:       url,
		Path:      path,
		Ns:        ns,
		Interval:  interval,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
