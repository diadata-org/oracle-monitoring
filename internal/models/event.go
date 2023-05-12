package models

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

type Event struct {
	Id                      uint64
	OracleCreationTimestamp uint64
	OracleAddress           string
	ChainID                 int64
	TxHash                  string
	BlockNumber             uint64
	BlockTimestamp          uint64
	Asset                   string
	AssetValue              string
	DataTimestamp           uint64
	GasUsed                 uint64
	GasPrice                string
	TotalCost               string
	RemainingGasFunds       string
	TxSender                string
}

func (e *Event) FromLog(log *types.Log) error {
	// Parse log data and populate the Event fields.
	// You will need to implement this based on your specific Oracle contract event structure.

	return nil
}

func (e *Event) BlockDateTime() time.Time {
	return time.Unix(int64(e.BlockTimestamp), 0)
}

func (e *Event) DataDateTime() time.Time {
	return time.Unix(int64(e.DataTimestamp), 0)
}

func calcFrom(tx *types.Transaction) common.Address {
	msg, err := core.TransactionToMessage(tx, types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		panic(err)
	}

	return msg.From
}

func ParseEvent(evmClient *evm.EVMClient, vLog *types.Log) (*Event, error) {
	client := evmClient.Client
	chainID := evmClient.Config.ChainID

	// Parse the log to get event data
	// You need to implement this function based on your specific Oracle contract event structure

	// Retrieve additional data for the event
	tx, _, err := client.TransactionByHash(context.Background(), vLog.TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), vLog.TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction receipt: %w", err)
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(0).SetUint64(vLog.BlockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block: %w", err)
	}

	data := vLog.Data
	if err != nil {
		return nil, fmt.Errorf("failed to obtain oracle creation timestamp: %w", err)
	}

	blockTime := block.Time()
	fundingWallet := calcFrom(tx).String()
	fundingWalletBalance, _ := evmClient.GetAddressBalance(fundingWallet)
	// Parse the log to get event data
	event := &Event{
		OracleAddress:  vLog.Address.String(),
		ChainID:        chainID,
		TxHash:         vLog.TxHash.String(),
		BlockNumber:    vLog.BlockNumber,
		BlockTimestamp: blockTime,
		AssetValue:     big.NewInt(0).SetBytes(data[32:64]).String(),
		DataTimestamp:  big.NewInt(0).SetBytes(data[64:96]).Uint64(),
		Asset: strings.TrimRightFunc(string(data[128:160]), func(r rune) bool {
			return r == rune(0)
		}),
		GasUsed:           receipt.GasUsed,
		GasPrice:          tx.GasPrice().String(),
		TotalCost:         new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed)).String(),
		RemainingGasFunds: fundingWalletBalance.String(),
		TxSender:          fundingWallet,
	}

	return event, nil
}

func GetLatestEventByOracleAddress(db database.DB, oracleAddress string) (*Event, error) {
	row := db.Pool.QueryRow(context.Background(), `
		SELECT 
			id,
			oracle_address,
			chain_id,
			tx_hash,
			block_number,
			EXTRACT(epoch FROM block_timestamp)::bigint as block_timestamp,
			asset,
			asset_value,
			data_timestamp,
			gas_used,
			gas_price,
			total_cost,
			remaining_gas_funds,
			tx_sender,
			EXTRACT(epoch FROM oracle_creation_timestamp)::bigint as oracle_creation_timestamp 
		FROM 
			events 
		WHERE 
			oracle_address = $1 
		ORDER BY 
			block_number DESC LIMIT 1
	`, oracleAddress)

	var event Event

	err := row.Scan(
		&event.Id,
		&event.OracleAddress,
		&event.ChainID,
		&event.TxHash,
		&event.BlockNumber,
		&event.BlockTimestamp,
		&event.Asset,
		&event.AssetValue,
		&event.DataTimestamp,
		&event.GasUsed,
		&event.GasPrice,
		&event.TotalCost,
		&event.RemainingGasFunds,
		&event.TxSender,
		&event.OracleCreationTimestamp,
	)

	return &event, err
}

func StoreEvent(db database.DB, event *Event) error {
	query := `
        INSERT INTO events (
            oracle_address,
            chain_id,
            tx_hash,
            block_number,
            block_timestamp,
            asset,
            asset_value,
            data_timestamp,
            gas_used,
            gas_price,
            total_cost,
            remaining_gas_funds,
            tx_sender,
			oracle_creation_timestamp
        ) VALUES (
            $1, $2, $3, $4, to_timestamp($5), $6, $7, $8, $9, $10, $11, $12, $13, to_timestamp($14)
        )
    `
	_, err := db.Pool.Exec(
		context.Background(),
		query,
		event.OracleAddress,
		event.ChainID,
		event.TxHash,
		event.BlockNumber,
		event.BlockTimestamp,
		event.Asset,
		event.AssetValue,
		event.DataTimestamp,
		event.GasUsed,
		event.GasPrice,
		event.TotalCost,
		event.RemainingGasFunds,
		event.TxSender,
		event.OracleCreationTimestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}
