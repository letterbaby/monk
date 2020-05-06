package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	. "src/pingsvr/config"
	"src/pingsvr/logic"

	"github.com/letterbaby/manzo/logger"
	"github.com/letterbaby/manzo/signal"
	"github.com/letterbaby/manzo/sys/console"
)

func exit(s os.Signal) {
	logger.Info("Handle signal: %v", s)

	logic.Close()

	logger.Close()

	os.Exit(0)
}

func main() {
	os.Setenv("GOTRACEBACK", "crash")
	//参数
	logc := flag.String("log", "pinglog.json", "log config")
	scfg := flag.String("cfg", "ping.json", "server config")
	flag.Parse()

	logger.Start(*logc)

	// 初始化本服务器配置
	InitConfig(*scfg)

	logic.Main()

	title := fmt.Sprintf("Pingsvr id:%s(%d)", G_Pingcfg.ServerInfo.Id, G_Pingcfg.ServerInfo.IId)
	console.SetConsoleTitle(title)

	logic.StartTest()

	h := []os.Signal{syscall.SIGINT,
		syscall.SIGTERM}

	signal.Watch(h, exit)
}
