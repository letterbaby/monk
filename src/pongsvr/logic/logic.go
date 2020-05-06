package logic

import (
	"github.com/letterbaby/manzo/logger"
)

var (
	G_BusLogic = &BusLogic{}
)

func Close() {
	logger.Info("Waiting for Pongsvr close, please wait...")

	G_BusLogic.Close()
	logger.Info("Pongsvr busLogic shutdown.")

	logger.Info("Pongsvr closed.")
}

func Main() {
	G_BusLogic.Init()
}
