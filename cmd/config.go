package cmd

import (
    "encoding/json"
    "os"
)

type Oracle struct {
    Address string `json:"address"`
    Chain   string `json:"chain"`
    Url     string `json:"url"`
}

func LoadConfig(filePath string) ([]Oracle, error) {
    var oracles []Oracle

    data, err := os.ReadFile(filePath)
    if err == nil {
        err = json.Unmarshal(data, &oracles)
    }

    return oracles, err
}
