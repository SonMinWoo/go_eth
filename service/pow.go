package service

import (
	"block_chain/types"
	"math/big"
)

type PowWork struct {
	Block      *types.Block `json:"block"`
	Target     *big.Int     `json:"target"`
	Difficulty *int64       `json:"difficulty"`
}

func (s *Service) NewPow(b *types.Block) *PowWork {
	t := new(big.Int).SetInt64(1)

	//1이 들어가면 2(10), 2가 들어가면 4(100)... -> 비트 밀기(비트마스크)
	t.Lsh(t, uint(256-s.difficulty))
	return &PowWork{Block: b, Target: t, Difficulty: &s.difficulty}
}
