package client

import (
	"context"
	"github.com/dictxwang/go-binance/futures"
	"golang.org/x/time/rate"
	"time"
	"wstrade/config"
	"wstrade/utils/logger"
)

type BinanceClient struct {
	FuturesClient   *futures.Client
	FuturesWsClient *futures.ClientWs
	limiter         *rate.Limiter
	limitProcess    int
}

func (cli *BinanceClient) Init(cfg *config.Config) bool {
	ctx := context.Background()
	if cfg.LocalBinanceIP == "" {
		cli.FuturesClient = futures.NewClient(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
		cli.FuturesWsClient = futures.NewTradingWsClient(ctx, cfg.BinanceWsAPIKey, cfg.BinanceWsSecretKey, "")
	} else {
		cli.FuturesClient = futures.NewClientWithIP(cfg.BinanceAPIKey, cfg.BinanceSecretKey, cfg.LocalBinanceIP)
		cli.FuturesWsClient = futures.NewTradingWsClient(ctx, cfg.BinanceWsAPIKey, cfg.BinanceWsSecretKey, cfg.LocalBinanceIP)
	}

	limit := rate.Every(1 * time.Second / time.Duration(cfg.APILimit))
	cli.limiter = rate.NewLimiter(limit, 60)
	cli.limitProcess = cfg.LimitProcess
	return true
}

func (cli *BinanceClient) CheckLimit(n int) bool {
	if cli.limitProcess == 1 {
		err := cli.limiter.WaitN(context.Background(), n)
		if err != nil {
			logger.Error("[OkxClient] reach to limit, error:%s", err.Error())
		}
		return true
	}
	ret := cli.limiter.AllowN(time.Now(), n)
	if !ret {
		logger.Warn("[OkxClient] reach to limit")
	}
	return ret
}
