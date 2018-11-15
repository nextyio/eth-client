package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MaxRequest maximum request
const MaxRequest = 100

func main() {
	// Some pre-fund accounts
	tos := [1]common.Address{
		common.HexToAddress("0e47Dcb26e0C3E8b7f363B738aE81aAe9FcE0004"),
	}

	// Getting account address from `keystore` folder
	files, err := ioutil.ReadDir("/home/ubuntu/.ethereum/keystore")
	if err != nil {
		log.Fatal(err)
	}
	workers := len(files)
	if workers > MaxRequest {
		workers = MaxRequest
	}
	froms := make([]common.Address, workers)
	for i, f := range files {
		if i >= workers {
			break
		}
		froms[i] = common.HexToAddress(f.Name()[37:])
	}

	// Reading data from file
	data, err := ioutil.ReadFile("./data")
	if err != nil {
		log.Fatal(err)
	}

	// Unlock all the account before sending txs
	ks := keystore.NewKeyStore("/home/ubuntu/.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	for k := 0; k < len(froms); k++ {
		ks.TimedUnlock(accounts.Account{Address: froms[k]}, "password", time.Duration(30)*time.Minute)
	}

	// Create an eth client
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NetworkID: ", networkID)

	// Send bunch of tnx to an endpoint
	for t := 0; t < len(tos); t++ {
		to := tos[t]
		value := big.NewInt(1)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var wg sync.WaitGroup
		wg.Add(workers)
		for k := 0; k < workers; k++ {
			go func(_from common.Address) {
				defer wg.Done()
				nonce, _ := client.PendingNonceAt(context.Background(), _from)
				for i := 0; i < 10000; i++ {
					msg := ethereum.CallMsg{
						From:     _from,
						To:       &to,
						GasPrice: gasPrice,
						Value:    value,
						Data:     data,
					}
					gasLimit, err := client.EstimateGas(context.Background(), msg)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
					signedTx, err := ks.SignTx(accounts.Account{Address: _from}, newTx, networkID)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					go client.SendTransaction(context.Background(), signedTx)
					fmt.Println("Send tnx succesfully...   " + strconv.Itoa(i))
					nonce++
				}
			}(froms[k])
		}
		wg.Wait()
		fmt.Println("Done!!!... " + strconv.Itoa(t))
	}
}
