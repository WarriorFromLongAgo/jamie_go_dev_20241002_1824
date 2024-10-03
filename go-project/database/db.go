package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DB struct {
	gorm *gorm.DB
	BlockInfo
	TransactionInfo
	Management
	TokenInfo
	TokenTransferLog
	WorkflowInfo
	WorkflowApprove
}

func NewDB(dbConfig *config.Config) (*DB, error) {
	writer := logs.MyLogWriter()
	DbLogger := logger.New(
		log.New(writer, "\r\n", log.Ldate|log.Ltime|log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConfig.DbHost, dbConfig.DbName)
	if dbConfig.DbPort != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConfig.DbPort)
	}
	if dbConfig.DbUser != "" {
		dsn += fmt.Sprintf(" user=%s", dbConfig.DbUser)
	}
	if dbConfig.DbPassword != "" {
		dsn += fmt.Sprintf(" password=%s", dbConfig.DbPassword)
	}
	gormConfig := gorm.Config{
		Logger:                 DbLogger,
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return gorm, nil
	})
	if err != nil {
		return nil, err
	}
	db := &DB{
		gorm:              gorm,
		Blocks:            common.NewBlocksDB(gorm),
		ContractEvent:     event.NewContractEventsDB(gorm),
		ActivityInfoDB:    activity.NewActivityDB(gorm),
		ActivityInfoExtDB: activity.NewActivityInfoExtDB(gorm),
		TokenNftDB:        token_nft.NewTokenNftDB(gorm),
		AccountNftInfoDB:  account_nft_info.NewAccountNftInfoDB(gorm),
		DropInfoDB:        drop.NewDropInfoDB(gorm),
		BlockListener:     block_listener.NewBlockListenerDB(gorm),
	}
	return db, nil
}

func (db *DB) Transaction(fn func(db *DB) error) error {
	return db.gorm.Transaction(func(tx *gorm.DB) error {
		txDB := &DB{
			gorm:              tx,
			Blocks:            common.NewBlocksDB(tx),
			ContractEvent:     event.NewContractEventsDB(tx),
			ActivityInfoDB:    activity.NewActivityDB(tx),
			ActivityInfoExtDB: activity.NewActivityInfoExtDB(tx),
			TokenNftDB:        token_nft.NewTokenNftDB(tx),
			DropInfoDB:        drop.NewDropInfoDB(tx),
			BlockListener:     block_listener.NewBlockListenerDB(tx),
			AccountNftInfoDB:  account_nft_info.NewAccountNftInfoDB(tx),
		}
		return fn(txDB)
	})
}

func (db *DB) Close() error {
	sql, err := db.gorm.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
