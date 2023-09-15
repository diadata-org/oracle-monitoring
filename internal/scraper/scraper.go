package scraper

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/diadata-org/oracle-monitoring/internal/helpers"
	"github.com/google/uuid"
)

// Scraper is an interface that represents the scraper functionality.
type Scraper interface {
	UpdateHistorical() error
	UpdateRecent() error
	UpdateEvents(oracleaddresses []common.Address) error
}

type scraperImpl struct {
	nodes            map[string]*ethclient.Client
	wsNodes          map[string]*ethclient.Client
	rpc              map[string]string
	wsurl            map[string]string
	mchan            chan helpers.OracleMetrics
	createChan       chan helpers.OracleUpdateEvent
	ctx              context.Context
	minblock         *big.Int
	maxblock         *big.Int
	oracles          []helpers.Oracle
	wg               *sync.WaitGroup
	chainID          string
	oraclesaddresses []common.Address
	oraclesmap       map[common.Address]helpers.Oracle
	client           *ethclient.Client
	wsClient         *ethclient.Client
	isHistorical     bool
	logger           *log.Logger
}

// NewScraper creates a new instance of the Scraper interface.
func NewScraper(context context.Context, mchan chan helpers.OracleMetrics, createChan chan helpers.OracleUpdateEvent, rpcmap, wsmap map[string]string, minblock *big.Int, maxblock *big.Int, oracles []helpers.Oracle, wg *sync.WaitGroup, chainID string) Scraper {

	id := uuid.Must(uuid.NewRandom()).String()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.SetPrefix(id)

	s := &scraperImpl{
		mchan:        mchan,
		nodes:        make(map[string]*ethclient.Client),
		wsNodes:      make(map[string]*ethclient.Client),
		rpc:          rpcmap,
		wsurl:        wsmap,
		ctx:          context,
		minblock:     minblock,
		maxblock:     maxblock,
		oracles:      oracles,
		wg:           wg,
		createChan:   createChan,
		chainID:      chainID,
		isHistorical: false,

		logger: logger,
	}

	s.oraclesmap = make(map[common.Address]helpers.Oracle)
	for _, oracle := range s.oracles {
		s.oraclesmap[oracle.ContractAddress] = oracle
		s.logger.Println("oracles", oracle.ContractAddress)

		s.oraclesaddresses = append(s.oraclesaddresses, oracle.ContractAddress)

	}
	var err error
	s.client, err = s.connectToNode()
	if err != nil {
		s.logger.Println("error connecting to rpc chainid ", chainID)
	}

	s.wsClient, err = s.connectToWsNode()
	if err != nil {
		s.logger.Println("error connecting to ws chainid  ", chainID)
		panic("")
	}
	return s

}

func (s *scraperImpl) connectToWsNode() (*ethclient.Client, error) {
	// Check if the client for the given URL is already connected

	url := s.wsurl[s.chainID]

	fmt.Println("url", url)

	if client, ok := s.wsNodes[url]; ok {
		return client, nil
	}

	client, err := ethclient.DialContext(s.ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the node: %v", err)
	}

	// Store the connected client for future use
	s.wsNodes[url] = client
	return client, nil
}

func (s *scraperImpl) connectToNode() (*ethclient.Client, error) {
	// Check if the client for the given URL is already connected

	url := s.rpc[s.chainID]
	if client, ok := s.nodes[url]; ok {
		return client, nil
	}

	client, err := ethclient.DialContext(s.ctx, url)
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

func (s *scraperImpl) isTargetingContract(tx *types.Transaction, addresss []common.Address) (common.Address, bool) {

	if tx.To() != nil && contains(addresss, *tx.To()) {
		return *tx.To(), true
	}
	return common.HexToAddress("0"), false
}

func (s *scraperImpl) isContractCreation(tx *types.Transaction, receipt *types.Receipt, addresses []common.Address) (common.Address, bool) {

	if tx.To() == nil {
		// Transaction is a contract creation

		if receipt != nil && contains(addresses, receipt.ContractAddress) {
			return receipt.ContractAddress, true
		}
	}

	return receipt.ContractAddress, false
}

func contains(slice []common.Address, item common.Address) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *scraperImpl) isOracleUpdate(tx *types.Transaction, oracle helpers.Oracle) bool {
	data := tx.Data()
	setValueSig := oracle.ContractABI.Methods["setValue"].ID

	return bytes.Equal(data[:4], setValueSig[:4])
}

