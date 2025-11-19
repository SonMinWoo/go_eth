package repository

import (
	"block_chain/types"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateNewWallet(w *types.Wallet) error {

	ctx := context.Background()
	w.Time = uint64(time.Now().Unix())
	opt := options.Update().SetUpsert(true)

	fmt.Printf("r.wallet: %#v\n", r.wallet)
	fmt.Printf("ctx is nil? %v\n", ctx == nil)

	filter := bson.M{"privateKey": w.PrivateKey}
	update := bson.M{"$set": bson.M{
		"privateKey": w.PrivateKey,
		"publicKey":  w.PublicKey,
		"time":       w.Time,
	}}

	if res, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	} else {

		r.log.Info("wallet upsert result",
			"matched", res.MatchedCount,
			"modified", res.ModifiedCount,
			"upserted", res.UpsertedCount,
			"upsertedID", res.UpsertedID,
		)
		return nil
	}

}

func (r *Repository) GetWallet(pk string) (*types.Wallet, error) {
	ctx := context.Background()

	var wallet types.Wallet

	filter := bson.M{"privateKey": pk}

	if err := r.wallet.FindOne(ctx, filter, options.FindOne()).Decode(&wallet); err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}
