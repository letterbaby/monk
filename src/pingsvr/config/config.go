package cfg

import (
	"encoding/json"
	"src/common"

	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/utils"
)

type bussvr struct {
	ID   string `json:"id"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type pingConfig struct {
	Bussvrs    []*bussvr          `json:"bussvrs"`
	ServerInfo *common.ServerInfo `json:"serverinfo"`

	Pongcnt int `json:"pongcnt"`
	Gocnt   int `json:"gocnt"`
	Msgcnt  int `json:"msgcnt"`
}

var (
	G_Pingcfg = &pingConfig{}
)

func InitConfig(fileName string) {
	// 解析文件标准文件
	data, err := utils.LoadFile(fileName)
	if err != nil {
		logger.Fatal("Config error:%v", err)
	}

	err = json.Unmarshal(data, G_Pingcfg)
	if err != nil {
		logger.Fatal("Config error:%v", err)
	}

	G_Pingcfg.ServerInfo.IId = bus.MakeServerIdByStr(G_Pingcfg.ServerInfo.Id)
}
