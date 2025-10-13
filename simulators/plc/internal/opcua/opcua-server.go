package opcua

import (
	"context"
	"log"

	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/server/attrs"
	"github.com/gopcua/opcua/ua"
)

func Setup() {
	log.Println("Iniciando o Servidor OPC-UA (gopcua v0.8.0 style)...")
	srv := server.New(
		server.EndPoint("0.0.0.0", 4841),
		server.EndPoint("172.31.212.236", 4841),
		server.EnableSecurity(ua.SecurityPolicyURINone, ua.MessageSecurityModeNone),
	)

	root_ns, _ := srv.Namespace(0)
	root_obj_node := root_ns.Objects()

	nodeNS := server.NewNodeNameSpace(srv, "NodeVars")
	nns_obj := nodeNS.Objects()
	srv.AddNamespace(nodeNS)
	root_obj_node.AddRef(nns_obj, id.HasComponent, true)

	vars := []string{"screwPosition", "injectPressure", "pressure", "temp"}
	for _, varName := range vars {
		nodeID := ua.NewStringNodeID(nodeNS.ID(), varName)
		node := server.NewNode(
			nodeID,
			map[ua.AttributeID]*ua.DataValue{
				ua.AttributeIDBrowseName: server.DataValueFromValue(attrs.BrowseName(varName)),
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
