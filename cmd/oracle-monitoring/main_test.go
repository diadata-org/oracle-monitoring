package main

import (
	"context"
	"testing"
	"time"

	"github.com/diadata-org/oracle-monitoring/internal/config"
	"github.com/diadata-org/oracle-monitoring/internal/database"
	"github.com/diadata-org/oracle-monitoring/internal/evm"
	"github.com/diadata-org/oracle-monitoring/internal/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Initialize the configuration and database connections
	config, err := config.LoadConfig()
	assert.NoError(t, err)

	db, err := database.New(config.Database)
	assert.NoError(t, err)
	defer db.Close()

	// Call main function in a goroutine
	go main()

	// Wait for main function to initialize and connect to the EVM client
	time.Sleep(5 * time.Second)

	// Delete previous DB records
	_, err = db.Pool.Exec(context.Background(), "DELETE FROM events")
	assert.NoError(t, err)

	// Get the EVM evm and contract address
	evm, err := evm.New(&config.EVM)
	assert.NoError(t, err)

	/*
		Calldata for: setValue(string key,uint128 value,uint128 timestamp)
		key = BTC/USD
		value = 1000
		timestamp = 1683234213
	*/
	data := common.Hex2Bytes("7898e0c2000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000003e80000000000000000000000000000000000000000000000000000000064541da500000000000000000000000000000000000000000000000000000000000000074254432f55534400000000000000000000000000000000000000000000000000")
	tx, err := evm.SendTransaction(config.EVM.OracleAddress, config.PrivateKey, data, 50000)
	assert.NoError(t, err)
	_, err = bind.WaitMined(context.Background(), evm.Client, tx)
	assert.NoError(t, err)

	// Wait until block is mined and check if record was created in database and contract creation timestamp got populated
	time.Sleep(15 * time.Second)

	event, err := models.GetLatestEventByOracleAddress(db, config.EVM.OracleAddress)
	assert.NoError(t, err)

	assert.NotEqual(t, uint64(0), event.OracleCreationTimestamp, "Oracle creation timestamp should not be 0")

	// Call contract again with same calldata
	tx, err = evm.SendTransaction(config.EVM.OracleAddress, config.PrivateKey, data, 50000)
	assert.NoError(t, err)
	_, err = bind.WaitMined(context.Background(), evm.Client, tx)
	assert.NoError(t, err)

	// Wait until block is mined and check if 2nd record was created
	time.Sleep(15 * time.Second)
	row := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM events WHERE oracle_address = $1", config.EVM.OracleAddress)
	var count int
	err = row.Scan(&count)
	assert.NoError(t, err)

	assert.Equal(t, 2, count, "There should be 2 records in the database for the contract address")
}
