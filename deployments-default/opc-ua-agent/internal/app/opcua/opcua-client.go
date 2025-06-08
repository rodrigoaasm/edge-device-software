package opcua

import (
	"context"
	"flag"

	"github.com/go-logr/logr"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"opc.ua.agent/internal/domain/entities"
)

type OPCUAClient struct {
	Device entities.Device

	log logr.Logger
}

func NewOPCUAClient(device entities.Device, log logr.Logger) *OPCUAClient {
	return &OPCUAClient{
		Device: device,

		log: log,
	}
}

func (c *OPCUAClient) Connect() error {
	var endpoint = flag.String("endpoint", "opc.tcp://"+c.Device.Url, "OPC UA Endpoint URL")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()

	c.log.Info("Connecting to endpoint: " + *endpoint)
	eps, err := opcua.GetEndpoints(context.Background(), *endpoint)
	if err != nil {
		return err
	}

	for _, ep := range eps {
		c.log.Info(ep.EndpointURL, ep.SecurityPolicyURI, ep.SecurityMode)
	}
	return nil
}

func (c *OPCUAClient) Publish(topic string, message string) error {
	c.log.Info("Publishing message: " + message + " to topic: " + topic)
	return nil
}
