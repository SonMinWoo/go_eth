package app

import (
	"block_chain/config"
	"block_chain/repository"
	"block_chain/service"
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"time"

	. "block_chain/types"
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

		a.log.Info("Module Started", "time", time.Now().Unix())

		sc := bufio.NewScanner(os.Stdin)

		useCase()

		for {
			sc.Scan()
			fmt.Println(sc.Text())
		}
	}

}

func useCase() {
	fmt.Println()

	fmt.Println("THos os Akaps Module For BLockChain Core With Mongo")
	fmt.Println()
	fmt.Println("Use Case")

	fmt.Println("1. ", CreateWallet)
	fmt.Println("2. ", TransferCoin, "<To> <Amount>")
	fmt.Println("3. ", MintCoin, "<To><Amount>")
	fmt.Println()
}
