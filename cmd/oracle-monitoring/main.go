package main

import (
	"log"

	"github.com/diadata-org/oracle-monitoring/internal/config"
	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/evm"
	"github.com/diadata-org/oracle-monitoring/internal/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func GetOracleCreationTimestamp(db database.DB, evmClient *evm.EVMClient, address string) (uint64, error) {
	event, err := models.GetLatestEventByOracleAddress(db, address)

	if err != nil {
		oracleCreationTimestamp := evmClient.CbFinder.GetContractCreationBlock(address)
		return oracleCreationTimestamp, nil
	}
	return event.OracleCreationTimestamp, nil
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close()

	// Initialize EVM client
	evmClient, err := evm.New(&cfg.EVM)
	if err != nil {
		log.Fatalf("Failed to initialize EVM client: %v", err)
	}

	// Run the oracle-monitoring service
	err = evmClient.MonitorEvents(func(rawEvent *types.Log) {
		log.Println("callback called")
		parsedEvent, err := models.ParseEvent(evmClient, rawEvent)

		oracleCreationTimestamp, _ := GetOracleCreationTimestamp(db, evmClient, parsedEvent.OracleAddress)
		parsedEvent.OracleCreationTimestamp = oracleCreationTimestamp

		if err != nil {
			log.Fatalf("Failed to parse event: %v", err)
		}

		err = models.StoreEvent(db, parsedEvent)
		if err != nil {
			log.Fatalf("Failed to store event: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to run oracle-monitoring service: %v", err)
	}
}
