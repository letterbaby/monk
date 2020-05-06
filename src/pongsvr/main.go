package main

import (
	"flag"
	"fmt"
	"os"
	"src/pongsvr/logic"
	"syscall"

	_ "net/http/pprof"

	. "src/pongsvr/config"

	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/network"
	"github.com/letterbaby/manzo/signal"
	"github.com/letterbaby/manzo/sys/console"

	"src/ss_proto"
)

var (
	G_TcpServer *network.TcpServer
)

func InitNetConfig() {

	// 其他没通过配置的需要初始化
	parser := network.NewProtocParser(-1)
	parser.Register(uint16(ss_proto.Cmd_SS_PING_PONG), ss_proto.PingPong{})

	G_Pongcfg.TcpConfig.Parser = parser
	G_Pongcfg.TcpConfig.Agent = logic.NewPongBus
}

func exit(s os.Signal) {
	logger.Info("Handle signal: %v", s)

	// 服务器内部先关闭逻辑，这样会等待客户端先退出
	logic.Close()

	// 当这个服务对外的话，请放在logic.Close之前
	if G_TcpServer != nil {
		G_TcpServer.Close()
	}

	logger.Close()

	os.Exit(0)
}

func main() {
	//参数
	logc := flag.String("log", "ponglog.json", "log config")
	scfg := flag.String("cfg", "pong.json", "server config")
	flag.Parse()

	logger.Start(*logc)

	// 初始化本服务器配置
	InitConfig(*scfg)
	InitNetConfig()

	logic.Main()

	//启动TCP服务器
	G_TcpServer = network.NewTcpServer(G_Pongcfg.TcpConfig)
	if G_TcpServer == nil {
		return
	}
	G_TcpServer.Serve(false)

	logger.Info("Main tcpconfig:%v,serverinfo:%v", G_Pongcfg.TcpConfig, G_Pongcfg.ServerInfo)

	title := fmt.Sprintf("Pongsvr id:%s(%d)", G_Pongcfg.ServerInfo.Id, G_Pongcfg.ServerInfo.IId)
	console.SetConsoleTitle(title)

	h := []os.Signal{syscall.SIGINT,
		syscall.SIGTERM}

	signal.Watch(h, exit)
}
