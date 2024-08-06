package context

import (
	"github.com/dictxwang/go-binance/futures"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
	"wstrade/client"
	"wstrade/config"
)

type GlobalContext struct {
	TelegramBot *tgbotapi.BotAPI // 电报机器人

	TickerToCancelCh chan string
	PlaceOrderChan   chan *futures.WsPlaceOrder
	CancelOrderChan  chan *futures.WsCancelOrder
	rwLock           sync.RWMutex
}

func (context *GlobalContext) Init(globalConfig *config.Config, globalBinanceClient *client.BinanceClient) {
	context.TickerToCancelCh = make(chan string)
	context.PlaceOrderChan = make(chan *futures.WsPlaceOrder)
	context.CancelOrderChan = make(chan *futures.WsCancelOrder)
	context.rwLock = sync.RWMutex{}
}
