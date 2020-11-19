package logic

import (
	"fmt"

	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"
	"github.com/letterbaby/manzo/rand"
	"github.com/letterbaby/manzo/task"

	"src/common"
	roletbl "src/common/roletble"
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
			self.SendRouteMsg(msg)
		case uint16(ss_proto.Cmd_SS_ROLE_LOAD):
			self.Hand_RoleLoad(msg)
		default:
			logger.Warning("PongBus:OnBusData conn:%v,msg:%v", 1, msg)
		}

	} else {
		tid := rand.RandInt(0, 100)

		/*
			switch msg.MsgId {
			case uint16(ss_proto.Cmd_SS_PING_PONG):
				ping := &ss_proto.PingPong{}
				err := proto.Unmarshal(msg.MsgData.([]byte), ping)
				if err != nil {
					logger.Fatal("PongBus:OnBusData conn:%v,msg:%v", 1, msg)
				}
				tid = ping.Gid
			case uint16(ss_proto.Cmd_SS_ROLE_LOAD):
				load := msg.MsgData.(*ss_proto.RoleLoad)
				tid = load.RoleId
			}
		*/
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

	logger.Debug("PongBus:Hand_Task conn:%v,t:%v,msg:%v", self.Conn, t, msg)

	switch msg.MsgId {
	case uint16(ss_proto.Cmd_SS_PING_PONG):
		self.SendRouteMsg(msg)
	case uint16(ss_proto.Cmd_SS_ROLE_LOAD):
		self.Hand_RoleLoad(msg)
	default:
		logger.Warning("PongBus:Hand_Task conn:%v,msg:%v", 1, msg)
	}
	// !!!回收
	t.Msg = nil
	t.Bus = nil

	taskPool.Put(t)
}

func (self *PongBus) Hand_RoleLoad(msg *network.RawMessage) {
	logger.Info("PongBus:Hand_RoleLoad conn:%v,msg:%v", self.Conn, msg)

	req := msg.MsgData.(*ss_proto.RoleLoad)

	rid := fmt.Sprintf("%d", req.RoleId)

	// cache 不允许从

	var roledata *roletbl.RoleData
	ret, err := G_Redis.Get(false, "pong", "base_"+rid)
	if err != nil {
		if err != redis.ErrNil {
			logger.Error("PongBus:Hand_RoleLoad rid:%v,i:%v", req.RoleId, err)
		}
	} else if ret != nil {
		//TODO：通过固定更新频率来判断，缓存是否失效
		roledata = roletbl.CreateData(req.RoleId, "role", G_MysqlMgr)
		roledata.SetData(ret.([]byte))
	}

	if roledata == nil {
		roledata = roletbl.LoadData(req.RoleId, "role", true, G_MysqlMgr)
	}

	if roledata != nil {
		//也可以通过req.Base修改
		//roledata.SetData(req.Base)

		// 模拟修改
		lvl := roledata.GetLevel()
		roledata.SetLevel(lvl + 1)

		name := roledata.GetRoleName()
		if len(name) == 0 {
			roledata.SetRoleName(fmt.Sprintf("rr%d", req.RoleId))
		}

		// 生成2进制
		req.Base = roledata.GetData()

		// err 看上面的TODO
		err = G_Redis.SetEx("pong", "base_"+rid, 30*60, req.Base)
		if err != nil {
			logger.Error("PongBus:Hand_RoleLoad rid:%v,i:%v", req.RoleId, err)
		}

		// 在线存一般都是异步，下线保存一般都是同步
		roledata.Save(false, false)
	}

	self.SendRouteMsg(msg)
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
