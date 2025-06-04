package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-logr/logr"
	"opc.ua.agent/config"
	domain_commands "opc.ua.agent/internal/domain/commands"
	"opc.ua.agent/internal/domain/dto"
)

type MQTTClient struct {
	mqttClient     mqtt.Client
	log            logr.Logger
	commandFactory *domain_commands.CommandFactory

	config *config.ContainerConfig
}

func NewMQTTClient(commandFactory *domain_commands.CommandFactory, config *config.ContainerConfig, log logr.Logger) *MQTTClient {
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
		var messageDTO dto.CommandDTO
		if err := json.Unmarshal(m.Payload(), &messageDTO); err != nil {
			c.PublishResult(false, "Failed to Unmarshal message. "+err.Error(), nil)
			c.log.Error(err, "Failed to Unmarshal message")
			return
		}

		if messageDTO.DeviceId == c.config.DeviceId {
			c.log.Info("Command is for me..")
			command, err := c.commandFactory.Make(messageDTO)
			if err == nil {
				data, cerr := command.Execute()
				if cerr != nil {
					c.PublishResult(false, "Failed to execute command."+err.Error(), data)
					c.log.Error(cerr, "Failed to execute command")
					return
				}

				c.log.Info("Command executed.")
				c.PublishResult(true, "", data)
				return
			}

			c.log.Error(err, "Failed to make command")
		}
	})
}

func (c *MQTTClient) PublishResult(Success bool, Message string, Data interface{}) error {
	if Success {
		c.log.Info("Command executed successfully")
	} else {
		c.log.Error(nil, "Command execution failed. Message: "+Message)
	}

	c.log.Info("Publishing result...")
	payload, err := json.MarshalIndent(dto.ResultDto{Success: Success, Data: Data, Message: Message}, "", "  ")
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
