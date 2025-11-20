package service

import (
	"block_chain/types"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBLock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genessis block will be created")
			//create genessis block

			//Create new block
			newBlock := createBlockInner(txs, pHash, height)
			//마이닝

			return newBlock
		} else {
			s.log.Error("Failed to get latest block", "Err", err)
			panic(err)
		}
	} else {
		pHash = latestBlock.Hash

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
