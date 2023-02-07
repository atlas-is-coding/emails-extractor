package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerInstance *zap.Logger
	once           sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		loggerInstance = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.InfoLevel,
		))
	})
	return loggerInstance
}
