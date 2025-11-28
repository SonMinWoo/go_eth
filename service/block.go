package service

import (
	"block_chain/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBLock(from, to, value string) {
	var block *types.Block
	var tx *types.Transaction
	var toBalance string

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genessis block will be created")

			genessisMessage := "This is genessis block"

			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx = createTransaction(genessisMessage, common.Address{}.String(), pk, to, value, 1)
				block = createBlockInner([]*types.Transaction{tx}, "", 1)
			}
		}
	} else {

		if common.HexToAddress(from) == (common.Address{}) {
			//Mint
			// todo ->private key 가져오기
			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx = createTransaction("MintCoin", common.Address{}.String(), pk, to, value, 1)
				toBalance = value
			}

		} else {
			//Transfer
			if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
				s.log.Error("Failed to get wallet by public key", "publicKey", from, "err", err)
				panic(err)
			} else if toWallet, err := s.repository.GetWalletByPublicKey(to); err != nil {
				if err == mongo.ErrNoDocuments {
					s.log.Error("Can't find new wallet", to, "err", err)
				} else {
					s.log.Error("Failed to get to wallet by public key", "publicKey", to, "err", err)
					panic(err)
				}
			} else {
				// todo -> from 밸런스 ㅔ크
				fromDecimalBalance, _ := decimal.NewFromString(wallet.Balance)
				valueDecimal, _ := decimal.NewFromString(value)
				toDecimalBalance, _ := decimal.NewFromString(toWallet.Balance)

				//wallet.Balance < value -> error
				if fromDecimalBalance.Cmp(valueDecimal) == -1 {
					s.log.Info("Insufficient balance", "balance", wallet.Balance, "transferAmount", value)
					return
				} else {
					toBalance = value
					toDecimalBalance = toDecimalBalance.Add(valueDecimal)
					toBalance = toDecimalBalance.String()

					fromDecimalBalance = fromDecimalBalance.Sub(valueDecimal)
					value = fromDecimalBalance.String()
				}

				//todo ->전송순서 체크
				// if err := s.repository.CreateNewWallet(wallet); err != nil {
				// 	s.log.Error("Failed to update wallet balance", "err", err)
				// 	panic(err)
				// }
				//todo -> update wallet balance
				tx = createTransaction("TransferCoin", from, wallet.PrivateKey, to, value, 1)
			}
		}
		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)

	}
	pow := s.NewPow(block)
	block.Nonce, block.Hash = pow.RunMining()

	if err := s.repository.UpsertWalletsWhenTransfer(from, to, value, toBalance); err != nil {
		s.log.Error("Failed to upsert when transfer", "to", to, "value", value, "err", err)
		panic(err)
	} else if err := s.repository.SaveBlock(block); err != nil {
		s.log.Error("Failed to save block", "err", err)
		panic(err)
	} else {
		s.log.Info("Block saved successfully", "height", block.Height, "hash", fmt.Sprintf("%x", block.Hash))
	}

	// var block *types.Block
	// var err error

	// block, err = s.repository.GetLatestBlock()

	// if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
	// 	s.log.Error("Failed to get wallet by public key", "publicKey", from, "err", err)
	// 	panic(err)
	// } else {
	// 	latestBlock, err := s.repository.GetLatestBlock()
	// 	var block *types.Block
	// 	if err != nil {
	// 		if err == mongo.ErrNoDocuments {
	// 			s.log.Info("Genessis block will be created")

	// 			genessisMessage := "This is genessis block"

	// 			tx := createTransaction(genessisMessage, from, wallet.PrivateKey, to, value, 1)

	// 			block = createBlockInner([]*types.Transaction{tx}, "", 1)

	// 		}
	// 	} else {
	// 		tx := createTransaction("New Block is created", from, wallet.PrivateKey, to, value, latestBlock.Height+1)
	// 		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)
	// 	}
	// 	pow := s.NewPow(block)

	// 	block.Nonce, block.Hash = pow.RunMining()
	// 	s.log.Info("New block created", "block", block)

	// 	if err := s.repository.SaveBlock(block); err != nil {
	// 		s.log.Error("Failed to save block", "err", err)
	// 		panic(err)
	// 	} else {
	// 		s.log.Info("Block saved successfully", "height", block.Height, "hash", fmt.Sprintf("%x", block.Hash))
	// 		return block
	// 	}
	// }

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

	pk = strings.TrimPrefix(pk, "0x")

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
