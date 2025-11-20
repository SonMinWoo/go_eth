package repository

import (
	"block_chain/config"
	"context"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	wallet *mongo.Collection
	tx     *mongo.Collection
	config *config.Config
	block  *mongo.Collection
	log    *slog.Logger
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{
		config: config,

		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "app"),
	}

	var err error
	ctx := context.Background()

	mConfig := config.Mongo
	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mConfig.Uri)); err != nil {
		r.log.Error("failed to connect to mongo", "uri", mConfig.Uri)
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		r.log.Error("failed to ping to mongo", "uri", mConfig.Uri)
		return nil, err
	} else {
		db := r.client.Database(mConfig.DB, nil)

		r.wallet = db.Collection("wallet")
		r.tx = db.Collection("tx")
		r.block = db.Collection("block")

		r.log.Info("success to connect Repository", "uri", mConfig.Uri, "db", mConfig.DB)
	}

	return r, nil
}
