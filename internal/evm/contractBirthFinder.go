// credits: https://levelup.gitconnected.com/how-to-get-smart-contract-creation-block-number-7f22f8952be0

package evm

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractBirthFinder struct {
	client      *ethclient.Client
	latestBlock int64
}

func NewContractBirthFinder(conn *ethclient.Client) (*ContractBirthFinder, error) {
	// get latest block number for reuse later
	latestBlock, err := conn.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &ContractBirthFinder{
		client:      conn,
		latestBlock: latestBlock.Number().Int64(),
	}, nil
}

func (c *ContractBirthFinder) codeLen(contractAddr string, blockNumber int64) int {
	log.Println(blockNumber)
	ctx := context.Background()
	data, err := c.client.CodeAt(ctx, common.HexToAddress(contractAddr), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal("Failed to call CodeAt", err)
	}

	return len(data)
}

func (c *ContractBirthFinder) GetContractCreationBlock(contractAddr string) uint64 {
	return uint64(c.getCreationBlock(contractAddr, 0, c.latestBlock))
}

// binary search
func (c *ContractBirthFinder) getCreationBlock(contractAddr string, startBlock int64, endBlock int64) int64 {
	if startBlock == endBlock {
		return startBlock
	}

	midBlock := (startBlock + endBlock) / 2
	codeLen := c.codeLen(contractAddr, midBlock)

	if codeLen > 2 {
		return c.getCreationBlock(contractAddr, startBlock, midBlock)
	} else {
		return c.getCreationBlock(contractAddr, midBlock+1, endBlock)
	}
}
