package config

import (
	"log"
	"os"
	"strconv"

	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/evm"
	"github.com/joho/godotenv"
)

type Config struct {
	Database   database.Config
	EVM        evm.Config
	PrivateKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Read and parse environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	rpcNode := os.Getenv("RPC_NODE")
	oracleAddress := os.Getenv("ORACLE_ADDRESS")
	chainID, _ := strconv.ParseInt(os.Getenv("CHAIN_ID"), 10, 64)

	privateKey := os.Getenv("PRIVATE_KEY")

	// Create the configuration struct
	config := &Config{
		Database: database.Config{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
		},
		EVM: evm.Config{
			RPCNode:       rpcNode,
			OracleAddress: oracleAddress,
			ChainID:       chainID,
		},
		PrivateKey: privateKey,
	}

	return config, nil
}
