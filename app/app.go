package app

import (
	"block_chain/config"
	"block_chain/repository"
	"block_chain/service"
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
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
		a.service = service.NewService(config, a.repository, 1)

		a.log.Info("Module Started", "time", time.Now().Unix())

		sc := bufio.NewScanner(os.Stdin)

		useCase()

		for {
			sc.Scan()

			input := strings.Split(sc.Text(), " ")
			if err = a.inputValueAssesment(input); err != nil {
				a.log.Error("Failed to call cli", "err", err, "input", input)
				panic(err)
			}
		}
	}

}

func (a *App) inputValueAssesment(input []string) error {
	msg := errors.New("check Use Case")

	if len(input) == 0 {
		return msg
	} else {
		switch input[0] {
		case CreateWallet:
			fmt.Println("CreateWallet in Switch")
			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("failed to create")
			} else {
				fmt.Println("Success to create wallet")
			}
		case TransferCoin:
			fmt.Println("TransferCoin in Switch")
		case MintCoin:
			fmt.Println("MintCoin in Switch")
		default:
			return msg
		}
		fmt.Println()
	}

	return nil
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
