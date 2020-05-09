package logic

import (
	"src/common"
	roletbl "src/common/roletble"
	"src/ss_proto"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/letterbaby/manzo/bus"
	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"

	. "src/pingsvr/config"
)

var (
	G______R0 *roletbl.Tbl // 玩家0校验
	G______W  sync.WaitGroup
)

func StartTest() {
	go func() {
		for atomic.LoadInt32(&G_OpenTest) == 0 {
			time.Sleep(time.Second * 1)
		}

		now := time.Now()
		logger.Info("------------StartTest start-------------")

		G______W.Add(G_Pingcfg.Gocnt * G_Pingcfg.Msgcnt * 3)

		Test_DBTest(true)

		for i := 0; i < G_Pingcfg.Gocnt; i++ {
			go Test_BusRpc(i)
			go Test_BusRep(i)
			go Test_DBRpc(i)
		}
		G______W.Wait()

		Test_DBTest(false)

		logger.Info("------------StartTest done-------------:%v", time.Now().Sub(now))
	}()
}

func Test_BusRpc(gid int) {
	for i := 0; i < G_Pingcfg.Msgcnt; i++ {
		ping := &ss_proto.PingPong{}
		ping.Seq = int32(i)
		ping.Gid = int64(gid)

		rmsg := &network.RawMessage{}
		rmsg.MsgId = uint16(ss_proto.Cmd_SS_PING_PONG)
		rmsg.MsgData = ping

		logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

		rt := G_BusLogic.BusMgr.SendRouteMsg(G_PongSvrId, false, rmsg, true, 20, G_PongSvrId, false)
		if rt == nil {
			logger.Fatal("Test id:%v,seq:%v", G_PongSvrId, i)
			return
		}

		pong := &ss_proto.PingPong{}
		err := proto.Unmarshal(rt.MsgData.([]byte), pong)
		if err != nil {
			logger.Fatal("Test id:%v,seq:%v,i:%v", G_PongSvrId, i, err)
		}

		if pong.Seq != ping.Seq {
			logger.Fatal("Test id:%v,ping:%v,pong:%v", ping, pong)
			return
		}

		G______W.Done()

		//time.Sleep(time.Millisecond * 1)
	}
}

func Pong_BusRep() {
	G______W.Done()
}

func Test_BusRep(gid int) {
	for i := 0; i < G_Pingcfg.Msgcnt; i++ {
		ping := &ss_proto.PingPong{}
		ping.Seq = int32(i)
		ping.Gid = int64(gid)

		rmsg := &network.RawMessage{}
		rmsg.MsgId = uint16(ss_proto.Cmd_SS_PING_PONG)
		rmsg.MsgData = ping

		logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

		G_BusLogic.BusMgr.SendRouteMsg(G_PongSvrId, false, rmsg, false, 20, G_PongSvrId, false)

		//time.Sleep(time.Millisecond * 1)
	}
}

func Test_DBRpc(gid int) {
	for i := 0; i < G_Pingcfg.Msgcnt; i++ {
		load := &ss_proto.RoleLoad{}
		load.RoleId = int64(i)

		rmsg := &network.RawMessage{}
		rmsg.MsgId = uint16(ss_proto.Cmd_SS_ROLE_LOAD)
		rmsg.MsgData = load

		logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

		/// ！！！！数据库测试时，开多个pong，同一个角色在同一个pong操作
		lid := i%G_Pingcfg.Pongcnt + 1
		pongSvrId := bus.MakeServerId(0, common.FUNC_PONG, int64(lid))
		rt := G_BusLogic.BusMgr.SendRouteMsg(pongSvrId, false, rmsg, true, 20, pongSvrId, false)
		if rt == nil {
			logger.Fatal("Test id:%v,rid:%v", pongSvrId, i)
			return
		}

		data := rt.MsgData.(*ss_proto.RoleLoad)
		if data.Base != nil {
			base := &roletbl.Tbl{}
			proto.Unmarshal(data.Base, base)

			logger.Debug("Cmd_SS_ROLE_LOAD rid:%v,base:%v", i, base)
		} else {
			logger.Fatal("Test id:%v,roleId:%v", pongSvrId, i)
		}

		G______W.Done()

		//time.Sleep(time.Millisecond * 1)
	}
}

func Test_DBTest(start bool) {
	i := 0

	load := &ss_proto.RoleLoad{}
	load.RoleId = int64(i)

	rmsg := &network.RawMessage{}
	rmsg.MsgId = uint16(ss_proto.Cmd_SS_ROLE_LOAD)
	rmsg.MsgData = load

	logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

	lid := i%G_Pingcfg.Pongcnt + 1
	pongSvrId := bus.MakeServerId(0, common.FUNC_PONG, int64(lid))
	rt := G_BusLogic.BusMgr.SendRouteMsg(pongSvrId, false, rmsg, true, 20, pongSvrId, false)
	if rt == nil {
		logger.Fatal("Test id:%v,rid:%v", pongSvrId, i)
		return
	}

	data := rt.MsgData.(*ss_proto.RoleLoad)
	if data.Base != nil {
		base := &roletbl.Tbl{}
		proto.Unmarshal(data.Base, base)

		if start {
			G______R0 = base
		} else {
			// +1 是Test_DBTest的
			if G______R0.GetLevel()+int32(G_Pingcfg.Gocnt)+1 != base.GetLevel() {
				logger.Fatal("Cmd_SS_ROLE_LOAD rid:%v,obase:%v,nbase:%v", i, G______R0, base)
			}
			logger.Debug("Cmd_SS_ROLE_LOAD rid:%v,base:%v", i, base)
		}

	} else {
		logger.Fatal("Test id:%v,roleId:%v", G_PongSvrId, i)
	}
}
