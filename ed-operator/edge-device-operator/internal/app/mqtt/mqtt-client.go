package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"ed-operator/config"
	domain_commands "ed-operator/internal/domain/commands"
	"ed-operator/internal/domain/dto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type MQTTClient struct {
	mqttClient     mqtt.Client
	log            logr.Logger
	commandFactory *domain_commands.CommandFactory

	config *config.ContainerConfig
}

func NewMQTTClient(commandFactory *domain_commands.CommandFactory, config *config.ContainerConfig, ctx context.Context) *MQTTClient {
	log := log.FromContext(ctx)
	log.Info("Init MQTT client...")
	opts := mqtt.NewClientOptions().AddBroker(config.BrokerUrl)

	log.Info("Generating client id...")
	rand.Seed(time.Now().UnixNano())
	opts.SetClientID(fmt.Sprintf("opc-ua-agent-%d", rand.Intn(10000)))
	log.Info("Client id generated and sotted")

	mqttClient := mqtt.NewClient(opts)

	return &MQTTClient{
		mqttClient:     mqttClient,
		log:            log,
		commandFactory: commandFactory,
		config:         config,
	}
}

func (c *MQTTClient) Connect() {
	c.log.Info("Connecting MQTT client...")
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		c.log.Error(token.Error(), "Failed to connect MQTT client")
		panic(token.Error())
	}

	c.log.Info("Subscribing MQTT client...")
	c.mqttClient.Subscribe(c.config.ConsumerTopic, 0, func(q mqtt.Client, m mqtt.Message) {
		var messageDTO dto.MessageDTO
		if err := json.Unmarshal(m.Payload(), &messageDTO); err != nil {
			c.log.Error(nil, "Failed to Unmarshal message. "+err.Error())
			return
		}

		c.log.Info(messageDTO.DeviceId + "==" + c.config.DeviceId)
		if messageDTO.DeviceId == c.config.DeviceId {
			for _, commandDTO := range messageDTO.Commands {
				c.log.Info("Command (" + commandDTO.CorrelationId + ") is for me..")
				if commandDTO.CorrelationId == "" {
					c.log.Error(nil, "The command not have a correlation id. Ignore.")
					continue
				}

				if command, err := c.commandFactory.Make(commandDTO); err == nil {
					data, cerr := command.Execute()
					if cerr != nil {
						c.PublishResult(command, false, "Failed to execute command."+cerr.Error(), data)
						continue
					}

					c.PublishResult(command, true, "", data)
				} else {
					c.PublishResult(command, false, "Command invalid. "+err.Error(), nil)
				}
			}
		}
	})
}

func (c *MQTTClient) PublishResult(cmd domain_commands.ICommand, Success bool, Message string, Data interface{}) error {
	c.log.Info(fmt.Sprintf("Publishing result: CorrelationId=%s, Success=%v, Message=%s", cmd.GetCorrelationId(), Success, Message))
	payload, err := json.MarshalIndent(dto.ResultDto{
		CorrelationId: cmd.GetCorrelationId(),
		Success:       Success,
		Data:          Data,
		Message:       Message,
		Timestamp:     time.Now().Unix(),
	},
		"", "  ")
	if err != nil {
		c.log.Error(err, "Failed to marshal result")
		return nil
	}
	if token := c.mqttClient.Publish(c.config.ResultsTopic, 0, false, payload); token.Wait() && token.Error() != nil {
		c.log.Error(token.Error(), "Failed to publish result")
		return nil
	}
	c.log.Info("Result published")
	return nil
}

func (c *MQTTClient) Disconnect() {
	c.log.Info("Disconnecting MQTT client...")
	c.mqttClient.Disconnect(250)
	c.log.Info("MQTT client disconnected")
}
