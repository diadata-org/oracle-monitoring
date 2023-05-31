package scraper

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/helpers"
)

// Scraper is an interface that represents the scraper functionality.
type Scraper interface {
	ScrapeSingleOracle(oracle *helpers.Oracle) error
	Update(oracles []helpers.Oracle) error
}

type scraperImpl struct {
	db    database.Database
	nodes map[string]*ethclient.Client
	rpc   map[string]string
}

// NewScraper creates a new instance of the Scraper interface.
func NewScraper(db database.Database, rpcmap map[string]string) Scraper {
	return &scraperImpl{
		db:    db,
		nodes: make(map[string]*ethclient.Client),
		rpc:   rpcmap,
	}
}

func (s *scraperImpl) connectToNode(url string) (*ethclient.Client, error) {
	// Check if the client for the given URL is already connected
	if client, ok := s.nodes[url]; ok {
		return client, nil
	}

	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the node: %v", err)
	}

	// Store the connected client for future use
	s.nodes[url] = client
	return client, nil
}

func (s *scraperImpl) getTransactionSender(tx *types.Transaction) (common.Address, error) {
	return types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
}

func (s *scraperImpl) isTargetingContract(tx *types.Transaction, address common.Address) bool {
	return tx.To() != nil && *tx.To() == address
}

func (s *scraperImpl) isContractCreation(tx *types.Transaction, receipt *types.Receipt, address common.Address) bool {
	if tx.To() == nil {
		// Transaction is a contract creation
		return receipt != nil && receipt.ContractAddress == address
	}

	return false
}

func (s *scraperImpl) isOracleUpdate(tx *types.Transaction, oracle *helpers.Oracle) bool {
	data := tx.Data()
	setValueSig := oracle.ContractABI.Methods["setValue"].ID

	return bytes.Equal(data[:4], setValueSig[:4])
}

func (s *scraperImpl) parseTransactionMetadata(client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt) (*helpers.TransactionMetadata, error) {
	metadata := &helpers.TransactionMetadata{}

	metadata.BlockNumber = block.Number().String()
	metadata.BlockTimestamp = strconv.FormatUint(block.Time(), 10)
	metadata.TransactionHash = strings.ToLower(tx.Hash().String())
	metadata.TransactionCost = new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice).String()
	metadata.GasUsed = strconv.FormatUint(receipt.GasUsed, 10)
	metadata.GasCost = receipt.EffectiveGasPrice.String()

	sender, err := s.getTransactionSender(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %v", err)
	}
	metadata.TransactionFrom = sender.String()

	if tx.To() != nil {
		metadata.TransactionTo = tx.To().String()
	}

	senderBalance, err := client.BalanceAt(context.Background(), sender, block.Number())
	if err != nil {
		return nil, fmt.Errorf("failed to get sender balance: %v", err)
	}
	metadata.SenderBalance = senderBalance.String()

	return metadata, nil
}

func (s *scraperImpl) parseOracleUpdate(tx *types.Transaction, receipt *types.Receipt, oracle *helpers.Oracle) (*helpers.OracleUpdate, error) {
	event := &helpers.OracleUpdate{}
	err := oracle.ContractABI.UnpackIntoInterface(event, "OracleUpdate", tx.Data()[4:])

	return event, err
}

func (s *scraperImpl) parseTransaction(client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt, oracle *helpers.Oracle) (bool, error) {
	done := false
	metadata, err := s.parseTransactionMetadata(client, block, tx, receipt)
	if err != nil {
		return false, fmt.Errorf("failed to parse transaction metadata: %v", err)
	}

	if s.isContractCreation(tx, receipt, oracle.ContractAddress) {
		// to is nil on contract creation
		metadata.TransactionTo = strings.ToLower(oracle.ContractAddress.String())

		// err := s.db.UpdateOracleCreation(oracle.ContractAddress.String(), block.Number().Uint64())
		// if err != nil {
		// 	return false, fmt.Errorf("failed to update oracle creation data: %v", err)
		// }

		// no more transactions to scrape after the oracle creation
		done = true
	}

	if s.isTargetingContract(tx, oracle.ContractAddress) && s.isOracleUpdate(tx, oracle) {
		oracleUpdate, err := s.parseOracleUpdate(tx, receipt, oracle)
		if err != nil {
			return false, fmt.Errorf("failed to parse oracle update: %v", err)
		}

		metrics := &helpers.OracleMetrics{
			TransactionMetadata: *metadata,
			AssetKey:            oracleUpdate.Key,
			AssetPrice:          oracleUpdate.Value.String(),
			UpdateTimestamp:     oracleUpdate.Timestamp.String(),
		}

		fmt.Println(helpers.PrettyPrint(metrics))

		err = s.db.InsertOracleMetrics(metrics)
		if err != nil {
			return false, fmt.Errorf("failed to insert oracle metrics: %v", err)
		}

		// reached the latest block scraped on the previous run
		done = block.Number().Cmp(oracle.LatestScrapedBlock) <= 0 // -1 or 0 if block.Number is smaller
	}

	return done, nil
}

func (s *scraperImpl) parseBlock(client *ethclient.Client, block *types.Block, oracle *helpers.Oracle) (bool, error) {
	done := false

	for _, tx := range block.Transactions() {
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return false, fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		creation, err := s.parseTransaction(client, block, tx, receipt, oracle)
		if err != nil {
			return false, fmt.Errorf("failed to scrape transaction: %v", err)
		}

		done = done || creation
	}

	return done, nil
}

func (s *scraperImpl) ScrapeSingleOracle(oracle *helpers.Oracle) error {
	done := false

	log.Println("oracle.ChainID", oracle.ChainID)

	client, err := s.connectToNode(s.rpc[oracle.ChainID])
	if err != nil {
		return fmt.Errorf("failed to connect to the node: %v", err)
	}

	current, err := client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve the latest block: %v", err)
	}

	for !done && current > 0 && current > oracle.LatestScrapedBlock.Uint64() {
		block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(current))
		if err != nil {
			return fmt.Errorf("failed to retrieve a block: %v, %v ", current, err)
		}

		fmt.Println(block.Number().Uint64())

		done, err = s.parseBlock(client, block, oracle)
		if err != nil {
			return fmt.Errorf("failed to scrape block: %v", err)
		}

		current = current - 1
	}

	return nil
}

func (s *scraperImpl) Update(oracles []helpers.Oracle) error {
	log.Println(len(oracles))
	for index, oracle := range oracles {
		log.Println(index)

		err := s.ScrapeSingleOracle(&oracle)
		if err != nil {
			fmt.Println("failed to scrape oracle: %v", err)
		}
	}

	return nil
}
