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

var allOracles []string

func main() {
	var wg sync.WaitGroup

	db := database.NewPostgresDB()

	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	rpcmap, err := db.GetRPCByChainID([]string{})
	if err != nil {
		log.Printf("failed to get rpcmap: %v", err)
		return
	}

	wsurlmap, err := db.GetWSByChainID([]string{})
	if err != nil {
		log.Printf("failed to get rpcmap: %v", err)
		return
	}

	log.Println("starting historical")
	for chainid := range rpcmap {
		if chainid == "123420111" {
			//skip
		} else {
			// go runScraper(db, chainid, true, rpcmap, wsurlmap)
			go runEventScraper(db, chainid, false, rpcmap, wsurlmap)
		}

	}

	wg.Add(1)
	wg.Wait()
}

func runEventScraper(db database.Database, chainID string, isHistorical bool, rpcmap, wsurlmap map[string]string) {
	var wg sync.WaitGroup
	metricsChan := make(chan helpers.OracleMetrics)
	updateEventChan := make(chan helpers.OracleUpdateEvent)

	oracles, err := getOracles(db, chainID)

	if err != nil {
		log.Fatalf("failed to get oracles: %v", err)
	}

	if len(oracles) > 0 {

		fmt.Printf("\n Event based Scrapping started for chain %s,  total oracles %d isHistorical %t", chainID, len(oracles), isHistorical)

		ctx := context.Background()
		sc, err := scraper.NewScraper(ctx, metricsChan, updateEventChan, rpcmap, wsurlmap, big.NewInt(0), big.NewInt(0), oracles, &wg, chainID)
		if err != nil {
			return
		}
		oraclesArray := []common.Address{}

		var latest time.Time

		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {

				oracles, err := getOraclesByCreationTime(db, chainID, latest)
				if err != nil {
					fmt.Println("getOraclesByCreationTime err", err)
					continue
				}

				fmt.Println("oracles added", len(oracles))

				for _, oracle := range oracles {

					if latest.Before(oracle.CreatedDate) {
						latest = oracle.CreatedDate
					}

					oraclesArray = append(oraclesArray, oracle.ContractAddress)

				}

				sc.UpdateEvents(oraclesArray)
				sc.UpdateDeployedDate(oracles)

			}

		}()

		go processMetrics(db, metricsChan)
		go processCreation(db, updateEventChan)
	}

}

func runScraper(db database.Database, chainID string, isHistorical bool, rpcmap, wsurlmap map[string]string) {
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
		sc, err := scraper.NewScraper(ctx, metricsChan, updateEventChan, rpcmap, wsurlmap, minimum, maximum, oracles, &wg, chainID)
		if err != nil {
			return
		}
		if isHistorical {
			go sc.UpdateHistorical()

		} else {
			sc.UpdateRecent()

		}
		go processMetrics(db, metricsChan)
		go processCreation(db, updateEventChan)
	}

}

func getOraclesByCreationTime(db database.Database, chainID string, createdtime time.Time) (oracles []helpers.Oracle, err error) {

	oracleConfigs, err := db.SelectOraclesWithCreationTime(chainID, createdtime)
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
		oracle.CreatedDate = oracleconfig.CreatedDate
		oracle.ChainID = oracleconfig.ChainId
		oracle.ContractAddress = common.HexToAddress(oracleconfig.ContractAddress)
		oracle.LatestScrapedBlock = new(big.Int).SetUint64(oracleconfig.LatestScrapedBlock)
		oracles = append(oracles, oracle)
	}
	return oracles, nil

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
