package src

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const (
	dbFile = "blockchain.db"
	blocksBucket = "blocks"
	lastTip = "lastTip"
)

// tip 指的是存储最后一个块的哈希
// db 存储数据库连接
type BlockChain struct {
	tip []byte
	Db  *bolt.DB
}

func (bc *BlockChain) AddBlock(data string)  {
	var lastHash []byte

	// 首先获取最后一个块的哈希用于生成新块的哈希
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastTip))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)


	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte(lastTip), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})

}
func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.serialize())
			if err != nil {
				log.Panic(err)
			}


			err = b.Put([]byte(lastTip),genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(lastTip))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}


	return &BlockChain{tip,db}
}

type  BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockchainIterator  {
	return &BlockchainIterator{bc.tip, bc.Db}

}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}