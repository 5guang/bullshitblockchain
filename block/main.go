
package main

import (
	"bullShitBlockChain/block/src"
)

func main() {
	bc := src.NewBlockChain()
	defer bc.Db.Close()
	cli := src.CLI{bc}
	cli.Run()
}