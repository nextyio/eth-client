package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// pre-fund account
	from := common.HexToAddress("0e47Dcb26e0C3E8b7f363B738aE81aAe9FcE0004")

	// Getting account address from `keystore` folder
	files, err := ioutil.ReadDir("/home/ubuntu/.ethereum/keystore")
	if err != nil {
		log.Fatal(err)
	}

	tos := make([]common.Address, len(files))
	for i, f := range files {
		tos[i] = common.HexToAddress(f.Name()[37:])
	}
	data := []byte("nexty testnet funding")

	// Unlock the account before sending txs
	ks := keystore.NewKeyStore("/home/ubuntu/.ethereum/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	ks.Unlock(accounts.Account{Address: from}, "password")

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

	// Fund to the accounts
	nonce, _ := client.PendingNonceAt(context.Background(), from)
	for i := 0; i < len(tos); i++ {
		to := tos[i]
		value := big.NewInt(100000)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		msg := ethereum.CallMsg{
			From:     from,
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
		signedTx, err := ks.SignTx(accounts.Account{Address: from}, newTx, networkID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := client.SendTransaction(context.Background(), signedTx); err != nil {
			fmt.Println(err.Error())
			return
		}
		nonce++
		fmt.Println("Send tnx succesfully...   " + strconv.Itoa(i))
	}
}
