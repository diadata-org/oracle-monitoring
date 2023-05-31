package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/diadata-org/oracle-monitoring/internal/config"
	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/helpers"
	"github.com/diadata-org/oracle-monitoring/internal/scraper"
)

func main() {
	var oracles []helpers.Oracle
	var rpcmap map[string]string

	// Parse the command-line arguments
	// configFilePath := flag.String("targets", "oracles.json", "Path to the configuration file")
	// flag.Parse()

	// read the configuration data from the file
	// targets, err := config.LoadTargetConfig(*configFilePath)
	// if err != nil {
	// 	log.Fatalf("failed to read config file: %v", err)
	// }

	// connect to db

	db := database.NewPostgresDB()

	err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// update the list of targets
	// err = db.InsertOracles(targets)
	// if err != nil {
	// 	log.Fatalf("failed to update the list of target oracles: %v", err)
	// }

	// fetch the full list of oracles with the known metadata
	_oracles, err := db.SelectOracles()
	if err != nil {
		log.Fatalf("failed to get the saved metadata from the DB: %v", err)
	}

	// add the ABI
	for _, _oracle := range _oracles {
		var oracle helpers.Oracle
		var oracleABI *abi.ABI

		oracle.ContractAddress = common.HexToAddress(_oracle.ContractAddress)
		oracle.NodeUrl = _oracle.NodeUrl
		oracle.LatestScrapedBlock = new(big.Int).SetUint64(_oracle.LatestScrapedBlock)

		// read the oracle contract ABI from a file
		oracleABI, err = config.LoadContractAbi(fmt.Sprintf("internal/abi/%s.json", "oracle-v2"))
		if err != nil {
			log.Fatalf("failed to load the oracle contract ABI: %v", err)
		}

		oracle.ContractABI = oracleABI
		oracle.ChainID = _oracle.ChainId
		oracles = append(oracles, oracle)
	}

	rpcmap, err = db.GetRPCByChainID()

	if err != nil {
		log.Printf("failed to get rpcmap: %v", err)

	}

	sc := scraper.NewScraper(db, rpcmap)

	// scrape the latest transactions for all the oracles
	if err = sc.Update(oracles); err != nil {
		log.Printf("failed to scraped the latest oracle transactions: %v", err)
	}
}
