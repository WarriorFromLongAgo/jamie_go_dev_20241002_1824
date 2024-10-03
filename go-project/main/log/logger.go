package log

import (
	"go-project/main/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go-project/util/file"
)

var (
	level   zapcore.Level
	options []zap.Option
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg *config.Configuration) (*ZapLogger, error) {
	createRootDir(cfg)
	setLogLevel(cfg)

	if cfg.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}
	zapLogger := zap.New(getZapCore(cfg), options...)
	return &ZapLogger{logger: zapLogger}, nil
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func createRootDir(cfg *config.Configuration) {
	if ok, _ := file.PathExists(cfg.Log.RootDir); !ok {
		_ = os.Mkdir(cfg.Log.RootDir, os.ModePerm)
	}
}

func setLogLevel(cfg *config.Configuration) {
	switch cfg.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

func getZapCore(cfg *config.Configuration) zapcore.Core {
	var encoder zapcore.Encoder

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(cfg.Server.Env + "." + l.String())
	}

	if cfg.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(cfg), level)
}

func getLogWriter(cfg *config.Configuration) zapcore.WriteSyncer {
	loggerFile := &lumberjack.Logger{
		Filename:   cfg.Log.RootDir + "/" + cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	}

	return zapcore.AddSync(loggerFile)
}
