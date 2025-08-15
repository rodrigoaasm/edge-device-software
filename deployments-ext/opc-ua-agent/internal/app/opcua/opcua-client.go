package opcua

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
	"opc.ua.agent/internal/domain/entities"
)

type OPCUAClient struct {
	Device entities.Device

	log    logr.Logger
	client *opcua.Client
}

func NewOPCUAClient(device entities.Device, log logr.Logger) *OPCUAClient {
	return &OPCUAClient{
		Device: device,

		log: log,
	}
}

// const maxDepth = 10

// func browse(ctx context.Context, n *opcua.Node, path string, level int) ([], error) {
// 	// fmt.Printf("node:%s path:%q level:%d\n", n, path, level)
// 	if level > maxDepth {
// 		return nil, nil
// 	}

// 	attrs, err := n.Attributes(ctx, ua.AttributeIDNodeClass, ua.AttributeIDBrowseName, ua.AttributeIDDescription, ua.AttributeIDAccessLevel, ua.AttributeIDDataType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var def = NodeDef{
// 		NodeID: n.ID,
// 	}

// 	switch err := attrs[0].Status; err {
// 	case ua.StatusOK:
// 		def.NodeClass = ua.NodeClass(attrs[0].Value.Int())
// 	default:
// 		return nil, err
// 	}

// 	switch err := attrs[1].Status; err {
// 	case ua.StatusOK:
// 		def.BrowseName = attrs[1].Value.String()
// 	default:
// 		return nil, err
// 	}

// 	switch err := attrs[2].Status; err {
// 	case ua.StatusOK:
// 		def.Description = attrs[2].Value.String()
// 	case ua.StatusBadAttributeIDInvalid:
// 		// ignore
// 	default:
// 		return nil, err
// 	}

// 	switch err := attrs[3].Status; err {
// 	case ua.StatusOK:
// 		def.AccessLevel = ua.AccessLevelType(attrs[3].Value.Int())
// 		def.Writable = def.AccessLevel&ua.AccessLevelTypeCurrentWrite == ua.AccessLevelTypeCurrentWrite
// 	case ua.StatusBadAttributeIDInvalid:
// 		// ignore
// 	default:
// 		return nil, err
// 	}

// 	switch err := attrs[4].Status; err {
// 	case ua.StatusOK:
// 		switch v := attrs[4].Value.NodeID().IntID(); v {
// 		case id.DateTime:
// 			def.DataType = "time.Time"
// 		case id.Boolean:
// 			def.DataType = "bool"
// 		case id.SByte:
// 			def.DataType = "int8"
// 		case id.Int16:
// 			def.DataType = "int16"
// 		case id.Int32:
// 			def.DataType = "int32"
// 		case id.Byte:
// 			def.DataType = "byte"
// 		case id.UInt16:
// 			def.DataType = "uint16"
// 		case id.UInt32:
// 			def.DataType = "uint32"
// 		case id.UtcTime:
// 			def.DataType = "time.Time"
// 		case id.String:
// 			def.DataType = "string"
// 		case id.Float:
// 			def.DataType = "float32"
// 		case id.Double:
// 			def.DataType = "float64"
// 		default:
// 			def.DataType = attrs[4].Value.NodeID().String()
// 		}
// 	case ua.StatusBadAttributeIDInvalid:
// 		// ignore
// 	default:
// 		return nil, err
// 	}

// 	def.Path = join(path, def.BrowseName)
// 	// fmt.Printf("%d: def.Path:%s def.NodeClass:%s\n", level, def.Path, def.NodeClass)

// 	var nodes []NodeDef
// 	if def.NodeClass == ua.NodeClassVariable {
// 		nodes = append(nodes, def)
// 	}

// 	browseChildren := func(refType uint32) error {
// 		refs, err := n.ReferencedNodes(ctx, refType, ua.BrowseDirectionForward, ua.NodeClassAll, true)
// 		if err != nil {
// 			return errors.Errorf("References: %d: %s", refType, err)
// 		}
// 		// fmt.Printf("found %d child refs\n", len(refs))
// 		for _, rn := range refs {
// 			children, err := browse(ctx, rn, def.Path, level+1)
// 			if err != nil {
// 				return errors.Errorf("browse children: %s", err)
// 			}
// 			nodes = append(nodes, children...)
// 		}
// 		return nil
// 	}

// 	if err := browseChildren(id.HasComponent); err != nil {
// 		return nil, err
// 	}
// 	if err := browseChildren(id.Organizes); err != nil {
// 		return nil, err
// 	}
// 	if err := browseChildren(id.HasProperty); err != nil {
// 		return nil, err
// 	}
// 	return nodes, nil
// }

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
	if c.client, err = opcua.NewClient(c.Device.Url, opcua.SecurityMode(ua.MessageSecurityModeNone)); err != nil {
		return err
	}

	return nil
}

func (c *OPCUAClient) Disconnect() error {
	return c.client.Close(context.Background())
}

func (c *OPCUAClient) Publish(topic string, data interface{}) error {
	c.log.Info("Write OPC. Data: " + fmt.Sprintf("%v", data) + " key: " + topic)
	return nil
}

func (c *OPCUAClient) GetClientId() string {
	return c.Device.DeviceId
}
