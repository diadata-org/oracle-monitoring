package main

import (
    "flag"
    "log"
    "oracle-monitoring/cmd"
)

func main() {
    // Parse the command-line arguments
    configFilePath := flag.String("config", "oracles.json", "Path to the configuration file")
    flag.Parse()

    // Read the configuration data from the file
    targets, err := config.LoadConfig(configFilePath)
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    // connect to db
    db, err := database.ConnectToDatabase()
    if err != nil {
        return err
    }
    defer db.Close()

    // update the list of targets
    err = database.InsertOracles(db, targets)
    if err != nil {
        return err
    }

    // fetch the full list of oracles with the known metadata
    oracles, err = database.SelectOracles(db)

    // add the ABI
    oracleABI, err := helpers.LoadContractAbi("internal/oracle-v2.abi.json")
    if err != nil {
        return err
    }

    for _, oracle := range oracles {
        oracle.ContractABI = oracleABI
    }

    // scrape the latest transactions for all the oracles
    err = scraper.Update(url, oracles, db)
    if err != nil {
        return err
    }
}
