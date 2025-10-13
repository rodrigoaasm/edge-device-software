package opcua

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type OPCUAClient struct {
	log    logr.Logger
	client *opcua.Client
	url    string
	id     string
}

func NewOPCUAClient(log logr.Logger, id string, url string) *OPCUAClient {
	return &OPCUAClient{
		log: log,
		url: url,
		id:  id,
	}
}

func (c *OPCUAClient) Connect() error {
	c.log.Info("Verify OPC UA local endpoint: " + c.url)
	eps, err := opcua.GetEndpoints(context.Background(), c.url)
	if err != nil {
		return err
	}
	for _, ep := range eps {
		c.log.Info(ep.EndpointURL, ep.SecurityPolicyURI, ep.SecurityMode)
	}

	c.log.Info("Connecting to OPC UA local endpoint: ")
	if c.client, err = opcua.NewClient(c.url, opcua.SecurityMode(ua.MessageSecurityModeNone)); err != nil {
		return err
	}
	if err = c.client.Connect(context.Background()); err != nil {
		return err
	}
	c.log.Info("OPC UA endpoint connected")

	return nil
}

func (c *OPCUAClient) ReadVar(opcNs uint16, path string, vary string) (interface{}, error) {
	nodeID := ua.NewStringNodeID(opcNs, fmt.Sprintf("%s.%s.%s", path, c.id, vary))

	resp, err := c.client.Read(context.Background(), &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDValue,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Results[0].Value.Value(), nil
}

func (c *OPCUAClient) Close() {
	c.client.Close(context.Background())
}

func (c *OPCUAClient) WriteVar(opcNs uint16, path string, vars map[string]interface{}) error {
	for key, val := range vars {
		nodeID := ua.NewStringNodeID(opcNs, fmt.Sprintf("%s%s", path, key))
		vData, err := ua.NewVariant(val.(float64))
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

		c.log.Info("Writing OPC: ", nodeID.StringID(), val) //nodeID.StringID()
		_, err = c.client.Write(context.Background(), req)
		if err != nil {
			return err
		}
		c.log.Info("Wrote OPC.")
	}

	return nil
}
