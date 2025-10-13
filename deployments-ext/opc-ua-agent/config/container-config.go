package config

import "os"

type ContainerConfig struct {
	ConsumerTopic string
	ResultsTopic  string
	DataTopic     string
	BrokerUrl     string
	DeviceId      string
	LogLevel      string
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
	con := ContainerConfig{}
	con.ConsumerTopic = getEnvOrDefault("CONSUMER_TOPIC", "device/config/opc")
	con.ResultsTopic = getEnvOrDefault("RESULTS_TOPIC", "device/results")
	con.DataTopic = getEnvOrDefault("DATA_TOPIC", "device/data")
	con.DeviceId = getEnv("DEVICE_ID")
	con.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")
	con.BrokerUrl = getEnvOrDefault("BROKER_URL", "tcp://nanomq:1883")

	return &con
}
