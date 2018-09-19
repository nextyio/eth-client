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
	value := big.NewInt(10000000000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	nonce, err := client.PendingNonceAt(context.Background(), from)
	for i := 0; i < 100; i++ {
		data := []byte(`When I handed my mother the green slip, she read the words back to me like she was returning borrowed things. 
		She had kept her ugly close to her once, a toothbrush she returned to twice a day, a common ritual for the women with our faces. 
		My mother knew she could smooth the ugly out of me as she did with her own ugly. 
		She showed me pictures of Angela Davis, bought records with Diana on the sleeve, staring a hole through the ugly in me.

		My mother threw me these women, bones that called me to attention. 
		But when I pressed my face against their faces, put my jaw to Angela’s jaw, 
		I couldn’t make out the similarities, couldn’t find the same shine in my temples. 
		I knew these women were supposed to remind me of myself, but there was an obvious disconnect between the Black girls presented to me as beautiful, 
		and the Black girl written on the green slip.

		When we talk about the importance of representation, we often start with children. 
		Every child deserves to see a version of themselves on a screen—to point at, to call dibs on for a playground reenactment. 
		Representation does the job of affirmation; it affirms the identity of the questioning child, the child in a predominantly white classroom, 
		the child accustomed to different spaces.

		When I came home from school to watch TV, it would have been valuable to see a cartoon of a Black girl in the suburbs eating grits for breakfast. 
		I would have cherished a sitcom with a Black girl with blue hair, if only as a wink and a nudge, a nod to my quiet existence.
		
		I was 18 when Black Panther came out, and like most Black Americans, I walked around with a glow. 
		I had walked through my entire youth without a movie like it. 
		I saw the movie for the first time in Atlanta, which made the whole experience Blacker than I thought possible. 
		I wore my pride like an unconcealed weapon, walked with it until Twitter got quiet, until my mother stopped talking about it. 
		The movie will always be a great example of the magic of representation, how it transfixes us.		
		`)
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
		nonce++
	}
}
