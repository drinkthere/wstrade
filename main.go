package main

import (
	"fmt"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
	"wstrade/client"
	"wstrade/config"
	"wstrade/context"
	"wstrade/utils"
	"wstrade/utils/logger"
)

var globalConfig config.Config
var globalBinanceClient client.BinanceClient
var globalContext context.GlobalContext

func ExitProcess() {
	os.Exit(1)
}

func main() {
	runtime.GOMAXPROCS(1)
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 监听退出消息，并调用ExitProcess进行处理
	utils.RegisterExitSignal(ExitProcess)

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 初始化okx客户端
	globalBinanceClient.Init(&globalConfig)

	// 解析config，加载杠杆和合约交易对，初始化context，账户初始化设置，拉取仓位、余额等
	globalContext.Init(&globalConfig, &globalBinanceClient)

	startTradeProcessor()

	//log.Println(http.ListenAndServe(":6060", nil))
	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
