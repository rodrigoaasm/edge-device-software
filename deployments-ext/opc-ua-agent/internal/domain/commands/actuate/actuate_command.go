package actuate_opc

import (
	"fmt"

	"github.com/go-logr/logr"
	"opc.ua.agent/internal/domain/services"
)

type ActuateCommand struct {
	timestamp int64
	data      map[string]interface{}
	deviceId  string

	clientMngService *services.OutputClientsManagerService
	log              logr.Logger
}

func New(
	timestamp int64,
	data map[string]interface{},
	deviceId string,
	log logr.Logger,
	clientMngService *services.OutputClientsManagerService,
) *ActuateCommand {
	return &ActuateCommand{
		timestamp:        timestamp,
		data:             data,
		log:              log,
		deviceId:         deviceId,
		clientMngService: clientMngService,
	}
}

func (c *ActuateCommand) GetCorrelationId() string {
	return fmt.Sprintf("%d", c.timestamp)
}

func (c *ActuateCommand) Execute() (interface{}, error) {
	c.log.Info("Verify if actuate belongs to a opc-device...")
	if client := c.clientMngService.GetClient(c.deviceId); client != nil {
		c.log.Info("Actuate belongs to a opc-device (" + c.deviceId + "). Sending actuate to opcua server.")
		delete(c.data, "deviceId")
		delete(c.data, "timestamp")
		for key, value := range c.data {
			if err := client.PublishData(key, value); err != nil {
				c.log.Error(err, "Unable to send actuate to opcua server.")
				return nil, err
			}
		}
		c.log.Info("Actuate to opcua server sent.")
		return nil, nil
	}

	c.log.Info("Actuate not belongs to a opc-device. Ignore")
	return nil, nil
}
