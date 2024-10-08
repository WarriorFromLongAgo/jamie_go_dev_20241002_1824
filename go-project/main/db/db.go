package db

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormLogger "gorm.io/gorm/logger"

	"go-project/main/config"
	log2 "go-project/main/log"
)

func InitializeDB(cfg *config.Configuration, log *log2.ZapLogger) *gorm.DB {
	switch cfg.MysqlDatabase.Driver {
	case "mysql":
		db := initMySqlGorm(cfg, log)
		return db
	default:
		db := initMySqlGorm(cfg, log)
		return db
	}
}

func initMySqlGorm(cfg *config.Configuration, log *log2.ZapLogger) *gorm.DB {
	dbConfig := cfg.MysqlDatabase

	if dbConfig.Database == "" {
		return nil
	}
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" +
		dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"

	log.Info("dsn " + dsn)

	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   getGormLogger(cfg),
	}); err != nil {
		log.Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		return db
	}
}

func getGormLogger(cfg *config.Configuration) ormLogger.Interface {
	var logMode ormLogger.LogLevel

	switch cfg.MysqlDatabase.LogMode {
	case "silent":
		logMode = ormLogger.Silent
	case "error":
		logMode = ormLogger.Error
	case "warn":
		logMode = ormLogger.Warn
	case "info":
		logMode = ormLogger.Info
	default:
		logMode = ormLogger.Info
	}

	return ormLogger.New(getGormLogWriter(cfg), ormLogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: false,
		Colorful:                  !cfg.MysqlDatabase.EnableFileLogWriter,
	})
}

func getGormLogWriter(cfg *config.Configuration) ormLogger.Writer {
	var writer io.Writer

	if cfg.MysqlDatabase.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   cfg.Log.RootDir + "/" + cfg.MysqlDatabase.LogFilename,
			MaxSize:    cfg.Log.MaxSize,
			MaxBackups: cfg.Log.MaxBackups,
			MaxAge:     cfg.Log.MaxAge,
			Compress:   cfg.Log.Compress,
		}
	} else {
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}
