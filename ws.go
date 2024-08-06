package main

import (
	"wstrade/trading"
	"wstrade/utils/logger"
)

//func startAccountWebSocket() {
//	// 监听Binance统一账户永续订单信息 和持仓信息
//	orderChan := make(chan *portfolio.WsOrderTradeUpdate)
//	accountChan := make(chan *portfolio.WsAccountUpdate)
//	message.StartBinanceUserDataWs(&globalBinanceClient, &globalConfig, &globalContext, orderChan, accountChan)
//
//	// 处理order update消息
//	message.StartGatherBinanceOrder(orderChan, &globalContext)
//	// 处理account update消息
//	message.StartGatherBinanceAccount(accountChan, &globalContext)
//}

func startTradeProcessor() {
	trading.StartBinanceTradingWebsocket(&globalConfig, &globalContext)
	logger.Info("[Trading] Start Trading Websocket")

	//trading.WaitingCancelOrder(&globalBinanceClient, &globalContext)
	//logger.Info("[Trading] Start Channel Of Cancel Order")
}
