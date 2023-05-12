package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMClient struct {
	Client   *ethclient.Client
	Config   *Config
	CbFinder *ContractBirthFinder
}

type Config struct {
	RPCNode       string
	OracleAddress string
	ChainID       int64
}

func New(config *Config) (*EVMClient, error) {
	client, err := ethclient.Dial(config.RPCNode)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	cbFinder, err := NewContractBirthFinder(client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &EVMClient{
		Config:   config,
		Client:   client,
		CbFinder: cbFinder,
	}, nil
}

func (evmClient *EVMClient) MonitorEvents(callback func(result *types.Log)) error {
	oracleAddress := common.HexToAddress(evmClient.Config.OracleAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{oracleAddress},
	}

	// Continuously fetch logs
	logs := make(chan types.Log)
	sub, err := evmClient.Client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case vLog := <-logs:
			log.Println("received the event")
			callback(&vLog)
		}
	}
}

func (evmClient *EVMClient) GetAddressBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := evmClient.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (evmClient *EVMClient) SendTransaction(to string, privateKeyStr string, data []byte, gasLimit uint64) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error while loading private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	senderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	contractAddress := common.HexToAddress(to)
	if !ok {
		return nil, fmt.Errorf("error while getting public key: %w", err)
	}
	// Get the chain ID
	chainID, err := evmClient.Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Get the nonce for the sender
	nonce, err := evmClient.Client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce for sender: %w", err)
	}

	// Get the suggested gas price and maximum priority fee per gas from the client
	suggestedGasPrice, err := evmClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas price: %w", err)
	}
	suggestedGasTip, err := evmClient.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas tip: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get suggested max fee per gas: %w", err)
	}

	// Create the transaction object
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		Data:      data,
		To:        &contractAddress,
		Gas:       gasLimit,
		GasFeeCap: suggestedGasPrice,
		GasTipCap: suggestedGasTip,
	})

	// Sign the transaction with the private key
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return tx, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction to the client
	err = evmClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return tx, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}
