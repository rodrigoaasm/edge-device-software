package entities

import (
	"errors"
	"regexp"
	"time"

	"github.com/gopcua/opcua/ua"
)

type Device struct {
	DeviceId     string
	Url          string
	Interval     int
	Path         string
	Ns           int
	SecurityMode ua.MessageSecurityMode

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDevice(deviceId string, url string, interval int, path string, ns int, securityMode int) (*Device, error) {
	hexRegex := regexp.MustCompile(`^[A-Fa-f0-9]{4,6}$`)

	if !hexRegex.MatchString(deviceId) {
		return nil, errors.New("invalid device id")
	}

	if interval <= 0 {
		return nil, errors.New("invalid interval seconds")
	}

	if securityMode < 1 && securityMode > 3 {
		return nil, errors.New("invalid security mode")
	}

	securityModeCasted := ua.MessageSecurityMode(securityMode)
	return &Device{
		DeviceId:     deviceId,
		Url:          url,
		Path:         path,
		Ns:           ns,
		Interval:     interval,
		SecurityMode: securityModeCasted,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
