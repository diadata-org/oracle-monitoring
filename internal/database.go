package database

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "math/big"
    "os"

    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5/pgxpool"
)

const INSERT_ORACLES_QUERY = "INSERT INTO oracles (address, creation_block) VALUES ($1, $2) ON CONFLICT DO NOTHING"
const UPDATE_ORACLES_CREATION_QUERY = "UPDATE oracles SET creation_block = $2 WHERE address = $1"
const SELECT_ORACLES = "SELECT oracles.address as address, latest.scraped_block as latest_scraped_block FROM oracles INNER JOIN (SELECT oracle_address, MAX(update_block) as scraped_block FROM metrics) latest ON oracles.address=latest.oracle_address"
const INSERT_METRICS_QUERY = "INSERT INTO metrics (oracle_address, transaction_hash, transaction_cost, asset_key, asset_price, update_block, update_from, from_balance, gas_cost, gas_used) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

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
        "postgres://%s:%s@%s:%s/%s?sslmode=require",
        dbUser, dbPass, dbHost, dbPort, dbName)

    // Create connection pool
    db, err := pgxpool.New(context.Background(), uri)
    if err != nil {
        return nil, fmt.Errorf("unable to connect to the database: %v", err)
    }

    return db, nil
}

func InsertOracles(db *pgxpool.Pool, oracles []config.Oracle) error {
    // Insert oracles into the database
    for _, oracle := range oracles {
        _, err := db.Exec(context.Background(), INSERT_ORACLES_QUERY, oracle.Address, 0)
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
        return fmt.Errorf("failed to query the database: %v", err)
    }

    return nil
}

func SelectOracles(db *pgxpool.Pool) ([]scraper.Oracle, error) {
    oracles := []scraper.Oracle{}

    // get all the oracles with the latest scraped block
    rows, err := db.Query(context.Background(), SELECT_ORACLES)
    if err != nil {
        return nil, fmt.Errorf("failed to cast the data from the query: %v", err)
    }

    // format as struct
    defer rows.Close()
    for rows.Next() {
        var o scraper.Oracle
        err := rows.Scan(&o.Address, &o.LatestScrapedBlock)
        if err != nil {
            return nil, fmt.Errorf("failed to update oracle creation block: %v", err)
        }

        oracles = append(oracles, o)
    }

    return oracles, nil
}

func InsertOracleMetrics(db *pgxpool.Pool, metrics *scraper.Metrics) error {
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
        return fmt.Errorf("failed to insert metrics: %v", err)
    }

    return nil
}
