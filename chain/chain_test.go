package chain

import (
	"reflect"
	"testing"

	"github.com/learn-blockchain/go-simple-coin/transaction"
)

func TestChain(t *testing.T) {
	chain, err := New()
	if err != nil {
		t.Errorf("could not create blockchain; err: %v", err)
	}

	valid, err := chain.IsChainValid()
	if err != nil {
		t.Errorf("could not validate blockchain; err: %v", err)
	}
	if !valid {
		t.Errorf("a new chain is not valid!")
	}

	t1 := transaction.Payment{
		From:   "adam",
		To:     "john",
		Amount: 100.0,
	}
	t2 := transaction.Payment{
		From:   "adam",
		To:     "jane",
		Amount: 100.0,
	}
	t3 := transaction.Payment{
		From:   "jane",
		To:     "adam",
		Amount: 5.0,
	}
	trans := []transaction.Payment{t1, t2, t3}

	err = chain.AddBlock(trans)
	if err != nil {
		t.Errorf("could not add transactions; err: %v", err)
	}
	valid, err = chain.IsChainValid()
	if err != nil {
		t.Errorf("could not validate blockchain; err: %v", err)
	}
	if !valid {
		t.Errorf("a chain with one added transaction is not valid!")
	}
	assert := reflect.DeepEqual(trans, chain.GetLatestBlock().GetProps().Transactions)
	if !assert {
		t.Errorf("the transactions on the current block don't match what was put there; trans: %v, current block trans: %v", trans, chain.GetLatestBlock().GetProps().Transactions)
	}

	numBlocks := 1
	dummyBlock := chain.GetLatestBlock()
	for dummyBlock.GetProps().PrevBlock != nil {
		dummyBlock = dummyBlock.GetProps().PrevBlock
		numBlocks++
	}
	if numBlocks != 2 {
		t.Errorf("Expected chain to have two blocks but found %d", numBlocks)
	}

	t4 := transaction.Payment{
		From:   "frank",
		To:     "beth",
		Amount: 100.0,
	}
	t5 := transaction.Payment{
		From:   "jimmy",
		To:     "timmy",
		Amount: 100.0,
	}
	t6 := transaction.Payment{
		From:   "dude1",
		To:     "dude2",
		Amount: 5.0,
	}
	trans2 := []transaction.Payment{t4, t5, t6}
	err = chain.AddBlock(trans2)
	if err != nil {
		t.Errorf("could not add transactions; err: %v", err)
	}
	valid, err = chain.IsChainValid()
	if err != nil {
		t.Errorf("could not validate blockchain; err: %v", err)
	}
	if !valid {
		t.Errorf("a chain with two added transactions is not valid!")
	}
	assert = reflect.DeepEqual(trans2, chain.GetLatestBlock().GetProps().Transactions)
	if !assert {
		t.Errorf("the transactions on the current block don't match what was put there; trans: %v, current block trans: %v", trans2, chain.GetLatestBlock().GetProps().Transactions)
	}
	numBlocks = 1
	dummyBlock = chain.GetLatestBlock()
	for dummyBlock.GetProps().PrevBlock != nil {
		dummyBlock = dummyBlock.GetProps().PrevBlock
		numBlocks++
	}
	if numBlocks != 3 {
		t.Errorf("Expected chain to have three blocks but found %d", numBlocks)
	}

	dummyBlock = chain.GetLatestBlock()
	dummyProps := dummyBlock.GetProps()
	var wrongIdx uint64
	wrongIdx = 5000
	dummyProps.Index = wrongIdx
	if chain.GetLatestBlock().GetProps().Index == wrongIdx || chain.GetLatestBlock().GetProps().Index != 2 {
		t.Errorf("Expected block to be immutable and for index of latest block to be 2. Found idx: %d", chain.GetLatestBlock().GetProps().Index)
	}
}
