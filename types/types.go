package types

const (
	CreateWallet = "CreateWallet"
	TransferCoin = "TransferCoin"
	MintCoin     = "MintCoin"
)

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
