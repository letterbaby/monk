package logic

import (
	"sync"

	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"
	"github.com/letterbaby/manzo/task"

	"src/common"
	. "src/pongsvr/config"
	"src/ss_proto"
)

type PongBusTask struct {
	SN int8
	Id int64

	Bus *PongBus
	Msg *network.RawMessage
}

func (self *PongBusTask) GetTaskSn(cnt int32) int8 {
	self.SN = int8(self.Id % int64(cnt))
	return self.SN
}

func (self *PongBusTask) OnExcute() {
	self.Bus.Hand_Task(self)
}

func (self *PongBusTask) TimeOut() {
	//??
	logger.Warning("PongBusTask:TimeOut sn:%v,msg:%v", self.SN, self.Msg)
}

var (
	taskPool = &sync.Pool{
		New: func() interface{} {
			return &PongBusTask{}
		},
	}
)

type PongBus struct {
	bus.BusServer

	term *task.TaskerMgr
}

func NewPongBus() network.IAgent {
	pbus := &PongBus{}

	pbus.term = task.NewTaskerMgr("PongBus", 100)

	// 重载??
	pbus.OnClose = pbus.Hand_Close
	pbus.OnMessage = pbus.Hand_Message

	pbus.OnData = pbus.OnBusData
	pbus.OnDisconnect = pbus.OnBusDisconnect

	pbus.Initx(G_Pongcfg.TcpConfig, G_BusLogic.SvrMgr)

	return pbus
}

func (self *PongBus) OnBusDisconnect() {
	logger.Info("PongBus:OnDisconnect conn:%v,cid:%v", self.Conn, self.Id)

	self.term.Close()

	funcId := bus.GetServerFuncId(self.Id)

	if common.FUNC_PING == funcId {
		// 清理改ID下的数据
	}
}

func (self *PongBus) OnBusData(msg *network.RawMessage) *network.RawMessage {

	logger.Debug("PongBus:message conn:%v,msg:%v", self.Conn, msg)

	funcId := bus.GetServerFuncId(self.Id)
	if common.FUNC_PING == funcId {
		// 可以分ID处理协议
	}

	if !G_Pongcfg.Bustask {
		//可以直接处理
		switch msg.MsgId {
		case uint16(ss_proto.Cmd_SS_PING_PONG):
			//pong := msg.MsgData.(*ss_proto.PingPong)
			self.SendRouteMsg(msg)
			return nil
		default:
			logger.Warning("PongBus:OnBusData conn:%v,msg:%v", 1, msg)
		}

	} else {
		tid := int32(0)

		switch msg.MsgId {
		case uint16(ss_proto.Cmd_SS_PING_PONG):
			ping := msg.MsgData.(*ss_proto.PingPong)
			tid = ping.Gid
		}

		//task异步处理
		t := taskPool.Get().(*PongBusTask)
		t.Id = int64(tid)
		t.Bus = self
		t.Msg = msg
		self.term.AddTask(t)
	}

	return nil
}

func (self *PongBus) Hand_Task(t *PongBusTask) {
	msg := t.Msg

	logger.Info("PongBus:Hand_Task conn:%v,t:%v,msg:%v", self.Conn, t, msg)

	switch msg.MsgId {
	case uint16(ss_proto.Cmd_SS_PING_PONG):
		//pong := msg.MsgData.(*ss_proto.PingPong)
		self.SendRouteMsg(msg)
	default:
		logger.Warning("PongBus:Hand_Task conn:%v,msg:%v", 1, msg)
	}
	// !!!回收
	t.Msg = nil
	t.Bus = nil

	taskPool.Put(t)
}

//----------------------------------------------------------------------
// 创建BusClientMgr 或者 BusServerMgr

type BusLogic struct {
	SvrMgr *bus.BusServerMgr
}

func (self *BusLogic) Init() {
	cfg := &bus.Config{}
	cfg.Parser = G_Pongcfg.TcpConfig.Parser

	cfg.SvrInfo = &bus.NewSvrInfo{}
	cfg.SvrInfo.Id = G_Pongcfg.ServerInfo.IId
	self.SvrMgr = bus.NewBusServerMgr(cfg)
}

func (self *BusLogic) Close() {

	if self.SvrMgr != nil {
		self.SvrMgr.Close()
	}
}
