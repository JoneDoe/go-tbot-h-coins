package main

import (
	"flag"
	"fmt"
	"os"

	blockchain "go-tbot-h-coins/src/blockchain"
)

// CLI cli
type CLI struct {
	bc *blockchain.Blockchain
}

// Run cli
func (cli *CLI) Run() {
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	if len(os.Args[1:]) != 0 {
		switch os.Args[1] {
		case "printchain":
			printChainCmd.Parse(os.Args[2:])
		}
	}

	if printChainCmd.Parsed() {
		cli.printChain()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
