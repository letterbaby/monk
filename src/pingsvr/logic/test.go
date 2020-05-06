package logic

import (
	"src/ss_proto"
	"sync"
	"sync/atomic"
	"time"

	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"
)

const (
	GO_COUNT  = 10
	MSG_COUNT = 1000
)

var (
	G______W sync.WaitGroup
)

func StartTest() {
	go func() {
		for atomic.LoadInt32(&G_OpenTest) == 0 {
			time.Sleep(time.Second * 1)
		}

		now := time.Now()
		logger.Info("------------StartTest start-------------")

		G______W.Add(GO_COUNT * MSG_COUNT * 2)
		for i := 0; i < GO_COUNT; i++ {
			go Test_Rpc()
			go Test_Rep(i)
		}
		G______W.Wait()

		logger.Info("------------StartTest done-------------:%v", time.Now().Sub(now))
	}()
}

func Test_Rpc() {
	for i := 0; i < MSG_COUNT; i++ {
		ping := &ss_proto.PingPong{}
		ping.Seq = int32(i)

		rmsg := &network.RawMessage{}
		rmsg.MsgId = uint16(ss_proto.Cmd_SS_PING_PONG)
		rmsg.MsgData = ping

		logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

		rt := G_BusLogic.BusMgr.SendRouteMsg(G_PongSvrId, false, rmsg, true, 2, G_PongSvrId, false)
		if rt == nil {
			logger.Fatal("Test id:%v,seq:%v", G_PongSvrId, i)
			return
		}

		pong := rt.MsgData.(*ss_proto.PingPong)
		if pong.Seq != ping.Seq {
			logger.Fatal("Test id:%v,ping:%v,pong:%v", ping, pong)
			return
		}

		G______W.Done()

		time.Sleep(time.Millisecond * 1)
	}
}

func Pong_Rep() {
	G______W.Done()
}

func Test_Rep(gid int) {
	for i := 0; i < MSG_COUNT; i++ {
		ping := &ss_proto.PingPong{}
		ping.Seq = int32(i)
		ping.Gid = int32(gid)

		rmsg := &network.RawMessage{}
		rmsg.MsgId = uint16(ss_proto.Cmd_SS_PING_PONG)
		rmsg.MsgData = ping

		logger.Debug("BusMgr.SendRouteMsg conn:%v,msg:%v", 1, rmsg)

		G_BusLogic.BusMgr.SendRouteMsg(G_PongSvrId, false, rmsg, false, 2, G_PongSvrId, false)

		time.Sleep(time.Millisecond * 1)
	}
}
