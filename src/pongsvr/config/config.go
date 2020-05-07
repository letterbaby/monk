package config

import (
	"encoding/json"
	"src/common"

	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/mysql"
	"github.com/letterbaby/manzo/network"
	"github.com/letterbaby/manzo/redis"
	"github.com/letterbaby/manzo/utils"
)

type pongConfig struct {
	ServerInfo *common.ServerInfo `json:"serverinfo"`

	Bustask   bool            `json:"bustask"`
	TcpConfig *network.Config `json:"tcpconfig"`

	RedisConfig *redis.Config `json:"redisconfig"`
	MysqlConfig *mysql.Config `json:"mysqlconfig"`
}

var (
	G_Pongcfg = &pongConfig{}
)

func InitConfig(fileName string) {
	// 解析文件标准文件
	data, err := utils.LoadFile(fileName)
	if err != nil {
		logger.Fatal("InitConfig file:%v,msg:%v", fileName, err)
	}

	err = json.Unmarshal(data, G_Pongcfg)
	if err != nil {
		logger.Fatal("InitConfig file:%v,msg:%v", fileName, err)
	}
	G_Pongcfg.ServerInfo.IId = bus.MakeServerIdByStr(G_Pongcfg.ServerInfo.Id)
}
