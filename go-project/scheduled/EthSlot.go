package scheduled

import (
	"context"

	"gorm.io/gorm"

	"go-project/chain/eth"
	"go-project/main/log"
)

type EthSlot struct {
	ctx       context.Context
	ethClient eth.EthClient
	db        *gorm.DB
	log       *log.ZapLogger
}

func NewEthSlot(ctx context.Context, client eth.EthClient, db *gorm.DB, log *log.ZapLogger) (*EthSlot, error) {
	return &EthSlot{
		ctx:       ctx,
		ethClient: client,
		db:        db,
		log:       log,
	}, nil
}

// slot 32 64
// change scan all block to LatestFinalizedBlockHeader
