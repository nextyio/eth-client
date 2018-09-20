package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Some pre-fund accounts
	from := [3]common.Address{
		common.HexToAddress("a3399f17f5ade94ff61c4c4adae586711cc4b043"),
		common.HexToAddress("5b154d28aeffb63602a326f140b8757246171546"),
		common.HexToAddress("2ccb075ade031ba82c48e6885da8577d57f3abc9"),
	}

	// Getting account address from `keystore` folder
	files, err := ioutil.ReadDir("./keystore")
	if err != nil {
		log.Fatal(err)
	}
	tos := make([]common.Address, len(files))
	for i, f := range files {
		tos[i] = common.HexToAddress(f.Name()[37:])
	}

	// Reading data from file
	data, err := ioutil.ReadFile("./data")
	if err != nil {
		log.Fatal(err)
	}

	// Reading passphrase from file
	file, err := ioutil.ReadFile("passwd")
	if err != nil {
		log.Fatal(err)
	}
	passphrase := string(file)

	// Send bunch of tnx to an endpoint
	for t := 0; t < len(tos); t++ {
		client, err := ethclient.Dial("http://198.13.47.125:8545")
		if err != nil {
			log.Fatal(err)
		}
		networkID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("NetworkID: ", networkID)
		ks := keystore.NewKeyStore("./keystore/from", keystore.StandardScryptN, keystore.StandardScryptP)
		to := tos[t]
		value := big.NewInt(10000000)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var wg sync.WaitGroup
		wg.Add(len(from))
		for k := 0; k < len(from); k++ {
			go func(_from common.Address) {
				defer wg.Done()
				nonce, _ := client.PendingNonceAt(context.Background(), _from)
				for i := 0; i < 1; i++ {
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
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
					signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: _from}, passphrase, newTx, networkID)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					if err := client.SendTransaction(context.Background(), signedTx); err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Println("Send tnx succesfully...   " + strconv.Itoa(i))
					nonce++
				}
			}(from[k])
		}
		wg.Wait()
		fmt.Println("Done!!!... " + strconv.Itoa(t))
	}
}
