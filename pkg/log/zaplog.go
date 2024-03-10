package log

import (
	"log"

	"go.uber.org/zap"
)

var (
	Log = NewZaLop()
)

func NewZaLop() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()
	return sugar
}

func Error(msg string) (*zap.SugaredLogger, string) {
	Log.Error("this is an error message")
	return Log, msg
}
