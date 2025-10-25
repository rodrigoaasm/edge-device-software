package opcua

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/server/attrs"
	"github.com/gopcua/opcua/ua"
	"simulators.plc/config"
)

func Setup(config *config.ContainerConfig) {
	opts := []server.Option{
		server.EndPoint("0.0.0.0", config.OpcInternalServerPort),
		server.EndPoint(config.AlternativeDomain, config.OpcInternalServerPort),
	}

	if config.Tls {
		TLSKeyPath := fmt.Sprintf("/certs/%s.key", config.DeviceId)
		TLSCertPath := fmt.Sprintf("/certs/%s.pem", config.DeviceId)

		log.Println("Fazendo a leitura do certificado TLS: " + TLSKeyPath + " e " + TLSCertPath)

		certObj, err := tls.LoadX509KeyPair(TLSCertPath, TLSKeyPath)
		if err != nil {
			log.Fatal("Failed to load certificate:" + err.Error())
		} else {
			var ok bool
			pkey, ok := certObj.PrivateKey.(*rsa.PrivateKey)
			if !ok {
				log.Fatalf("Invalid private key")
			}
			cert := certObj.Certificate[0]
			opts = append(opts,
				server.EnableSecurity(ua.SecurityPolicyURIBasic256Sha256, ua.MessageSecurityModeSignAndEncrypt),
				server.EnableAuthMode(ua.UserTokenTypeAnonymous),
				server.PrivateKey(pkey),
				server.Certificate(cert),
			)
		}
		log.Println("Certificados foram lidos")
	} else {
		opts = append(opts, server.EnableSecurity(ua.SecurityPolicyURINone, ua.MessageSecurityModeNone))
	}

	log.Println("Iniciando o Servidor OPC-UA (gopcua v0.8.0 style)...")
	log.Println("Server opts.: ")
	for _, opt := range opts {
		log.Println(opt)
	}
	srv := server.New(opts...)

	root_ns, _ := srv.Namespace(0)
	root_obj_node := root_ns.Objects()

	nodeNS := server.NewNodeNameSpace(srv, "NodeVars")
	srv.AddNamespace(nodeNS)

	nns_obj := server.NewNode(
		ua.NewStringNodeID(nodeNS.ID(), "NodeVars"),
		map[ua.AttributeID]*ua.DataValue{
			ua.AttributeIDBrowseName: server.DataValueFromValue(attrs.BrowseName("NodeVars")),
			ua.AttributeIDNodeClass:  server.DataValueFromValue(uint32(ua.NodeClassObject)),
		},
		nil, nil,
	)
	nodeNS.AddNode(nns_obj)
	root_obj_node.AddRef(nns_obj, id.HasComponent, true)

	vars := []string{"screwPosition", "injectPressure", "pressure", "temp"}
	prefix := "opc_"
	for _, varName := range vars {
		nodeID := ua.NewStringNodeID(nodeNS.ID(), prefix+varName)
		node := server.NewNode(
			nodeID,
			map[ua.AttributeID]*ua.DataValue{
				ua.AttributeIDBrowseName: server.DataValueFromValue(attrs.BrowseName(prefix + varName)),
				ua.AttributeIDNodeClass:  server.DataValueFromValue(uint32(ua.NodeClassVariable)),
			},
			nil,
			func() *ua.DataValue { return server.DataValueFromValue(0.0) },
		)
		nodeNS.AddNode(node)
		nns_obj.AddRef(node, id.HasComponent, true)
	}

	// nodeNS := server.NewNodeNameSpace(srv, "NodeNamespace")
	// nns_obj := nodeNS.Objects()
	// root_obj_node.AddRef(nns_obj, id.HasComponent, true)

	if err := srv.Start(context.Background()); err != nil {
		log.Fatalf("Error starting server, exiting: %s", err)
	}
	log.Println("Servidor OPC-UA")

	// 4. Configurar e Iniciar o Servidor
	// srv := server.New(
	// 	server.WithEndpoint(endpointURL),
	// 	server.WithAddressSpace(nm.AddressSpace()), // Passa o Address Space do NodeManager
	// 	server.WithAddress(nm.Address()),
	// 	server.WithDisableSecurity(true), // Desabilita segurança
	// 	server.WithSessionTimeout(30*time.Minute),
	// 	server.WithApplicationURI("urn:gopcua:server:app"),
	// )
	// if err != nil {
	// 	log.Fatalf("Falha ao criar o servidor: %v", err)
	// }

	// // Inicia o servidor e bloqueia
	// log.Printf("Servidor pronto em %s", endpointURL)
	// log.Fatal(srv.Run())
}
