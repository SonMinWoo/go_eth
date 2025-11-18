package service

import (
	"block_chain/config"
	"block_chain/repository"
	"block_chain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"log/slog"
	"os"
)

type Service struct {
	config *config.Config

	repository *repository.Repository

	log *slog.Logger

	difficulty int64
}

func NewService(config *config.Config, repository *repository.Repository, difficulty int64) *Service {
	s := &Service{
		config:     config,
		repository: repository,
		log:        slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "app"),
	}

	return s
}

func (s *Service) newWallet() (string, string, error) {
	p256 := elliptic.P256()

	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else {
		if private == nil {
			return "", "", errors.New(types.PkNil)
		} else {

		}
	}
	return "", "", nil
}
