package opcua

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
	"opc.ua.agent/internal/utils"
)

var (
	OPC_PATH        = "|var|CODESYS Control Win V3 x64.Application.PLC_PRG"
	OPC_NS   uint16 = 4
)

type OPCUAClient struct {
	Device entities.Device

	output domain_interfaces.IOutputDriver
	ticker *time.Ticker
	log    logr.Logger
	client *opcua.Client
}

func NewOPCUAClient(device entities.Device, log logr.Logger, outputDriver domain_interfaces.IOutputDriver) *OPCUAClient {
	return &OPCUAClient{
		Device: device,

		output: outputDriver,
		log:    log,
	}
}

func (c *OPCUAClient) Connect() error {
	var endpoint = flag.String("endpoint", "opc.tcp://"+c.Device.Url, "OPC UA Endpoint URL")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()

	c.log.Info("Verify OPC UA endpoint: " + *endpoint)
	eps, err := opcua.GetEndpoints(context.Background(), *endpoint)
	if err != nil {
		return err
	}
	for _, ep := range eps {
		c.log.Info(ep.EndpointURL, ep.SecurityPolicyURI, ep.SecurityMode)
	}

	c.log.Info("Connecting to OPC UA endpoint: " + c.Device.Url)
	if c.client, err = opcua.NewClient(*endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone)); err != nil {
		return err
	}
	if err = c.client.Connect(context.Background()); err != nil {
		return err
	}
	c.log.Info("OPC UA endpoint connected " + c.Device.Url)

	c.log.Info("Starting data polling for device: " + c.Device.DeviceId)
	c.ticker = time.NewTicker(time.Duration(c.Device.IntervalSeconds) * time.Second)
	go c.onData()
	c.log.Info("Data polling for device: " + c.Device.DeviceId + " started")

	return nil
}

func (c *OPCUAClient) Disconnect() error {
	c.ticker.Stop()
	return c.client.Close(context.Background())
}

func (c *OPCUAClient) PublishData(topic string, data interface{}) error {

	nodeID := ua.NewStringNodeID(4, OPC_PATH+"."+c.Device.DeviceId+"_"+topic)
	vData, err := ua.NewVariant(data.(bool))
	if err != nil {
		return err
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        vData,
				},
			},
		},
	}

	c.log.Info("Writing OPC. Data: " + fmt.Sprintf("%v", data) + " key: " + c.Device.DeviceId + "_" + topic)
	_, err = c.client.Write(context.Background(), req)
	if err != nil {
		return utils.EmitError(c.log, err.Error())
	}

	c.log.Info("Write OPC. Data: " + fmt.Sprintf("%v", data) + " key: " + topic)
	return nil
}

func (c *OPCUAClient) GetClientId() string {
	return c.Device.DeviceId
}

func (c *OPCUAClient) getNodeList(nodeId *ua.NodeID) (*ua.BrowseResponse, error) {
	browseReq := &ua.BrowseRequest{
		View: &ua.ViewDescription{},
		NodesToBrowse: []*ua.BrowseDescription{
			{
				NodeID:          nodeId,
				BrowseDirection: ua.BrowseDirectionForward,
				IncludeSubtypes: true,
				NodeClassMask:   uint32(ua.NodeClassAll),
				ResultMask:      uint32(ua.BrowseResultMaskAll),
			},
		},
	}

	return c.client.Browse(context.Background(), browseReq)
}

func (c *OPCUAClient) getNodeValue(nodeId *ua.NodeID) (interface{}, error) {
	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{
				NodeID:      nodeId,
				AttributeID: ua.AttributeIDValue,
			},
		},
	}

	resp, err := c.client.Read(context.Background(), req)
	if err != nil {
		return nil, utils.EmitError(c.log, "Value not read "+"from "+nodeId.String())
	} else {
		if len(resp.Results) > 0 && resp.Results[0].Status == ua.StatusOK {
			return resp.Results[0].Value.Value(), nil
		} else {
			return nil, utils.EmitError(c.log, "Value not read "+"from "+nodeId.String())
		}
	}
}

func (c *OPCUAClient) onData() {
	for t := range c.ticker.C {
		var payload map[string]interface{} = make(map[string]interface{})
		payload["timestamp"] = time.Now().Unix()
		c.log.Info("Polling data for device: " + c.Device.DeviceId + " at " + t.String())
		nodeID := ua.NewStringNodeID(OPC_NS, OPC_PATH)

		c.log.Info(fmt.Sprintf("Getting all vars in NodeID=%s...", nodeID.String()))
		browseResp, err := c.getNodeList(nodeID)
		if err != nil {
			c.log.Error(err, "Erro no Browse:")
		}

		c.log.Info(fmt.Sprintf("Found %d nodes", len(browseResp.Results)))
		for _, result := range browseResp.Results {
			c.log.Info(fmt.Sprintf("Found %d references", len(result.References)))
			for _, ref := range result.References {
				if strings.HasPrefix(ref.DisplayName.Text, c.Device.DeviceId) {
					c.log.Info(fmt.Sprintf("Getting value from NodeID=%s, BrowseName=%s, DisplayName=%s\n",
						ref.NodeID.String(),
						ref.BrowseName.Name,
						ref.DisplayName.Text,
					))
					if value, err := c.getNodeValue(ref.NodeID.NodeID); err == nil {
						attrKey := strings.TrimPrefix(ref.DisplayName.Text, c.Device.DeviceId+"_")
						payload[attrKey] = value
					} else {
						c.log.Error(err, "Unable to getting value from NodeID=%s, BrowseName=%s, DisplayName=%s\n")
					}
				}
			}
		}

		c.log.Info("Publishing data in mqtt broker...")
		if err := c.output.PublishData(c.Device.DeviceId, payload); err != nil {
			c.log.Error(err, "Failed to publish data in mqtt broker")
		}
		c.log.Info("Data published in mqtt broker")
	}
}
