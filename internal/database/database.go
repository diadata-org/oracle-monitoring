package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/diadata-org/oracle-monitoring/internal/helpers"
)

const (
	oraclesTable       = "oracles"
	feederupdatesTable = "feederupdates"
	chainconfig        = "chainconfig"
)

const (
	updateOraclesCreationQuery = "UPDATE oracleconfig SET creation_block = $2, creation_block_time=$3 WHERE address = $1 and chainid =$4"
	selectOraclesQuery         = `SELECT address, chainid,  COALESCE(latest.scraped_block, 0) AS latest_scraped_block FROM oracleconfig LEFT JOIN (SELECT oracle_address, chain_id, MAX(update_block) AS scraped_block FROM feederupdates GROUP BY oracle_address,chain_id) latest ON (oracleconfig.address = latest.oracle_address and oracleconfig.chainid = latest.chain_id) WHERE  oracleconfig.chainid = '%s'`
	selectLatestOraclesQuery   = `SELECT address, chainid,  createddate,COALESCE(latest.scraped_block, 0) AS latest_scraped_block FROM oracleconfig LEFT JOIN (SELECT oracle_address, chain_id, MAX(update_block) AS scraped_block FROM feederupdates GROUP BY oracle_address,chain_id) latest ON (oracleconfig.address = latest.oracle_address and oracleconfig.chainid = latest.chain_id) WHERE  oracleconfig.chainid = '%s' and oracleconfig.createddate > '%s'`
)

// Database is an interface that represents the required database operations.
type Database interface {
	Connect() error
	InsertOracles(targets []helpers.Target) error
	UpdateOracleCreation(address string, block string, blocktime time.Time, chainid string) error
	SelectOracles(string) ([]helpers.Target, error)
	InsertOracleMetrics(metrics *helpers.OracleMetrics) error
	GetRPCByChainID() (map[string]string, error)
	GetWSByChainID() (map[string]string, error)
	SelectOraclesWithCreationTime(chainID string, lastCreatedTime time.Time) ([]helpers.Target, error)

	Close()
}

type postgresDB struct {
	db *pgxpool.Pool
}

// NewPostgresDB creates a new instance of the Database interface with PostgreSQL implementation.
func NewPostgresDB() Database {
	return &postgresDB{}
}

func (pdb *postgresDB) Connect() error {
	// Load database connection info from .env file
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Construct the connection URI
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Create connection pool
	pdb.db, err = pgxpool.New(context.Background(), uri)
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %v", err)
	}

	return nil
}

func (pdb *postgresDB) InsertOracles(targets []helpers.Target) error {
	// Insert targets into the database
	// for _, target := range targets {
	// 	_, err := pdb.db.Exec(context.Background(), insertOraclesQuery, target.ContractAddress, target.ContractABI, target.ChainId, target.NodeUrl, 0)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert oracle: %v", err)
	// 	}
	// }

	return nil
}

