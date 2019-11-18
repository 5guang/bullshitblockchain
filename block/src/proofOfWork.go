package src

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int
}

// 准备hash的数据
func (pow *ProofOfWork) prepareData(nonce int) []byte  {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			Int64ToBytes(pow.block.Timestamp),
			Int64ToBytes(int64(targetBits)),
			Int64ToBytes(int64(nonce)),
		},
		[]byte{},
		)
}

// 挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	// 使用了一个 大整数 ，把hash转换成一个大整数， 然后进行比较
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		// 将hash转换成一个大整数
		hashInt.SetBytes(hash[:])
		// 比较x和y的大小。x<y时返回-1；x>y时返回+1；否则返回0。
		// 就是说生成的hash前3个字节（24位）必须为0
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")
	return nonce, hash[:]
}

// 验证工作量
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 说明 这里将target 左移256 - targetBits位 （256 是一个 SHA-256 哈希的位数）
	// 保证hash前targetBits位必须为0
	target.Lsh(target, uint(256 - targetBits))
	return &ProofOfWork{block, target}
}