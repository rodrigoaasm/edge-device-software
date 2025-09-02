package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	stdr "github.com/go-logr/stdr"
	"opc.ua.agent/config"
	app_mqtt "opc.ua.agent/internal/app/mqtt"
	app_opcua "opc.ua.agent/internal/app/opcua"
	domain_commands "opc.ua.agent/internal/domain/commands"
	"opc.ua.agent/internal/domain/services"
	"opc.ua.agent/internal/persistence"
	"opc.ua.agent/internal/persistence/repositories"
)

func main() {
	// setup signal handler
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// load config and instance logger
	containerConfig := config.NewContainerConfig()
	log := stdr.New(log.Default())
	ctx := logr.NewContext(context.Background(), log)
	defer ctx.Done()

	// start database
	log.Info("Connecting to database...")
	db, dbErr := persistence.SqliteConnection()
	if dbErr != nil {
		log.Error(dbErr, "Failed to connect to database")
		panic(dbErr)
	}
	log.Info("Database connected. Run migrations...")
	persistence.Migrate(db)
	defer db.Close()

	// Init repositories
	deviceRepository := repositories.NewDeviceRepository(db)

	// services
	outputClientsManagerService := services.NewOutputClientsManagerService(log)

	// factory
	opcDriverFactory := app_opcua.NewOPCUAClientFactory(log)
	commandFactory := domain_commands.NewCommandFactory(
		ctx,
		opcDriverFactory,
		deviceRepository,
		outputClientsManagerService,
	)

	// mqtt client
	mqttClient := app_mqtt.NewMQTTClient(commandFactory, containerConfig, log)
	mqttClient.Connect()
	defer mqttClient.Disconnect()
	opcDriverFactory.OutputDriver = mqttClient

	<-signalChan
	log.Info("Shutting down...")
}
