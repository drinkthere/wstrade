package context

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
	"wstrade/client"
	"wstrade/config"
	"wstrade/container"
)

type GlobalContext struct {
	TelegramBot *tgbotapi.BotAPI // 电报机器人

	TickerToCancelCh chan string
	PlaceOrderChan   chan *container.Order
	CancelOrderChan  chan *container.Order
	rwLock           sync.RWMutex
}

func (context *GlobalContext) Init(globalConfig *config.Config, globalBinanceClient *client.BinanceClient) {
	context.TickerToCancelCh = make(chan string)
	context.PlaceOrderChan = make(chan *container.Order)
	context.CancelOrderChan = make(chan *container.Order)
	context.rwLock = sync.RWMutex{}
}
