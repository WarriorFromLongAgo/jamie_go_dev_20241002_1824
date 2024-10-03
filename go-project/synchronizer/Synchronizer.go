package synchronizer

import (
	"context"
	"go-project/chain/eth"
	"go-project/database"
)

type Synchronizer struct {
}

func NewSynchronizer(db *database.DB, client eth.EthClient, shutdown context.CancelCauseFunc) (*Synchronizer, error) {
	//latestHeader, err := db.Blocks.LatestBlockHeader()

	return &Synchronizer{}, nil
}
