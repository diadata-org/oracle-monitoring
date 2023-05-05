package main

import (
	"context"
	"log"

	"github.com/diadata-org/oracle-monitoring/internal/config"
	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/evm"
	"github.com/diadata-org/oracle-monitoring/internal/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func GetOracleCreationTimestamp(db *database.DB, evmClient *evm.EVMClient, address string) (uint64, error) {
	row := db.Pool.QueryRow(context.Background(), `
		SELECT 
			id,
			oracle_address,
			chain_id,
			tx_hash,
			block_number,
			EXTRACT(epoch FROM block_timestamp)::bigint as block_timestamp,
			asset,
			asset_value,
			data_timestamp,
			gas_used,
			gas_price,
			total_cost,
			remaining_gas_funds,
			tx_sender,
			EXTRACT(epoch FROM oracle_creation_timestamp)::bigint as oracle_creation_timestamp 
		FROM 
			events 
		WHERE 
			oracle_address = $1 
		ORDER BY 
			block_number DESC LIMIT 1
	`, address)

	var event models.Event

	err := row.Scan(
		&event.Id,
		&event.OracleAddress,
		&event.ChainID,
		&event.TxHash,
		&event.BlockNumber,
		&event.BlockTimestamp,
		&event.Asset,
		&event.AssetValue,
		&event.DataTimestamp,
		&event.GasUsed,
		&event.GasPrice,
		&event.TotalCost,
		&event.RemainingGasFunds,
		&event.TxSender,
		&event.OracleCreationTimestamp,
	)

	if err != nil {
		oracleCreationTimestamp, _ := evmClient.IterateAndGetContractCreationTimestamp(address)
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
