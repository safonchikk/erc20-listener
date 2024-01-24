package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	ethereumNodeURL     = "wss://mainnet.infura.io/ws/v3/76e1210cd141441da26a50f5f14735eb"
	usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	transferEventName   = "Transfer"
)

type TransferEvent struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func main() {
	client, err := ethclient.Dial(ethereumNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(usdtContractAddress)
	abiBytes, err := os.ReadFile("usdt.abi")
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatal(err)
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {

		event := new(TransferEvent)
		err := contractAbi.UnpackIntoInterface(event, transferEventName, vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received USDT Transfer Event: From %s, To %s, Tokens %s\n", event.From.Hex(), event.To.Hex(), event.Tokens.String())

	}

}