func (pdb *postgresDB) UpdateOracleCreation(address string, block string, blocktime time.Time, chainid string) error {
	// Update the creation block of an oracle
	_, err := pdb.db.Exec(context.Background(), updateOraclesCreationQuery, address, block, blocktime, chainid)
	if err != nil {
		return fmt.Errorf("failed to update the creation block in the DB: %v", err)
	}

	return nil
}
func (pdb *postgresDB) SelectOraclesWithCreationTime(chainID string, lastCreatedTime time.Time) ([]helpers.Target, error) {
	targets := []helpers.Target{}

	// Retrieve oracles with the latest scraped block

	query := fmt.Sprintf(selectLatestOraclesQuery, chainID, lastCreatedTime.Format("2006-01-02 15:04:05.999999-07:00"))

	rows, err := pdb.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the DB query: %v", err)
	}
	defer rows.Close()

	// Format the results as struct
	for rows.Next() {
		var target helpers.Target
		err := rows.Scan(&target.ContractAddress, &target.ChainId, &target.CreatedDate, &target.LatestScrapedBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to get the list of targets from the DB: %v", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (pdb *postgresDB) SelectOracles(chainID string) ([]helpers.Target, error) {
	targets := []helpers.Target{}

	// Retrieve oracles with the latest scraped block

	query := fmt.Sprintf(selectOraclesQuery, chainID)
	rows, err := pdb.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the DB query: %v", err)
	}
	defer rows.Close()

	// Format the results as struct
	for rows.Next() {
		var target helpers.Target
		err := rows.Scan(&target.ContractAddress, &target.ChainId, &target.LatestScrapedBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to get the list of targets from the DB: %v", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (pdb *postgresDB) InsertOracleMetrics(metrics *helpers.OracleMetrics) error {
	insertMetricsQuery := fmt.Sprintf("INSERT INTO %s (oracle_address,transaction_hash,transaction_cost,asset_key,asset_price,update_block, update_from, from_balance, gas_cost, gas_used,creation_block,chain_id,update_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11,$12,$13)", feederupdatesTable)

	fmt.Printf(
		"-- Inserted Metrics --\n"+
			"TransactionTo: %s"+
			"TransactionHash: %s"+
			"TransactionCost: %s"+
			"AssetKey: %s"+
			"AssetPrice: %s"+
			"BlockNumber: %s"+
			"TransactionFrom: %s"+
			"SenderBalance: %s"+
			"GasCost: %s"+
			"GasUsed: %s"+
			"ChainID: %s"+
			"BlockTimestamp: %s",
		metrics.TransactionTo,
		metrics.TransactionHash,
		metrics.TransactionCost,
		metrics.AssetKey,
		metrics.AssetPrice,
		metrics.BlockNumber,
		metrics.TransactionFrom,
		metrics.SenderBalance,
		metrics.GasCost,
		metrics.GasUsed,
		metrics.ChainID,
		metrics.BlockTimestamp,
	)

	metrics.CreationBlock = "0"
	fmt.Println("metrics.CreationBlock", metrics.CreationBlock)

	// Insert metrics into the database
	_, err := pdb.db.Exec(
		context.Background(),
		insertMetricsQuery,
		metrics.TransactionTo.String(),
		metrics.TransactionHash,
		metrics.TransactionCost,
		metrics.AssetKey,
		metrics.AssetPrice,
		metrics.BlockNumber,
		metrics.TransactionFrom.String(),
		metrics.SenderBalance,
		metrics.GasCost,
		metrics.GasUsed,
		metrics.CreationBlock,
		metrics.ChainID,
		metrics.BlockTimestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to insert metrics in the DB: %v", err)
	}

	return nil
}

// GetRPCByChainID returns the RPC URL for the given chain ID.
func (pdb *postgresDB) GetRPCByChainID() (rpc map[string]string, err error) {

	query := `SELECT rpcurl, chainid from %s`

	rpc = make(map[string]string)

	rows, err := pdb.db.Query(context.Background(), fmt.Sprintf(query, chainconfig))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the DB query: %v", err)
	}

	for rows.Next() {

		var rpcurl, chainid string
		err := rows.Scan(&rpcurl, &chainid)
		if err != nil {
			return nil, fmt.Errorf("failed to get the list of rpc from the DB: %v", err)
		}
		rpc[chainid] = rpcurl
	}
	return rpc, nil
}

// GetWSByChainID returns the WS URL for the given chain ID.
func (pdb *postgresDB) GetWSByChainID() (rpc map[string]string, err error) {

	query := `SELECT wsurl, chainid from %s`

	rpc = make(map[string]string)

	rows, err := pdb.db.Query(context.Background(), fmt.Sprintf(query, chainconfig))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from the DB query: %v", err)
	}

	for rows.Next() {

		var rpcurl, chainid string
		err := rows.Scan(&rpcurl, &chainid)
		if err != nil {
			return nil, fmt.Errorf("failed to get the list of rpc from the DB: %v", err)
		}
		rpc[chainid] = rpcurl
	}
	return rpc, nil
}

func (pdb *postgresDB) Close() {
	pdb.db.Close()
}

// SELECT address, chainid,  COALESCE(latest.scraped_block, 0) AS latest_scraped_block FROM oracleconfig LEFT JOIN (SELECT oracle_address,chain_id, MAX(update_block) AS scraped_block FROM feederupdates GROUP BY oracle_address,chain_id) latest ON oracleconfig.address = latest.oracle_address;
