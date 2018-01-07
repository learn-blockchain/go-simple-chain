package chain

import (
	"bytes"
	"errors"

	"github.com/learn-blockchain/go-simple-coin/block"
	"github.com/learn-blockchain/go-simple-coin/transaction"
)

const difficulty uint64 = 1

// Chain is the linked list of blocks
type Chain struct {
	head *block.Block
}

// New constructs a new blockchain
func New() (*Chain, error) {
	c := &Chain{}

	g, err := c.createGenesisBlock()
	if err != nil {
		return nil, err
	}

	c.head = g

	return c, nil
}

func (c *Chain) createGenesisBlock() (*block.Block, error) {
	return block.New(block.Props{
		Index:        0,
		Transactions: nil,
		PrevBlock:    nil,
		PrevHash:     nil,
	}, difficulty)
}

// GetLatestBlock returns the most recent block on the chain
func (c *Chain) GetLatestBlock() *block.Block {
	return c.head
}

// AddBlock adds new transactions to the blockchain as a block
func (c *Chain) AddBlock(t []transaction.Payment) error {
	prevBlock := c.GetLatestBlock()
	prevProps := prevBlock.GetProps()
	newBlock, err := block.New(block.Props{
		Index:        prevProps.Index + 1,
		Transactions: t,
		PrevBlock:    prevBlock,
		PrevHash:     prevBlock.GetProps().Hash,
	}, difficulty)
	if err != nil {
		return err
	}

	if c.head == nil {
		return errors.New("chain has not been initialized")
	}

	c.head = newBlock
	return nil
}

// IsChainValid checks if the chain is not corrupted
func (c Chain) IsChainValid() (bool, error) {
	dummyBlock := c.GetLatestBlock()

	for dummyBlock.GetProps().PrevBlock != nil {
		// Check #1: is the current dummy block's hash correct?
		hash, err := dummyBlock.CalcHash()
		if err != nil {
			return false, err
		}
		if !bytes.Equal(dummyBlock.GetProps().Hash, hash) {
			return false, nil
		}

		// Check #2: does the prevHash prop match the hash of the previous block?
		if !bytes.Equal(dummyBlock.GetProps().PrevBlock.GetProps().Hash, dummyBlock.GetProps().PrevHash) {
			return false, nil
		}

		// reset the dummyBlock to the prevBlock
		dummyBlock = dummyBlock.GetProps().PrevBlock
	}

	return true, nil
}
