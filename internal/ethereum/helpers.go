package ethereum

// loadABI loads the ABI of the oracle contract from a JSON file.
func loadABI(filePath string) (*abi.ABI, error) {
    abiBytes, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    abiObj, err := abi.JSON(bytes.NewReader(abiBytes))
    if err != nil {
        return nil, err
    }
    return &abiObj, nil
}