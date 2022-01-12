package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"git.ont.io/waas/ERC721"
)

type Config struct {
	Network         string
	ContractAddress string
}

// 参考https://geth.ethereum.org/docs/dapp/native-bindings
func main() {
	conf := Config{
		Network:         "https://rinkeby.infura.io/v3/13f7ccb9852e48c99af2bcc47d8445f3",
		ContractAddress: "0xC3329E18f65dE7F07Aa88e2eEF5D36F9943E072F",
	}
	client, err := ethclient.Dial(conf.Network)
	if err != nil {
		fmt.Printf("Failed to connect to eth: %v", err)
		return
	}

	token, err := ERC721.NewERC721(common.HexToAddress(conf.ContractAddress), client)
	if err != nil {
		fmt.Printf("Failed to instantiate a Token contract: %v", err)
		return
	}

	privateKey, err := crypto.HexToECDSA("a604b62553868ca82e12efcc50193f6d27ed12392b910abe2609156b3076ec13")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("your address is: %v", fromAddress)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// auth, err := bind.NewTransactor(strings.NewReader(key), "my awesome super secret password")
	auth := bind.NewKeyedTransactor(privateKey)

	session := &ERC721.ERC721Session{
		Contract: token,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: uint64(10000000),
			GasPrice: gasPrice,
			Nonce:    big.NewInt(int64(nonce)),
		},
	}

	name, _ := session.Name()
	fmt.Printf("name: %v\n", name)
	// owner, _ := session.OwnerOf(big.NewInt(1))
	// fmt.Printf("owner: %v\n", owner)
	// totalSupply, _ := session.TotalSupply()
	// fmt.Printf("TotalSupply: %v\n", totalSupply)
	// mintTx, err := session.Mint(fromAddress)
	// if err != nil {
	// 	fmt.Printf("Failed to Mint: %v", err)
	// 	return
	// }
	// fmt.Printf("mintTx: %v\n", mintTx.Hash())
	tokenId := big.NewInt(1)
	owner, _ := session.OwnerOf(tokenId)
	fmt.Println("Pre-transfer tokenId[1] Owner:", owner)

	tx, err := session.TransferFrom(fromAddress, common.HexToAddress("0xCE6D5aD5B7Ca5AE121A17bb5296016d5B6b5e150"), tokenId)
	if err != nil {
		fmt.Printf("Failed to Transfer: %v", err)
		return
	}

	fmt.Printf("txHash: %v\n", tx.Hash())
}
