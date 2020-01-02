package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var once sync.Once
var onceSugar sync.Once
var sugar *zap.SugaredLogger

//GetLogger get singleton logger
func GetLogger() *zap.Logger {
	once.Do(func() {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	})
	defer logger.Sync()
	return logger
}

//GetSugarLogger get singleton sugar logger
func GetSugarLogger() *zap.SugaredLogger {
	logger := GetLogger()
	onceSugar.Do(func() {
		sugar = logger.Sugar()
	})
	return sugar
}
