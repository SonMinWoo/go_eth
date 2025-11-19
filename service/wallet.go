package service

import (
	"block_chain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *Service) newKeyPair() (string, string, error) {

	p256 := elliptic.P256()

	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else if private == nil {
		return "", "", fmt.Errorf("private key is nil")
	} else {
		privateKeyBytes := crypto.FromECDSA(private)
		privateKey := hexutil.Encode(privateKeyBytes)

		//privateKey: 0xb9753351d4eb715b1bb6c7c734d94990f74642fca81d8f83d9ec16c04e47a64d
		againPrivateKey, err := crypto.HexToECDSA(privateKey[2:])
		if err != nil {
			return "", "", err
		}

		cPublicKey := againPrivateKey.Public()
		publicKeyECDSA, ok := cPublicKey.(*ecdsa.PublicKey)

		if !ok {
			return "", "", errors.New("error casting public key type")
		}

		publicKey := crypto.PubkeyToAddress(*publicKeyECDSA)

		// publicKey := private.Public()
		// publicKeyECDSA, ok := publicKey.(ecdsa.PublicKey)
		// if !ok {
		// 	return "", "", errors.New("error casting public key type")
		// }

		// address := crypto.PubkeyToAddress(*&publicKeyECDSA)
		//fmt.Println(privateKeyBytes)

		//[:] -> 대문자
		return privateKey, hexutil.Encode(publicKey[:]), nil
	}

}

func (s *Service) MakeWallet() *types.Wallet {

	var wallet types.Wallet
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = s.newKeyPair(); err != nil {
		panic(err)
	} else if err = s.repository.CreateNewWallet(&wallet); err != nil {
		//todo - connect repository to store wallet info

		s.repository.CreateNewWallet(&wallet)
		return nil
	} else {
		return &wallet
	}
}

func (s *Service) GetWallet(pk string) (*types.Wallet, error) {
	if wallet, err := s.repository.GetWallet(pk); err != nil {
		return nil, err
	} else {
		return wallet, nil
	}
}
