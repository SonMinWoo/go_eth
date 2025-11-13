package service

import (
	"block_chain/config"
	"block_chain/repository"
	"log/slog"
	"os"
)

type Service struct {
	config *config.Config

	repository *repository.Repository

	log *slog.Logger
}

func NewService(config *config.Config, repository *repository.Repository) *Service {
	s := &Service{
		config: config,

		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "app"),
	}

	return s
}
