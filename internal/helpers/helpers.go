package helpers

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Target struct {
	ContractAddress    string `json:"contract-address"`
	ContractABI        string `json:"contract-abi"`
	ChainId            string `json:"chain-id"`
	NodeUrl            string `json:"node-url"`
	CreationBlock      uint64 `json:"creation-block"`
	LatestScrapedBlock uint64 `json:"latest-scraped-block"`
}

type Oracle struct {
	ContractAddress    common.Address
	ContractABI        *abi.ABI
	NodeUrl            string
	ChainID            string
	LatestScrapedBlock *big.Int
}

// Event emitted by the oracle contract
type OracleUpdate struct {
	Key       string
	Value     *big.Int
	Timestamp *big.Int
}

// Metadata on any transaction
type TransactionMetadata struct {
	BlockNumber     string
	BlockTimestamp  string
	TransactionFrom string
	TransactionTo   string
	TransactionHash string
	TransactionCost string
	SenderBalance   string
	GasUsed         string
	GasCost         string
	CreationBlock   string
}

// All the data scraped
type OracleMetrics struct {
	TransactionMetadata
	AssetKey        string
	AssetPrice      string
	UpdateTimestamp string
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
