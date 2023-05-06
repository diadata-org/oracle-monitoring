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
    oracles, err := LoadConfig(configFilePath)
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    // Do something with the configuration data...
}
