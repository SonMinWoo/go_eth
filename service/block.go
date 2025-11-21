package service

import (
	"block_chain/types"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBLock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {

		if err == mongo.ErrNoDocuments {
			s.log.Info("Genessis block will be created")
			//create genessis block
			genessisMessage := "This is genessis block"

			tx := createTransaction(genessisMessage, "0x50885a7528dab3a7bf09b1ec9e92c05865f5e2370b0e8d218d7cfbb45dac6913", "", "", 1)
			//Create new block
			newBlock := createBlockInner([]*types.Transaction{tx}, pHash, height)
			//마이닝
			pow := s.NewPow(newBlock)
			newBlock.Nonce, newBlock.Hash = pow.RunMining()

			return newBlock
		} else {
			s.log.Error("Failed to get latest block", "Err", err)
			panic(err)
		}
	} else {
		pHash = latestBlock.Hash

		newBlock := createBlockInner(txs, pHash, height)

		pow := s.NewPow(newBlock)

		newBlock.Nonce, newBlock.Hash = pow.RunMining()

		//create new block
	}
	return nil
}

func createBlockInner(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, to, amount string, block int64) *types.Transaction {
	data := struct {
		Message string `json:"message"`
		From    string `json:"from"`
		To      string `json:"to"`
		Amount  string `json:"amount"`
	}{
		Message: message,
		From:    from,
		To:      to,
		Amount:  amount,
	}

	dataToSign := fmt.Sprintf("%x\n", data)

	pk := "0x50885a7528dab3a7bf09b1ec9e92c05865f5e2370b0e8d218d7cfbb45dac6913"

	if ecdsaPrivatekKey, err := crypto.HexToECDSA(pk); err != nil {
		panic(err)
	} else if r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivatekKey, []byte(dataToSign)); err != nil {
		panic(err)
	} else {
		signiture := append(r.Bytes(), s.Bytes()...)

		return &types.Transaction{
			Block:   block,
			Time:    time.Now().Unix(),
			From:    from,
			To:      to,
			Amount:  amount,
			Message: message,
			Tx:      hex.EncodeToString(signiture),
		}
	}

}
