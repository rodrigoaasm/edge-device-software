package mqtt

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
	"opc.ua.agent/internal/utils"
)

type CommandError struct {
	Error         error
	CorrelationId string
}

type MQTTClient struct {
	mqttClient     mqtt.Client
	log            logr.Logger
	commandFactory *domain_commands.CommandFactory

	config   *config.ContainerConfig
	clientId string
}

func NewMQTTClient(commandFactory *domain_commands.CommandFactory, config *config.ContainerConfig, log logr.Logger) *MQTTClient {
	log.Info("Init MQTT client...")
	opts := mqtt.NewClientOptions().AddBroker(config.BrokerUrl)

	log.Info("Generating client id...")
	rand.Seed(time.Now().UnixNano())
	clientId := fmt.Sprintf("opc-ua-agent-%d", rand.Intn(10000))
	opts.SetClientID(clientId)
	log.Info("Client id generated and sotted")

	mqttClient := mqtt.NewClient(opts)

	return &MQTTClient{
		mqttClient:     mqttClient,
		log:            log,
		commandFactory: commandFactory,
		config:         config,
		clientId:       clientId,
	}
}

func (c *MQTTClient) makeCmdError(correlationId string, message string) *CommandError {
	return &CommandError{
		Error:         utils.EmitError(c.log, message),
		CorrelationId: correlationId,
	}
}

func (c *MQTTClient) parseMessageToCommands(messageDTO map[string]interface{}) ([]domain_commands.ICommand, *CommandError) {
	commands := make([]domain_commands.ICommand, 0)
	if messageDTO["deviceId"] == c.config.DeviceId {
		if messageDTO["commands"] == nil {
			return nil, c.makeCmdError(messageDTO["correlationId"].(string), "The message not have commands. Ignore.")
		}

		for _, commandG := range messageDTO["commands"].([]interface{}) {
			commandDTO := commandG.(map[string]interface{})
			correlationId := commandDTO["correlationId"].(string)
			c.log.Info("Command (" + correlationId + ") is for me..")
			if correlationId == "" {
				c.log.Error(nil, "The command not have a correlation id. Ignore.")
				continue
			}

			if command, err := c.commandFactory.Make(commandDTO); err == nil {
				commands = append(commands, command)
			} else {
				return nil, c.makeCmdError(correlationId, "Command invalid. "+err.Error())
			}
		}
	} else {
		c.log.Info("Actuate will be passed on..")
		timestampf, ok := messageDTO["timestamp"].(float64)
		if !ok {
			return nil, c.makeCmdError(fmt.Sprintf("%d", timestampf), "The message not have a timestamp or it is not a epoch. Ignore.")
		}

		if messageDTO["deviceId"] == nil {
			return nil, c.makeCmdError(fmt.Sprintf("%d", timestampf), "The message not have a device id. Ignore.")
		}

		if command, err := c.commandFactory.MakeActuateCmd(
			int64(timestampf),
			messageDTO["deviceId"].(string),
			messageDTO,
		); err == nil {
			commands = append(commands, command)
		} else {
			return nil, c.makeCmdError(fmt.Sprintf("%d", timestampf), "Actuate invalid. "+err.Error())
		}
	}

	return commands, nil
}

func (c *MQTTClient) handlerMessage(q mqtt.Client, m mqtt.Message) {
	var messageDTO map[string]interface{}
	if err := json.Unmarshal(m.Payload(), &messageDTO); err != nil {
		c.log.Error(nil, "Failed to Unmarshal message. "+err.Error(), nil)
		return
	}

	commands, err := c.parseMessageToCommands(messageDTO)
	if err != nil {
		c.PublishResult(err.CorrelationId, false, err.Error.Error(), nil)
		return
	}

	for _, command := range commands {
		data, cerr := command.Execute()
		if cerr != nil {
			c.PublishResult(command.GetCorrelationId(), false, "Failed to execute command."+cerr.Error(), data)
			continue
		}

		c.PublishResult(command.GetCorrelationId(), true, "", data)
	}
}

func (c *MQTTClient) Connect() error {
	c.log.Info("Connecting MQTT client...")
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		c.log.Error(token.Error(), "Failed to connect MQTT client")
		panic(token.Error())
	}

	c.log.Info("Subscribing MQTT client...")
	c.mqttClient.Subscribe(c.config.ConsumerTopic, 0, c.handlerMessage)
	return nil
}

func (c *MQTTClient) PublishResult(CorrelationId string, Success bool, Message string, Data interface{}) error {
	c.log.Info(fmt.Sprintf("Publishing result: CorrelationId=%s, Success=%v, Message=%s", CorrelationId, Success, Message))
	payload, err := json.MarshalIndent(dto.ResultDto{
		CorrelationId: CorrelationId,
		Success:       Success,
		Data:          Data,
		Message:       Message,
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

func (c *MQTTClient) PublishData(deviceId string, message interface{}) error {
	c.log.Info("Publishing data in " + deviceId)
	payload, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		c.log.Error(err, "Failed to marshal result")
		return nil
	}
	if token := c.mqttClient.Publish(c.config.DataTopic+"/"+deviceId, 0, false, payload); token.Wait() && token.Error() != nil {
		c.log.Error(token.Error(), "Failed to publish result")
		return nil
	}
	c.log.Info("Data published in " + deviceId)
	return nil
}

func (c *MQTTClient) Disconnect() error {
	c.log.Info("Disconnecting MQTT client...")
	c.mqttClient.Disconnect(250)
	c.log.Info("MQTT client disconnected")
	return nil
}

func (c *MQTTClient) GetClientId() string {
	return c.clientId
}
