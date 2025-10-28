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
	Tls                   bool
	AlternativeDomain     string
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
	con.DeviceId = getEnvOrDefault("DEVICE_ID", "cc0001")
	con.CaptFrequecy, err = strconv.Atoi(getEnvOrDefault("CAPT_FREQUENCY", "100"))
	if err != nil {
		panic("CAPT_FREQUENCY must be an integer")
	}
	con.ParentId = getEnvOrDefault("DEVICE_PARENT_ID", "dd0001")
	con.ParentURL = getEnvOrDefault("DEVICE_PARENT_URL", "opc.tcp://192.168.0.5:4840")
	con.OpcInternalServerPort, err = strconv.Atoi(getEnvOrDefault("OPC_INTERNAL_SERVER_PORT", "4841"))
	if err != nil {
		panic("OPC_INTERNAL_SERVER_PORT must be an integer")
	}

	con.Tls, err = strconv.ParseBool(getEnvOrDefault("TLS_ENABLE", "false"))
	if err != nil {
		panic("TLS must be a boolean")
	}
	con.AlternativeDomain = getEnvOrDefault("ALTERNATIVE_DOMAIN", "172.31.212.236")

	return &con
}
