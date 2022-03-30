package initialize

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panic("Error Failed to initializing logger: ", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
