package opcua

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"opc.ua.agent/config"
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
	config *config.ContainerConfig
}

func NewOPCUAClient(device entities.Device, log logr.Logger, outputDriver domain_interfaces.IOutputDriver, config *config.ContainerConfig) *OPCUAClient {
	return &OPCUAClient{
		Device: device,
		config: config,

		output: outputDriver,
		log:    log,
	}
}

func (c *OPCUAClient) getCerts() (*rsa.PrivateKey, []byte, error) {
	return nil, nil, nil
}

func (c *OPCUAClient) Connect() error {
	var endpoint = "opc.tcp://" + c.Device.Url

	c.log.Info("Verify OPC UA endpoint: " + endpoint)
	eps, err := opcua.GetEndpoints(context.Background(), endpoint)
	if err != nil {
		return err
	}
	for _, ep := range eps {
		c.log.Info(ep.EndpointURL, ep.SecurityPolicyURI, ep.SecurityMode)
	}

	opts := []opcua.Option{
		opcua.SecurityMode(c.Device.SecurityMode),
	}
	if c.Device.SecurityMode >= 2 {
		TLSKeyPath := fmt.Sprintf("%s/opc-ua-agent.key", c.config.CertsDir)
		TLSCertPath := fmt.Sprintf("%s/opc-ua-agent.pem", c.config.CertsDir)
		c.log.Info("Loading certificates from " + TLSCertPath + " and " + TLSKeyPath)

		certObj, err := tls.LoadX509KeyPair(TLSCertPath, TLSKeyPath)
		if err != nil {
			return err
		} else {
			var ok bool
			pkey, ok := certObj.PrivateKey.(*rsa.PrivateKey)
			if !ok {
				c.log.Info("Invalid private key")
			}

			cert := certObj.Certificate[0]
			opts = append(opts,
				opcua.SecurityPolicy(ua.SecurityPolicyURIBasic256Sha256),
				opcua.PrivateKey(pkey),
				opcua.Certificate(cert),
			)
		}
	}

	c.log.Info("Connecting to OPC UA endpoint: " + c.Device.Url)
	if c.client, err = opcua.NewClient(endpoint, opts...); err != nil {
		return err
	}
	if err = c.client.Connect(context.Background()); err != nil {
		panic(err)
	}
	c.log.Info("OPC UA endpoint connected " + c.Device.Url)

	c.log.Info("Starting data polling for device: " + c.Device.DeviceId)
	c.ticker = time.NewTicker(time.Duration(c.Device.Interval) * time.Millisecond)
	go c.onData()
	c.log.Info("Data polling for device: " + c.Device.DeviceId + " started")

	return nil
}

func (c *OPCUAClient) Disconnect() error {
	c.ticker.Stop()
	return c.client.Close(context.Background())
}

func (c *OPCUAClient) PublishData(topic string, data interface{}) error {
	nodeID := ua.NewStringNodeID(uint16(c.Device.Ns), c.Device.Path+".opc_"+topic)
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

	c.log.Info("Writing OPC. Data: " + fmt.Sprintf("%v", data) + " key: " + nodeID.String())
	_, err = c.client.Write(context.Background(), req)
	if err != nil {
		return utils.EmitError(c.log, err.Error())
	}
	c.log.Info("Write OPC. Data: " + fmt.Sprintf("%v", data) + " key: " + nodeID.String())
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
				ReferenceTypeID: ua.NewNumericNodeID(0, 33),
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
		payload["timestamp"] = time.Now().UnixMicro()
		c.log.Info("Polling data for device: " + c.Device.DeviceId + " at " + t.String())
		nodeID := ua.NewStringNodeID(uint16(c.Device.Ns), c.Device.Path)

		c.log.Info(fmt.Sprintf("Getting all vars in NodeID=%s...", nodeID.String()))
		browseResp, err := c.getNodeList(nodeID)
		if err != nil {
			c.log.Error(err, "Erro no Browse:")
			return
		}

		c.log.Info(fmt.Sprintf("Found %d nodes", len(browseResp.Results)))
		for _, result := range browseResp.Results {
			c.log.Info(fmt.Sprintf("Found %d references", len(result.References)))
			for _, ref := range result.References {
				if strings.HasPrefix(ref.DisplayName.Text, "opc_") {
					c.log.Info(fmt.Sprintf("Getting value from NodeID=%s, BrowseName=%s, DisplayName=%s\n",
						ref.NodeID.String(),
						ref.BrowseName.Name,
						ref.DisplayName.Text,
					))
					if value, err := c.getNodeValue(ref.NodeID.NodeID); err == nil {
						attrKey := strings.TrimPrefix(ref.DisplayName.Text, "opc_")
						payload[attrKey] = value
					} else {
						c.log.Error(err, "Unable to getting value from NodeID=%s, BrowseName=%s, DisplayName=%s\n")
					}
				}
			}
		}

		if len(payload) > 1 {
			c.log.Info("Publishing data in mqtt broker...")
			if err := c.output.PublishData(c.Device.DeviceId, payload); err != nil {
				c.log.Error(err, "Failed to publish data in mqtt broker")
			}
			c.log.Info("Data published in mqtt broker")
		}
	}
}
