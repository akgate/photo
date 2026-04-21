package main

import (
	"context"
	"log"

	"github.com/akgate/photo/internal/application"
	"github.com/akgate/photo/internal/infrastructure/persistence"
	"github.com/akgate/platform/pkg/db/pg"
	"github.com/akgate/platform/pkg/db/tx"
)

func main() {

	ctx := context.Background()

	db, err := pg.NewPgClient(ctx, "")

	if err != nil {
		panic(err)
	}

	err = db.DB().Ping(ctx)

	if err != nil {
		log.Fatalf("db.Ping(): %v", err)
	}

	txManager := tx.NewTxManager(db.DB())

	repo := persistence.NewPhotoRepository(db.DB())

	cp := application.NewCommandProcessor(repo, txManager)
	qp := application.NewQueryProcessor(repo)

	_, _ = cp, qp
}
