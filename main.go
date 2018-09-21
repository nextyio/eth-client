package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Some pre-fund accounts
	tos := [3]common.Address{
		common.HexToAddress("a3399f17f5ade94ff61c4c4adae586711cc4b043"),
		common.HexToAddress("5b154d28aeffb63602a326f140b8757246171546"),
		common.HexToAddress("2ccb075ade031ba82c48e6885da8577d57f3abc9"),
	}

	// Getting account address from `keystore` folder
	files, err := ioutil.ReadDir("./keystore")
	if err != nil {
		log.Fatal(err)
	}
	froms := make([]common.Address, len(files))
	for i, f := range files {
		froms[i] = common.HexToAddress(f.Name()[37:])
	}

	// Reading data from file
	data, err := ioutil.ReadFile("./data")
	if err != nil {
		log.Fatal(err)
	}

	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	to := tos[0]
	value := big.NewInt(10000)
	for k := 0; k < len(froms); k++ {
		start := time.Now()
		client, err := ethclient.Dial("http://198.13.47.125:8545")
		if err != nil {
			log.Fatal(err)
		}
		networkID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("NetworkID: ", networkID)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		nonce, _ := client.PendingNonceAt(context.Background(), froms[k])
		elapsed := time.Since(start)
		fmt.Println("PendingNonceAt took " + elapsed.String())
		msg := ethereum.CallMsg{
			From:     froms[k],
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
		elapsed = time.Since(start)
		fmt.Println("EstimateGas took " + elapsed.String())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
		elapsed = time.Since(start)
		fmt.Println("NewTransaction took " + elapsed.String())
		signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: froms[k]}, "i3nxx1rk", newTx, networkID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		elapsed = time.Since(start)
		fmt.Println("SignTxWithPassphrase took " + elapsed.String())
		go client.SendTransaction(context.Background(), signedTx)
		elapsed = time.Since(start)
		fmt.Println("SendTransaction took " + elapsed.String())
	}
}
