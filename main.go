package main

import (
	"context"
	"fmt"
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
	// pre-fund accounts
	from := [23]common.Address{
		common.HexToAddress("a3399f17f5ade94ff61c4c4adae586711cc4b043"),
		common.HexToAddress("5b154d28aeffb63602a326f140b8757246171546"),
		common.HexToAddress("2ccb075ade031ba82c48e6885da8577d57f3abc9"),
		common.HexToAddress("a5b3c09177c196ad2e49b93b5b735c938e599453"),
		common.HexToAddress("55c02a134cff5de25a28c6a81866cc72fc9ace8b"),
		common.HexToAddress("f4fded12e9f02f6b2005746a73318f11ee1acc02"),
		common.HexToAddress("cf24de6697c41311fc4a0300f4e244645dfd248b"),
		common.HexToAddress("51cfe244d4a3abc1e92ea55f6487b9da0a08dc4a"),
		common.HexToAddress("18f7c409e8311edde92072a0b08c89a13aea9479"),
		common.HexToAddress("f1cf16affcf8a3116f380d7ddc04fb91f0081385"),
		common.HexToAddress("f549eb467afaed268e4b0679ddc32b3af3c1e563"),
		common.HexToAddress("d407ae870a3082d168e167205d751fbe6f86173d"),
		common.HexToAddress("3a8918b3d821a7e1c0943edd679b58e73e2d80ab"),
	}

	tos := [20]common.Address{
		common.HexToAddress("a5b3c09177c196ad2e49b93b5b735c938e599453"),
		common.HexToAddress("55c02a134cff5de25a28c6a81866cc72fc9ace8b"),
		common.HexToAddress("f4fded12e9f02f6b2005746a73318f11ee1acc02"),
		common.HexToAddress("cf24de6697c41311fc4a0300f4e244645dfd248b"),
		common.HexToAddress("51cfe244d4a3abc1e92ea55f6487b9da0a08dc4a"),
		common.HexToAddress("18f7c409e8311edde92072a0b08c89a13aea9479"),
		common.HexToAddress("f1cf16affcf8a3116f380d7ddc04fb91f0081385"),
		common.HexToAddress("f549eb467afaed268e4b0679ddc32b3af3c1e563"),
		common.HexToAddress("d407ae870a3082d168e167205d751fbe6f86173d"),
		common.HexToAddress("3a8918b3d821a7e1c0943edd679b58e73e2d80ab"),
		common.HexToAddress("6d1b144654075e655141795c61c73118befb1be6"),
		common.HexToAddress("dfc18143836a56533c846bbc06f9d22fb3deedbb"),
		common.HexToAddress("987e89f39c3893e1fffe86a7ad7c2cd1da9013c1"),
		common.HexToAddress("557fb8310447c2d7508641cdeab59a48708e7846"),
		common.HexToAddress("dea3b32d85c1486056ea22c8796f6e7dbfa77423"),
		common.HexToAddress("0a1e36073281e1cfead456869523c4378dfb46bb"),
		common.HexToAddress("a1dec2b1916c41a222931dce1188f5ebcf237b0e"),
		common.HexToAddress("d392feb6b5a2efc911138668f7ec722e52bb8812"),
		common.HexToAddress("5eb91a2604a4433221bcf0e60e2d1cc7b1c4cf05"),
		common.HexToAddress("6cc7a4bfb8d885ffe0a3561655519c992d6619cd"),
	}

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
		ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
		to := tos[t]
		value := big.NewInt(10000000000)
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
				for i := 0; i < 10; i++ {
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
					signedTx, err := ks.SignTxWithPassphrase(accounts.Account{Address: _from}, "i3nxx1rk", newTx, networkID)
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
