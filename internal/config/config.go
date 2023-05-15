package config

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"

    "github.com/ethereum/go-ethereum/accounts/abi"

    "github.com/diadata-org/oracle-monitoring/internal/helpers"
)

// read the list of oracle targets
func LoadTargetConfig(filePath string) ([]helpers.Target, error) {
    var targets []helpers.Target

    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open the config file: %v", err)
    }

    err = json.Unmarshal(data, &targets)
    if err != nil {
        return nil, fmt.Errorf("failed to parse the JSON in the config file: %v", err)
    }

    return targets, nil
}

// loads the ABI of the oracle contract from a JSON file.
func LoadContractAbi(filePath string) (*abi.ABI, error) {
    abiBytes, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open the ABI file: %v", err)
    }

    abiObj, err := abi.JSON(bytes.NewReader(abiBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to parse the JSON for the ABI: %v", err)
    }

    return &abiObj, nil
}
