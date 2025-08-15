package utils

import (
	"errors"

	"github.com/go-logr/logr"
)

func EmitError(log logr.Logger, message string) error {
	log.Error(nil, message)

	return errors.New(message)
}