func (s *scraperImpl) parseTransactionMetadata(ctx context.Context, client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt) (*helpers.TransactionMetadata, error) {
	metadata := &helpers.TransactionMetadata{}

	metadata.ChainID = s.chainID
	metadata.BlockNumber = block.Number().String()
	metadata.BlockTimestamp = time.Unix(int64(block.Time()), 0)
	metadata.TransactionHash = strings.ToLower(tx.Hash().String())
	metadata.TransactionCost = new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice).String()
	metadata.GasUsed = strconv.FormatUint(receipt.GasUsed, 10)
	metadata.GasCost = receipt.EffectiveGasPrice.String()

	sender, err := s.getTransactionSender(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %v", err)
	}

	metadata.TransactionFrom = sender

	if tx.To() != nil {
		metadata.TransactionTo = *tx.To()
	}

	// use lates balance instead of block number as that need archieve node
	senderBalance, err := client.BalanceAt(ctx, sender, nil)
	if err != nil {
		s.logger.Printf("failed to get sender balance: %v", err)
	}

	metadata.SenderBalance = senderBalance.String()

	return metadata, nil
}

func (s *scraperImpl) parseOracleUpdate(tx *types.Transaction, receipt *types.Receipt, oracle helpers.Oracle) (*helpers.OracleUpdate, error) {
	event := &helpers.OracleUpdate{}
	err := oracle.ContractABI.UnpackIntoInterface(event, "OracleUpdate", tx.Data()[4:])

	return event, err
}

func (s *scraperImpl) parseTransaction(ctx context.Context, client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt) (bool, error) {
	done := false
	metadata, err := s.parseTransactionMetadata(ctx, client, block, tx, receipt)
	if err != nil {
		return false, fmt.Errorf("failed to parse transaction metadata: %v", err)
	}

	contract, iscreated := s.isContractCreation(tx, receipt, s.oraclesaddresses)

	if iscreated {
		s.logger.Println("isContractCreation", contract.String())
		// to is nil on contract creation
		metadata.TransactionTo = contract

		ue := helpers.OracleUpdateEvent{}
		ue.Address = contract.String()
		ue.Block = metadata.BlockNumber
		ue.ChainID = s.chainID
		ue.BlockTimestamp = metadata.BlockTimestamp

		s.createChan <- ue

		// err := s.db.UpdateOracleCreation(oracle.ContractAddress.String(), block.Number().Uint64())
		// if err != nil {
		// 	return false, fmt.Errorf("failed to update oracle creation data: %v", err)
		// }

		// no more transactions to scrape after the oracle creation
		done = true
	}

	contract, istarget := s.isTargetingContract(tx, s.oraclesaddresses)

	if istarget {
		if s.isOracleUpdate(tx, s.oraclesmap[contract]) {

			oracleUpdate, err := s.parseOracleUpdate(tx, receipt, s.oraclesmap[contract])
			if err != nil {
				return false, fmt.Errorf("failed to parse oracle update: %v", err)
			}

			metrics := &helpers.OracleMetrics{
				TransactionMetadata: *metadata,
				AssetKey:            oracleUpdate.Key,
				AssetPrice:          oracleUpdate.Value.String(),
				UpdateTimestamp:     oracleUpdate.Timestamp.String(),
			}

			// err = s.db.InsertOracleMetrics(metrics)
			// if err != nil {
			// 	return false, fmt.Errorf("failed to insert oracle metrics: %v", err)
			// }

			// reached the latest block scraped on the previous run
			// done = block.Number().Cmp(oracle.LatestScrapedBlock) <= 0 // -1 or 0 if block.Number is smaller
			s.mchan <- *metrics

		}
	}
	return done, nil
}

func (s *scraperImpl) listenEvents(addresses []common.Address) {

	ctx := context.Background()

	updateeventchan := make(chan types.Log)

	if len(s.oraclesaddresses) <= 0 {
		return
	}

	topic := s.oraclesmap[s.oraclesaddresses[0]].ContractABI.Events["OracleUpdate"].ID

	subscription, err := s.wsClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    [][]common.Hash{{topic}},
	}, updateeventchan)
	if err != nil {
		log.Fatalf("Failed to subscribe to event logs: %v", err)
	}

	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			log.Println("Subscription error: %v", err)
		case eventLog := <-updateeventchan:

			eventData := make(map[string]interface{})

			err := s.oraclesmap[eventLog.Address].ContractABI.UnpackIntoMap(eventData, "OracleUpdate", eventLog.Data)
			if err != nil {
				log.Printf("ailed to unpack log data: %v", err)
			}

			receipt, err := s.client.TransactionReceipt(s.ctx, eventLog.TxHash)
			if err != nil {
				fmt.Printf("failed to get transaction receipt: %v", err)
			}

			block, err := s.client.BlockByNumber(s.ctx, new(big.Int).SetUint64(eventLog.BlockNumber))
			if err != nil {
				s.logger.Printf("failed to retrieve a block: %v, %v, chainid %s ", eventLog.BlockNumber, err, s.chainID)
				continue
			}

			metadata := helpers.TransactionMetadata{}
			metadata.BlockNumber = strconv.Itoa(int(eventLog.BlockNumber))
			// metadata.BlockTimestamp = eventLog.
			metadata.TransactionHash = eventLog.TxHash.Hex()
			metadata.TransactionCost = new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice).String()
			metadata.BlockTimestamp = time.Unix(int64(block.Time()), 0)
			metadata.GasUsed = strconv.FormatUint(receipt.GasUsed, 10)
			metadata.GasCost = receipt.EffectiveGasPrice.String()
			metadata.ChainID = s.chainID

			tx, _, err := s.client.TransactionByHash(context.Background(), eventLog.TxHash)
			if err != nil {
				fmt.Printf("failed to get sender: %v", err)
			}

			sender, err := s.getTransactionSender(tx)
			if err != nil {
				fmt.Printf("failed to get sender: %v", err)
			}

			metadata.TransactionFrom = sender

			if tx.To() != nil {
				metadata.TransactionTo = *tx.To()
			}

			// use lates balance instead of block number as that need archieve node
			senderBalance, err := s.client.BalanceAt(ctx, sender, nil)
			if err != nil {
				s.logger.Printf("failed to get sender balance: %v", err)
			}

			metadata.SenderBalance = senderBalance.String()

			metadata.TransactionTo = eventLog.Address

			metrics := &helpers.OracleMetrics{
				TransactionMetadata: metadata,
				AssetKey:            eventData["key"].(string),
				AssetPrice:          eventData["value"].(*big.Int).String(),
				UpdateTimestamp:     eventData["timestamp"].(*big.Int).String(),
			}

			s.mchan <- *metrics

		}
	}
}

