package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	merkle "github.com/learn-blockchain/go-merkle-tree"
	"github.com/learn-blockchain/go-simple-coin/transaction"
)

// Props are the Block's properties
type Props struct {
	Index        uint64
	Timestamp    time.Time // note: this is set at the time the block is created
	Transactions []transaction.Payment
	PrevBlock    *Block
	Nonce        uint64
	Hash         []byte // note: this will not be set if passed to the New function, but will be calc'd, instead
	PrevHash     []byte
}

// Block is an individual block on the blockchain
type Block struct {
	p Props
}

// New returns a block
func New(props Props, difficulty uint64) (*Block, error) {
	b := &Block{
		p: Props{
			Index:        props.Index,
			Timestamp:    time.Now(),
			Transactions: props.Transactions,
			PrevBlock:    props.PrevBlock,
			PrevHash:     props.PrevHash,
		},
	}

	err := b.mine(difficulty)
	return b, err
}

// GetProps returns the block's properties
func (b Block) GetProps() Props {
	return b.p
}

// CalcHash computes the hash for this block
func (b *Block) CalcHash() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 1. Add the index of the block
	err := binary.Write(buf, binary.LittleEndian, b.p.Index)
	if err != nil {
		return nil, err
	}
	tmpBytes := buf.Bytes()

	// 2. Add the timestamp of the block
	tmpBytes = append(tmpBytes, []byte(b.p.Timestamp.String())...)

	// 3. Add the transactions of the block
	// but first, convert the transactions into a merkle tree
	transactions := merkle.Data{}
	for _, t := range b.p.Transactions {
		tBytes := []byte(fmt.Sprintf("%v", t))
		transactions = append(transactions, tBytes)
	}
	node, err := merkle.New(transactions)
	if err != nil {
		return nil, err
	}
	tmpBytes = append(tmpBytes, node.Hash...)

	// 4. Add the hash of the previous block
	if b.p.PrevHash != nil {
		tmpBytes = append(tmpBytes, b.p.PrevHash...)
	}

	// 5. Add nonce
	tmpBytes = append(tmpBytes, []byte(strconv.Itoa(int(b.p.Nonce)))...)

	hasher := sha256.New()
	hasher.Write(tmpBytes)
	return hasher.Sum(nil), nil
}

func (b *Block) setHash() error {
	h, err := b.CalcHash()
	if err != nil {
		return err
	}

	b.p.Hash = h
	return nil
}

func (b *Block) mine(difficulty uint64) error {
	var s string
	for i := 0; i < int(difficulty); i++ {
		s += "0"
	}

	err := b.setHash()
	if err != nil {
		return err
	}

	for string(b.p.Hash)[:int(difficulty)] != s {
		b.p.Nonce++
		err = b.setHash()
		if err != nil {
			return err
		}
	}

	return nil
}
