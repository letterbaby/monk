package logic

import (
	"src/common"

	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
)

var (
	G_BusLogic  = &BusLogic{}
	G_PongSvrId int64

	G_OpenTest = int32(0)
)

func Close() {
	logger.Info("Waiting for Pingsvr close, please wait...")

	G_BusLogic.Close()
	logger.Info("Pingsvr busLogic shutdown.")

	logger.Info("Pingsvr closed.")
}

func Main() {
	G_PongSvrId = bus.MakeServerId(0, common.FUNC_PONG, 0)

	G_BusLogic.Init()
}
