package logic

import (
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/mysql"
	"github.com/letterbaby/manzo/redis"

	roletbl "src/common/roletble"
	. "src/pongsvr/config"
)

var (
	G_BusLogic = &BusLogic{}

	G_MysqlOps *mysql.MyOps
	G_MysqlMgr *mysql.DBMgr

	G_Redis redis.IRedis
)

func Close() {
	logger.Info("Waiting for Pongsvr close, please wait...")

	G_BusLogic.Close()
	logger.Info("Pongsvr busLogic shutdown.")

	G_MysqlOps.Close()
	logger.Info("Pongsvr dbops shutdown.")

	G_MysqlMgr.Close()
	logger.Info("Pongsvr dbmgr shutdown.")

	G_Redis.Close()
	logger.Info("Pongsvr redis shutdown.")

	logger.Info("Pongsvr closed.")
}

func Main() {
	G_MysqlMgr = mysql.NewDBMgr(G_Pongcfg.MysqlConfig)
	G_MysqlOps = mysql.NewMyOps(G_Pongcfg.MysqlConfig)

	// 分表的时候 第一个参数有用
	roletbl.Init("", G_MysqlOps)

	G_Redis = redis.NewRedisCluster(G_Pongcfg.RedisConfig)

	G_BusLogic.Init()
}
