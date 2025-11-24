package app

import (
	"block_chain/config"
	"block_chain/global"
	"block_chain/repository"
	"block_chain/service"
	. "block_chain/types"
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	service *service.Service

	config *config.Config

	repository *repository.Repository

	log *slog.Logger
}

func NewApp(config *config.Config, difficulty int64) {
	a := &App{
		config: config,

		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "app"),
	}
	var err error

	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(config, a.repository, difficulty)

		a.log.Info("Module Started", "time", time.Now().Unix())

		sc := bufio.NewScanner(os.Stdin)

		useCase()

		for {
			from := global.FROM()

			if from != "" {
				a.log.Info("CUrrent Connected wallet", "from", from)
				fmt.Println()
			}

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
		from := global.FROM()

		switch input[0] {
		case CreateWallet:
			fmt.Println("CreateWallet in Switch")
			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("failed to create")
			} else {
				fmt.Println()
				a.log.Info("Success to create wallet", "pk", wallet.PrivateKey, "pu", wallet.PublicKey)
				fmt.Println()
			}
		case TransferCoin:
			if from == "" {
				a.log.Info("Not Connected")
				fmt.Println()
				return nil
			} else {
				a.service.CreateBLock([]*Transaction{}, []byte{}, from)
			}

		case MintCoin:
			if from == "" {
				a.log.Debug("Not Connected")
				fmt.Println()
				return nil
			} else {

			}
		case ConnectionWallet:
			fmt.Println(from)
			if from != "" {
				a.log.Info("Already Connected", "from", from)
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("Failed to find wallet PK is Nil", "pk", input[1])
					} else {
						a.log.Error("Failed to find wallet", "pk", input[1], "err", err)
					}
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("Success to connect wallet", "From", wallet.PublicKey)
					fmt.Println()
				}
			}
		case ChangeWallet:
			if from == "" {
				a.log.Debug("Not Connected")
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("Failed to find wallet PK is Nil", "pk", input[1])
					} else {
						a.log.Error("Failed to find wallet", "pk", input[1], "err", err)
					}
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("Success to connect wallet", "From", wallet.PublicKey)
					fmt.Println()
				}
			}
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
	fmt.Println("2. ", ConnectionWallet, " <PK>")
	fmt.Println("3. ", ChangeWallet, " <PK>")
	fmt.Println("4. ", TransferCoin, "<To> <Amount>")
	fmt.Println("5. ", MintCoin, "<To><Amount>")
	fmt.Println()
}
