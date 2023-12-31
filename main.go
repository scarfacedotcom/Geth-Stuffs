package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

const (
	SK   = "0xaf5ead4413ff4b78bc94191a2926ae9ccbec86ce099d65aaf469e9eb1a0fa87f"
	ADDR = "0x6177843db3138ae69679A54b95cf345ED759450d"
)

func sendTransaction(cl *ethclient.Client) error {
	var (
		sk       = crypto.ToECDSAUnsafe(common.FromHex(SK))
		to       = common.HexToAddress("0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D")
		value    = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
		sender   = common.HexToAddress(ADDR)
		gasLimit = uint64(21000)
	)
	// Retrieve the chainid (needed for signer)
	chainid, err := cl.ChainID(context.Background())
	if err != nil {
		return err
	}
	// Retrieve the pending nonce
	nonce, err := cl.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return err
	}
	// Get suggested gas price
	tipCap, _ := cl.SuggestGasTipCap(context.Background())
	feeCap, _ := cl.SuggestGasPrice(context.Background())
	// Create a new transaction
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainid,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &to,
			Value:     value,
			Data:      nil,
		})
	// Sign the transaction using our keys
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainid), sk)
	// Send the transaction to our node
	return cl.SendTransaction(context.Background(), signedTx)
}
func main() {
	// create instance of ethclient and assign to cl
	client, err := ethclient.Dial("/tmp/geth.ipc")
	if err != nil {
		panic(err)
	}
	chainid, _ := client.ChainID(context.Background())
	if err != nil {
		return
	}
	fmt.Println("Chain ID", chainid.String())

	addr := common.HexToAddress("0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D")
	nonce, err := client.NonceAt(context.Background(), addr, big.NewInt(14000000))
	fmt.Println("Nonce", nonce)

	err = sendTransaction(client)
	if err != nil {
		return
	}
	fmt.Println("Transaction sent sucessfully")
}
