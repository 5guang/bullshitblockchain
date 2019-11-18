package src

import (
	"fmt"
	"strconv"
)

func (cli *CLI) addBlock(data string)  {
	cli.Bc.AddBlock(data)
	fmt.Println("Success")
}
func (cli *CLI) printChain()  {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
