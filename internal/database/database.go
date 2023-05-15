package database

import (
    "context"
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/diadata-org/oracle-monitoring/internal/helpers"
)

const INSERT_ORACLES_QUERY = "INSERT INTO oracles (contract_address, contract_abi, chain_id, node_url, creation_block) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"
const UPDATE_ORACLES_CREATION_QUERY = "UPDATE oracles SET creation_block = $2 WHERE contract_address = $1"
const SELECT_ORACLES = "SELECT oracles.contract_address AS contract_address, oracles.contract_abi AS contract_abi, oracles.chain_id AS chain_id, oracles.node_url AS node_url, oracles.creation_block AS creation_block, COALESCE(latest.scraped_block, 0) AS latest_scraped_block FROM oracles LEFT JOIN (SELECT oracle_address, MAX(update_block) AS scraped_block FROM metrics GROUP BY oracle_address) latest ON oracles.contract_address=latest.oracle_address"
const INSERT_METRICS_QUERY = "INSERT INTO metrics (oracle_address, transaction_hash, transaction_cost, asset_key, asset_price, update_block, update_from, from_balance, gas_cost, gas_used) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

func ConnectToDatabase() (*pgxpool.Pool, error) {
    // Load database connection info from .env file
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    // uri := "postgres://user:password@localhost:5432/dbname?sslmode=require"
    uri := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        dbUser, dbPass, dbHost, dbPort, dbName)

    // Create connection pool
    db, err := pgxpool.New(context.Background(), uri)
    if err != nil {
        return nil, fmt.Errorf("unable to connect to the database: %v", err)
    }

    return db, nil
}

func InsertOracles(db *pgxpool.Pool, targets []helpers.Target) error {
    // Insert targets into the database
    for _, target := range targets {
        _, err := db.Exec(context.Background(), INSERT_ORACLES_QUERY, target.ContractAddress, target.ContractABI, target.ChainId, target.NodeUrl, 0)
        if err != nil {
            return fmt.Errorf("failed to insert oracle: %v", err)
        }
    }

    return nil
}

func UpdateOracleCreation(db *pgxpool.Pool, address string, block uint64) error {
    // Update the creation block of an oracle
    _, err := db.Exec(context.Background(), UPDATE_ORACLES_CREATION_QUERY, address, block)
    if err != nil {
        return fmt.Errorf("failed to update the creation block in the DB: %v", err)
    }

    return nil
}

func SelectOracles(db *pgxpool.Pool) ([]helpers.Target, error) {
    targets := []helpers.Target{}

    // get all the oracles with the latest scraped block
    rows, err := db.Query(context.Background(), SELECT_ORACLES)
    if err != nil {
        return nil, fmt.Errorf("failed to cast the data from the DB query: %v", err)
    }

    // format as struct
    defer rows.Close()
    for rows.Next() {
        var target helpers.Target
        err := rows.Scan(&target.ContractAddress, &target.ContractABI, &target.ChainId, &target.NodeUrl, &target.CreationBlock, &target.LatestScrapedBlock)
        if err != nil {
            return nil, fmt.Errorf("failed to get the list of targets from the DB: %v", err)
        }

        targets = append(targets, target)
    }

    return targets, nil
}

func InsertOracleMetrics(db *pgxpool.Pool, metrics *helpers.OracleMetrics) error {
    // Insert metrics into the database
    // transaction_hash, transaction_cost, asset_key, asset_price, update_block, update_from, gas_cost, gas_used
    _, err := db.Exec(
        context.Background(),
        INSERT_METRICS_QUERY,
        metrics.TransactionTo,
        metrics.TransactionHash,
        metrics.TransactionCost,
        metrics.AssetKey,
        metrics.AssetPrice,
        metrics.BlockNumber,
        metrics.TransactionFrom,
        metrics.SenderBalance,
        metrics.GasCost,
        metrics.GasUsed)
    if err != nil {
        return fmt.Errorf("failed to insert metrics in the DB: %v", err)
    }

    return nil
}
