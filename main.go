package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://198.13.47.125:8545")
	if err != nil {
		log.Fatal(err)
	}

	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NetworkID: ", networkID)

	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// pre-fund account
	from := common.HexToAddress("a3399f17f5ade94ff61c4c4adae586711cc4b043")
	to := common.HexToAddress("a5b3c09177c196ad2e49b93b5b735c938e599453")
	data := []byte("Nexty Testnet")
	value := big.NewInt(10000000000000000)

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

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: from}, "i3nxx1rk", newTx, networkID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Send tnx succesfully!")
}
