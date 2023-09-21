package helpers

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Target struct {
	ContractAddress    string    `json:"contract-address"`
	ContractABI        string    `json:"contract-abi"`
	ChainId            string    `json:"chain-id"`
	NodeUrl            string    `json:"node-url"`
	CreationBlock      uint64    `json:"creation-block"`
	LatestScrapedBlock uint64    `json:"latest-scraped-block"`
	CreatedDate        time.Time `json:"createddate"`
}

type Oracle struct {
	ContractAddress    common.Address
	ContractABI        *abi.ABI
	NodeUrl            string
	ChainID            string
	LatestScrapedBlock *big.Int
	CreatedDate        time.Time
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
	ChainID         string
	BlockTimestamp  time.Time
	TransactionFrom common.Address
	TransactionTo   common.Address
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

type OracleMetricsState struct {
	ChainID   string
	LastBlock uint64
}

type OracleUpdateEvent struct {
	Address        string
	Block          string
	ChainID        string
	BlockTimestamp time.Time
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
