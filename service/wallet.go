package service

import (
	"block_chain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *Service) newWallet() (string, string, error) {

	p256 := elliptic.P256()

	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else if private == nil {
		return "", "", fmt.Errorf("private key is nil")
	} else {
		privateKeyBytes := crypto.FromECDSA(private)
		hexutil.Encode(privateKeyBytes)
		fmt.Println(privateKeyBytes)
	}

	return "private_key_example", "public_key_example", nil
}

func (s *Service) MakeWallet() (*types.Wallet, error) {

	fmt.Println("들어옴")
	var wallet types.Wallet
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = s.newWallet(); err != nil {
		panic(err)
	} else {

		//todo - connect repository to store wallet info
		return &wallet, nil
	}
}
