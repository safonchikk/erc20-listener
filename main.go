package main

import (
	"context"
	"erc20-listener/util"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

const (
	contractAbi = `[{
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
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading app.env file" + err.Error())
	}

	client, err := ethclient.Dial(config.EthNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(config.ContractAddr)
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

				tokens, _ := event.Value.Float64()
				tokDec, err := strconv.Atoi(config.TokenDecimal)
				if err != nil {
					log.Println("Invalid token decimal")
					tokDec = 0
				}
				tokens *= math.Pow(0.1, float64(tokDec))

				fmt.Printf("Received Transfer Event: From %s, To %s, Tokens %f\n", from, to, tokens)
			}
		}
	}()

	wg.Wait()

}
