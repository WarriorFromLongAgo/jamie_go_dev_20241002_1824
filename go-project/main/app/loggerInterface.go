package app

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Log  Logger
	DB   *gorm.DB
	Conf *Configuration
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func SetLogger(l Logger) {
	Log = l
}
func SetDB(db *gorm.DB) {
	DB = db
}

func SetConfig(config *Configuration) {
	Conf = config
}
