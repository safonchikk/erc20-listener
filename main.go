package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"sync"
)

const (
	ethereumNodeURL     = "wss://mainnet.infura.io/ws/v3/76e1210cd141441da26a50f5f14735eb"
	usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	contractAbi         = `[{
							  "anonymous": false,
							  "inputs": [
								{"indexed": true, "name": "from", "type": "address"},
								{"indexed": true, "name": "to", "type": "address"},
								{"indexed": false, "name": "value", "type": "uint256"}
							  ],
							  "name": "Transfer",
							  "type": "event"
							}]`
)

type TransferEvent struct {
	Value *big.Int
}

func main() {
	client, err := ethclient.Dial(ethereumNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(usdtContractAddress)
	usdtAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		log.Fatal(err)
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{usdtAbi.Events["Transfer"].ID}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs:
				event := new(TransferEvent)
				err := usdtAbi.UnpackIntoInterface(event, "Transfer", vLog.Data)
				if err != nil {
					log.Println("Error unpacking log data:", err)
					continue
				}

				from := common.HexToAddress(vLog.Topics[1].Hex())
				to := common.HexToAddress(vLog.Topics[2].Hex())

				fmt.Printf("Received USDT Transfer Event: From %s, To %s, Tokens %s\n", from,
					to, (event.Value.Div(event.Value, big.NewInt(1000000))).String())
			}
		}
	}()

	wg.Wait()

}
