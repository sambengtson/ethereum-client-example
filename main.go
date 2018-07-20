package main

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onrik/ethrpc"
)

func main() {
	client := ethrpc.New("http://192.168.1.71:8545")
	_, err := client.Web3ClientVersion()
	if err != nil {
		panic(err)
	}

	address, privateKey := generateWallet()

	coinbase, err := client.EthCoinbase()
	if err != nil {
		panic(err)
	}

	/* Typically not necessary.  Only for coinbase account */
	_, err = client.Call("personal_unlockAccount", coinbase, "password")
	if err != nil {
		panic(err)
	}

	balance, err := client.EthGetBalance(coinbase, "latest")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Coinbase account balance: %v", balance.Div(&balance, ethrpc.Eth1()).String()))

	balance, err = client.EthGetBalance(address, "latest")
	if err != nil {
		panic(err)
	}


	txId, err := client.EthSendTransaction(ethrpc.T{
		From:  coinbase,
		To:    address,
		Value: ethrpc.Eth1().Mul(ethrpc.Eth1(), big.NewInt(20)),
	})

	for {
		receipt, err := client.EthGetTransactionReceipt(txId)
		if err != nil {
			panic(err)
		}

		fmt.Println("Waiting for transaction confirmation...")
		if receipt.Status != 0 {
			fmt.Println("Transaction confirmed!")
			break
		}

		time.Sleep(time.Second * 10)
	}

	tx := signAndSendBack(coinbase, address, privateKey)
	str := "0x" + hex.EncodeToString(tx)
	txId, err = client.EthSendRawTransaction(str)

	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func generateWallet() (string, string) {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())
	return address, privateKey
}

func signAndSendBack(to, from, privKey string) []byte {
	chainId := big.NewInt(1994)
	nonce := uint64(0)
	amount := ethrpc.Eth1().Mul(ethrpc.Eth1(), big.NewInt(10))
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(0)

	tx := types.NewTransaction(nonce, common.HexToAddress(to), amount, gasLimit, gasPrice, nil)

	senderPrivKey, _ := crypto.HexToECDSA(privKey)
	signer := types.NewEIP155Signer(chainId)
	signedTx, _ := types.SignTx(tx, signer, senderPrivKey)

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)

	return buff.Bytes()
}
