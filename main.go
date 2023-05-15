package main

import (
    "flag"
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

    // Parse the command-line arguments
    configFilePath := flag.String("targets", "oracles.json", "Path to the configuration file")
    flag.Parse()

    // read the configuration data from the file
    targets, err := config.LoadTargetConfig(*configFilePath)
    if err != nil {
        log.Fatalf("failed to read config file: %v", err)
    }

    // connect to db
    db, err := database.ConnectToDatabase()
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    defer db.Close()

    // update the list of targets
    err = database.InsertOracles(db, targets)
    if err != nil {
        log.Fatalf("failed to update the list of target oracles: %v", err)
    }

    // fetch the full list of oracles with the known metadata
    _oracles, err := database.SelectOracles(db)
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
        oracleABI, err = config.LoadContractAbi(fmt.Sprintf("internal/abi/%s.json", _oracle.ContractABI))
        if err != nil {
            log.Fatalf("failed to load the oracle contract ABI: %v", err)
        }

        oracle.ContractABI = oracleABI
        oracles = append(oracles, oracle)
    }

    // scrape the latest transactions for all the oracles
    if err = scraper.Update(oracles, db); err != nil {
        log.Fatalf("failed to scraped the latest oracle transactions: %v", err)
    }
}