func (s *scraperImpl) parseBlock(block *types.Block) (bool, error) {
	done := false

	s.logger.Printf(" parsing block  %s, for chain  %s", block.Number(), s.chainID)

	for _, tx := range block.Transactions() {

		receipt, err := s.client.TransactionReceipt(s.ctx, tx.Hash())
		if err != nil {
			return false, fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		creation, err := s.parseTransaction(s.ctx, s.client, block, tx, receipt)
		if err != nil {
			return false, fmt.Errorf("failed to scrape transaction: %v", err)
		}

		done = done || creation
	}

	s.logger.Printf("parsed block  %s, for chain  %s isHistorical %t", block.Number().String(), s.chainID, s.isHistorical)

	return done, nil
}

func (s *scraperImpl) recent(current uint64) (err error) {

	if current == 0 {
		current, err = s.client.BlockNumber(s.ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve the latest block: %v", err)
		}

	}

	for current > s.maxblock.Uint64() {
		block, err := s.client.BlockByNumber(s.ctx, new(big.Int).SetUint64(current))
		if err != nil {
			s.logger.Printf("failed to retrieve a block: %v, %v, chainid %s ", current, err, s.chainID)
			continue
		}

		_, err = s.parseBlock(block)
		if err != nil {
			s.logger.Printf("failed to scrape block: %v", err)
		}

		current = current - 1

	}

	s.logger.Printf("done  oracles  ")

	return nil
}

func (s *scraperImpl) historical() error {

	current, err := s.client.BlockNumber(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve the latest block: %v", err)
	}
	// switch s.chainID {
	// case "80001":
	// 	current = 34889087

	// case "5":
	// 	current = 8263129

	// }
	// change it to minimum

	for current > 0 && current > s.minblock.Uint64() {
		block, err := s.client.BlockByNumber(s.ctx, new(big.Int).SetUint64(current))
		if err != nil {
			s.logger.Printf("failed to retrieve a block: %v, %v, chainid %s ", current, err, s.chainID)
			continue
		}

		_, err = s.parseBlock(block)
		if err != nil {
			s.logger.Printf("failed to scrape block: %v", err)
		}

		current = current - 1
		s.logger.Printf("decrease current %d  ", current)

	}
	s.wg.Done()

	s.logger.Printf(" done  oracles  ")

	return nil
}

func (s *scraperImpl) UpdateHistorical() error {
	log.Printf("Scrapping started for chain %s, up to minimum block %s, maximum block %s and total oracles %d UpdateHistorical ", s.chainID, s.minblock, s.maxblock, len(s.oracles))
	s.isHistorical = true
	go s.historical()
	return nil
}

func (s *scraperImpl) UpdateRecent() error {
	log.Printf("Scrapping started for chain %s, up to minimum block %s, maximum block %s and total oracles %d UpdateRecent ", s.chainID, s.minblock, s.maxblock, len(s.oracles))
	go s.recent(0)

	return nil
}

func (s *scraperImpl) UpdateEvents(oracleaddresses []common.Address) error {
	log.Printf("UpdateEvents started for chain %s, up to minimum block %s, maximum block %s and total oracles %d UpdateRecent ", s.chainID, s.minblock, s.maxblock, len(s.oracles))

	go s.listenEvents(oracleaddresses)

	return nil
}
