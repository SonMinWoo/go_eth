package app

import (
	"block_chain/config"
	"block_chain/repository"
	"block_chain/service"
	"log/slog"
	"os"
)

type App struct {
	service *service.Service

	config *config.Config

	repository *repository.Repository

	log *slog.Logger
}

func NewApp(config *config.Config) {
	a := &App{
		config: config,

		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "app"),
	}
	var err error

	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(config, a.repository)
	}

}
