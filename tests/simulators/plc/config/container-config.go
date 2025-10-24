package config

import (
	"os"
	"strconv"
)

type ContainerConfig struct {
	DeviceId              string
	ParentId              string
	ParentURL             string
	CaptFrequecy          int
	OpcInternalServerPort int
}

func getEnvOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment variable " + key + " is not set")
}

func NewContainerConfig() *ContainerConfig {
	var err error
	con := ContainerConfig{}
	con.DeviceId = getEnv("DEVICE_ID")
	con.CaptFrequecy, err = strconv.Atoi(getEnvOrDefault("CAPT_FREQUENCY", "100"))
	if err != nil {
		panic("CAPT_FREQUENCY must be an integer")
	}
	con.ParentId = getEnv("DEVICE_PARENT_ID")
	con.ParentURL = getEnvOrDefault("DEVICE_PARENT_URL", "opc.tcp://dell-g15:4840")
	con.OpcInternalServerPort, err = strconv.Atoi(getEnvOrDefault("OPC_INTERNAL_SERVER_PORT", "4840"))
	if err != nil {
		panic("OPC_INTERNAL_SERVER_PORT must be an integer")
	}

	return &con
}
