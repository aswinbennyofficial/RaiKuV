package main

import (
	"fmt"
	"os"

	"github.com/aswinbennyofficial/raikuv/configs"
	"github.com/aswinbennyofficial/raikuv/internal/logging"
	"github.com/aswinbennyofficial/raikuv/internal/node"
	"github.com/aswinbennyofficial/raikuv/internal/storage"
)

func main(){
	yamlConfigs, err:=configs.LoadConfig()
	if err!=nil{
		fmt.Printf("Error loading YAML configs: %v\n", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(yamlConfigs.LogConfig)

	// Initialise a new data store
	ds := storage.NewDataStore()

	// Initialise new TCP server for client to node comms
	tcpServer := node.NewTCPServer("3232",ds, &logger)
	err=tcpServer.ListenAndServe()
	if err !=nil {
		logger.Error().Err(err).Msg("Failed to listen to tcp")
		os.Exit(1)
	}

	
}