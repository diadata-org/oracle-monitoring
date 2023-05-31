package scraper

import (
    "bytes"
    "context"
    "fmt"
    "math/big"
    "strconv"
    "strings"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"

    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/diadata-org/oracle-monitoring/internal/database"
    "github.com/diadata-org/oracle-monitoring/internal/helpers"
)

func ConnectToNode(url string) (*ethclient.Client, error) {
    return ethclient.DialContext(context.Background(), url)
}

func GetTransactionSender(tx *types.Transaction) (common.Address, error) {
    return types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
}

func IsTargetingContract(tx *types.Transaction, address common.Address) bool {
    return tx.To() != nil && *tx.To() == address
}

func IsContractCreation(tx *types.Transaction, receipt *types.Receipt, address common.Address) bool {
    if tx.To() == nil {
        // Transaction is a contract creation
        return receipt != nil && receipt.ContractAddress == address
    }

    return false
}

func IsOracleUpdate(tx *types.Transaction, oracle *helpers.Oracle) bool {
    data := tx.Data()
    setValueSig := oracle.ContractABI.Methods["setValue"].ID

    return bytes.Equal(data[:4], setValueSig[:4])
}

func ParseTransactionMetadata(client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt) (*helpers.TransactionMetadata, error) {
    metadata := &helpers.TransactionMetadata{}

    metadata.BlockNumber = block.Number().String()
    metadata.BlockTimestamp = strconv.FormatUint(block.Time(), 10)
    metadata.TransactionHash = strings.ToLower(tx.Hash().String())
    metadata.TransactionCost = new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice).String()
    metadata.GasUsed = strconv.FormatUint(receipt.GasUsed, 10)
    metadata.GasCost = receipt.EffectiveGasPrice.String()

    sender, err := GetTransactionSender(tx)
    if err != nil {
        return nil, fmt.Errorf("failed to get sender: %v", err)
    }
    metadata.TransactionFrom = strings.ToLower(sender.String())

    if tx.To() != nil {
        metadata.TransactionTo = strings.ToLower(tx.To().String())
    }

    senderBalance, err := client.BalanceAt(context.Background(), sender, block.Number())
    if err != nil {
        return nil, fmt.Errorf("failed to get sender balance: %v", err)
    }
    metadata.SenderBalance = senderBalance.String()

    return metadata, nil
}

func ParseOracleUpdate(tx *types.Transaction, oracle *helpers.Oracle) (*helpers.OracleUpdate, error) {
    event := &helpers.OracleUpdate{}
    err := oracle.ContractABI.UnpackIntoInterface(event, "setValue", tx.Data()[4:])

    return event, err
}

func ParseTransaction(client *ethclient.Client, block *types.Block, tx *types.Transaction, receipt *types.Receipt, oracle *helpers.Oracle, db *pgxpool.Pool) (bool, error) {
    done := false
    metadata, err := ParseTransactionMetadata(client, block, tx, receipt)
    if err != nil {
        return false, fmt.Errorf("failed to parse transaction metadata: %v", err)
    }

    if IsContractCreation(tx, receipt, oracle.ContractAddress) {
        // to is nil on contract creation
        metadata.TransactionTo = strings.ToLower(oracle.ContractAddress.String())

        err := database.UpdateOracleCreation(db, oracle.ContractAddress.String(), block.Number().Uint64())
        if err != nil {
            return false, fmt.Errorf("failed to update oracle creation data: %v", err)
        }

        // no more transactions to scrape after the oracle creation
        done = true
    }

    if IsTargetingContract(tx, oracle.ContractAddress) && IsOracleUpdate(tx, oracle) {
        oracleUpdate, err := ParseOracleUpdate(tx, oracle)
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

        err = database.InsertOracleMetrics(db, metrics)
        if err != nil {
            return false, fmt.Errorf("failed to insert oracle metrics: %v", err)
        }

        // reached the latest block scraped on the previous run
        done = block.Number().Cmp(oracle.LatestScrapedBlock) <= 0 // -1 or 0 if block.Number is smaller
    }

    return done, nil
}

func ParseBlock(client *ethclient.Client, block *types.Block, oracle *helpers.Oracle, db *pgxpool.Pool) (bool, error) {
    done := false

    for _, tx := range block.Transactions() {
        receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
            return false, fmt.Errorf("failed to get transaction receipt: %v", err)
        }

        creation, err := ParseTransaction(client, block, tx, receipt, oracle, db)
        if err != nil {
            return false, fmt.Errorf("failed to scrape transaction: %v", err)
        }

        done = done || creation
    }

    return done, nil
}

func ScrapeSingleOracle(oracle *helpers.Oracle, db *pgxpool.Pool) error {
    done := false

    client, err := ConnectToNode(oracle.NodeUrl)
    if err != nil {
        return fmt.Errorf("failed to connect to the node: %v", err)
    }

    current, err := client.BlockNumber(context.Background())
    if err != nil {
        return fmt.Errorf("failed to retrieve the latest block: %v", err)
    }

    current = 17151362

    for !done && current > 0 && current > oracle.LatestScrapedBlock.Uint64() {
        block, err := client.BlockByNumber(context.Background(), new(big.Int).SetUint64(current))
        if err != nil {
            return fmt.Errorf("failed to retrieve a block: %v", err)
        }

        fmt.Println(block.Number().Uint64())

        done, err = ParseBlock(client, block, oracle, db)
        if err != nil {
            return fmt.Errorf("failed to scrape block: %v", err)
        }

        current = current - 1
    }

    return nil
}

func Update(oracles []helpers.Oracle, db *pgxpool.Pool) error {
    for _, oracle := range oracles {
        err := ScrapeSingleOracle(&oracle, db)
        if err != nil {
            return fmt.Errorf("failed to scrape oracle: %v", err)
        }
    }

    return nil
}
