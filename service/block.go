package service

import (
	"block_chain/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBLock(txs []*types.Transaction, prevHash []byte, from string) *types.Block {

	var block *types.Block

	if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
		s.log.Error("Failed to get wallet by public key", "publicKey", from, "err", err)
		panic(err)
	} else {
		latestBlock, err := s.repository.GetLatestBlock()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.log.Info("Genessis block will be created")

				genessisMessage := "This is genessis block"

				tx := createTransaction(genessisMessage, from, wallet.PrivateKey, "", "", 1)

				block = createBlockInner([]*types.Transaction{tx}, "", 1)

			}
		} else {
			block = createBlockInner(txs, latestBlock.Hash, latestBlock.Height+1)
		}
		pow := s.NewPow(block)

		block.Nonce, block.Hash = pow.RunMining()
		s.log.Info("New block created", "block", block)

		if err := s.repository.SaveBlock(block); err != nil {
			s.log.Error("Failed to save block", "err", err)
			panic(err)
		} else {
			s.log.Info("Block saved successfully", "height", block.Height, "hash", fmt.Sprintf("%x", block.Hash))
			return block
		}
	}

}

func createBlockInner(txs []*types.Transaction, prevHash string, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         "",
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, pk, to, amount string, block int64) *types.Transaction {
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

func HashTransactions(b *types.Block) []byte {

	var txHashes [][]byte

	for _, tx := range b.Transactions {
		var encoded bytes.Buffer

		enc := gob.NewEncoder(&encoded)

		if err := enc.Encode(tx); err != nil {
			panic(err)
		} else {
			txHashes = append(txHashes, encoded.Bytes())
		}
	}

	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}
