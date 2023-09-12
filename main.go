package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/diadata-org/oracle-monitoring/internal/config"
	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/helpers"
	"github.com/diadata-org/oracle-monitoring/internal/scraper"
)

func main() {
	db := database.NewPostgresDB()

	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	rpcmap, err := db.GetRPCByChainID()
	if err != nil {
		log.Printf("failed to get rpcmap: %v", err)
		return
	}

	log.Println("starting historical")
	for chainid, _ := range rpcmap {
		go runScraper(db, chainid, true, rpcmap)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for chainid, _ := range rpcmap {
			go runScraper(db, chainid, false, rpcmap)
		}
	}
}

func runScraper(db database.Database, chainID string, isHistorical bool, rpcmap map[string]string) {
	var wg sync.WaitGroup
	metricsChan := make(chan helpers.OracleMetrics)
	updateEventChan := make(chan helpers.OracleUpdateEvent)

	oracles, err := getOracles(db, chainID)
	if err != nil {
		log.Fatalf("failed to get oracles: %v", err)
	}

	if len(oracles) > 0 {
		minimum, maximum := calculateMinMaxBlocks(oracles)

		fmt.Printf("\n Scrapping started for chain %s, up to minimum block %s, maximum block %s and total oracles %d isHistorical %t", chainID, minimum, maximum, len(oracles), isHistorical)

		ctx := context.Background()
		sc := scraper.NewScraper(ctx, metricsChan, updateEventChan, rpcmap, minimum, maximum, oracles, &wg, chainID)
		if isHistorical {
			go sc.UpdateHistorical()

		} else {
			sc.UpdateRecent()
		}
		go processMetrics(db, metricsChan)
		go processCreation(db, updateEventChan)
	}

}

func getOracles(db database.Database, chainID string) (oracles []helpers.Oracle, err error) {
	oracleConfigs, err := db.SelectOracles(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the saved metadata from the DB: %v", err)
	}

	for _, oracleconfig := range oracleConfigs {
		var oracle helpers.Oracle
		var oracleABI *abi.ABI
		oracleABI, err := config.LoadContractAbi(fmt.Sprintf("internal/abi/%s.json", "oracle-v2"))
		if err != nil {
			return nil, fmt.Errorf("failed to load the oracle contract ABI: %v", err)
		}
		oracle.NodeUrl = oracleconfig.NodeUrl
		oracle.ContractABI = oracleABI
		oracle.ChainID = oracleconfig.ChainId
		oracle.ContractAddress = common.HexToAddress(oracleconfig.ContractAddress)
		oracle.LatestScrapedBlock = new(big.Int).SetUint64(oracleconfig.LatestScrapedBlock)
		oracles = append(oracles, oracle)
	}

	return oracles, nil
}

func calculateMinMaxBlocks(oracles []helpers.Oracle) (*big.Int, *big.Int) {
	minimum := new(big.Int).Set(oracles[0].LatestScrapedBlock)
	maximum := new(big.Int).Set(oracles[0].LatestScrapedBlock)

	for _, oracle := range oracles {
		if minimum.Cmp(oracle.LatestScrapedBlock) == 1 {
			minimum.Set(oracle.LatestScrapedBlock)
		}

		if maximum.Cmp(oracle.LatestScrapedBlock) == -1 {
			maximum.Set(oracle.LatestScrapedBlock)
		}
	}

	return minimum, maximum
}

func processMetrics(db database.Database, metricsChan chan helpers.OracleMetrics) {
	for metrics := range metricsChan {
		if err := db.InsertOracleMetrics(&metrics); err != nil {
			log.Println("Error inserting oracle metrics:", err)
		}
	}
}

func processCreation(db database.Database, updateEvent chan helpers.OracleUpdateEvent) {
	for ue := range updateEvent {
		log.Println("updating oracle creation date", ue)
		if err := db.UpdateOracleCreation(ue.Address, ue.Block, ue.BlockTimestamp, ue.ChainID); err != nil {
			log.Println("Error inserting oracle creation:", err)
		}
	}
}
