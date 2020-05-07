package logic

import (
	"src/common"
	"src/ss_proto"
	"sync/atomic"

	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"

	. "src/pingsvr/config"
)

// 创建BusClientMgr 或者 BusServerMgr
type BusLogic struct {
	BusMgr *bus.BusClientMgr
}

func (self *BusLogic) Init() {
	cfg := &bus.Config{}
	cfg.BusCfg = make([]*bus.NewSvrInfo, 0)
	cfg.OnData = self.OnBusData
	cfg.OnNewBus = self.OnNewBus
	cfg.OnBusReg = self.OnBusReg

	cfg.SvrInfo = &bus.NewSvrInfo{}
	cfg.SvrInfo.Id = G_Pingcfg.ServerInfo.IId

	cfg.Parser = network.NewProtocParser(-1)
	cfg.Parser.Register(uint16(ss_proto.Cmd_SS_PING_PONG), ss_proto.PingPong{})
	cfg.Parser.Register(uint16(ss_proto.Cmd_SS_ROLE_LOAD), ss_proto.RoleLoad{})

	for _, v := range G_Pingcfg.Bussvrs {
		binfo := &bus.NewSvrInfo{}
		binfo.Id = bus.MakeServerIdByStr(v.ID)
		binfo.Ip = v.IP
		binfo.Port = v.Port
		cfg.BusCfg = append(cfg.BusCfg, binfo)
	}

	self.BusMgr = bus.NewBusClientMgr(cfg)
}

func (self *BusLogic) Close() {
	if self.BusMgr != nil {
		self.BusMgr.Close()
	}
}

func (self *BusLogic) OnBusReg(info *bus.NewSvrInfo, flag int64) {
	logger.Info("BusLogic.OnBusReg id:%v,s:%v,f:%v", info.Id, bus.GetServerIdStr(info.Id), flag)

	atomic.AddInt32(&G_OpenTest, 1)
}

func (self *BusLogic) OnNewBus(id int64) bool {
	logger.Info("BusLogic.OnNewBus id:%v,s:%v", id, bus.GetServerIdStr(id))

	fid := bus.GetServerFuncId(id)
	if fid != common.FUNC_PONG {
		return false
	}

	return true
}

func (self *BusLogic) OnBusData(msg *network.RawMessage) *network.RawMessage {
	logger.Debug("BusLogic:OnBusData conn:%v,msg:%v", 1, msg)
	switch msg.MsgId {
	case uint16(ss_proto.Cmd_SS_PING_PONG):
		// 多个Pong不是同步的
		///!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		Pong_BusRep()
	default:
		logger.Warning("BusLogic:OnBusData conn:%v,msg:%v", 1, msg)
	}

	return nil
}
