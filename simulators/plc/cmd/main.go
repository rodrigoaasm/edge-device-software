package main

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"simulators.plc/config"
	"simulators.plc/internal/opcua"
	plc_runtime "simulators.plc/internal/runtime"
)

func main() {
	log := stdr.New(log.Default())
	ctx := logr.NewContext(context.Background(), log)
	defer ctx.Done()

	config := config.NewContainerConfig()

	opcua.Setup()
	log.Info("PLC Runtime for ed-operator tests...")
	selfClient := opcua.NewOPCUAClient(log, config.DeviceId, "opc.tcp://localhost:4841")
	if err := selfClient.Connect(); err != nil {
		log.Error(err, "Failed to connect to OPC UA")
		return
	}

	linClient := opcua.NewOPCUAClient(log, config.ParentId, config.ParentURL)
	if err := linClient.Connect(); err != nil {
		log.Error(err, "Failed to connect to OPC UA")
		return
	}

	plcRunner := plc_runtime.NewPLCRunner(config.CaptFrequecy)
	plcRunner.Run(selfClient, linClient, log)
}

// }
