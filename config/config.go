package config

import (
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
)

type InstIDConfig struct {
	ContractNum            float64 // 每单委托数量（单位：张）
	VolPerCont             float64 // 每张对应币的数量
	BaseAsset              string  // eg. Base Asset: BTC
	Leverage               int     // 杠杆倍数，初始化时给交易对设置好
	MaxContractNum         float64 // 最大可以开的张数， 也是用来限制单方向最大持仓数量的
	Precision              [2]int  // BTCUSDT => [4, 2] amount数量精度， price数量精度
	EffectiveNum           float64 // 获取交易对报价时，quantity 需要大于这个值才认为有效（特别是从depth消息中获取价格时）
	FirstOrderMargin       float64 // 第一单不是根据gapSize来订单的，而是在min(对照数据)的基础上，再*/这个margin来计算的。
	FirstOrderRangePercent float64 // 挂第一单的时候，与当前订单的距离必须大于这个比例
	GapSizePercent         float64 // 每单之间的默认间隔，e.g. 0.0002 就是间隔万分之二
	ForgivePercent         float64 // 套利价差比 > forgive才下单，forgive会随着spread的变化而变动
	MinimumTickerShift     float64 // tickershift 只有当|position| > MinimumTickerShift && 有盈利时才开始应用
	MaximumTickerShift     float64 // tickershift 只有当|position| > MaximumTickerShift 就开始应用（不考虑是否盈利）
	BreakEvenX             float64 // 当前币价与盈亏平衡价比较时用的的参数，只有盈利才会启动tickerShift
	PositionReduceFactor   float64 // 减仓时使用x:(1 + 0.1 * (position / x)) * ContractNum
	TickerShift            float64 // 根据仓位修正现货和U本位合约买卖价格时的系数
	MaxOrderNum            int     // 每个方向上挂单的数量
	FarOrderNum            int     // 单边挂单最多限制
	VolatilityE            float64 // Math.pow(volatility / 100, volatilityE) / volatilityD
	VolatilityD            float64 // Math.pow(volatility / 100, volatilityE) / volatilityD
	VolatilityG            float64 // 计算adjustedGapSize = gapSize + gapSize * volatility * VolatilityG
}

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// 电报配置
	TgBotToken string
	TgChatID   int64

	// 币安配置
	BinanceAPIKey    string
	BinanceSecretKey string
	// 币安配置
	BinanceWsAPIKey    string
	BinanceWsSecretKey string

	// Okx配置
	OkxAPIKey    string
	OkxSecretKey string
	OkxPassword  string

	// 交易手续费，正数交易所给补贴，负数要扣我们手续费
	FeeRate float64

	// 频率控制
	APIRateInterval int // 时间间隔，10就是将1s拆成10份，时间间隔就是100ms
	APILimit        int // 每份允许通过的事件数量，10就是100ms内可以发起10次请求
	LimitProcess    int // 超过限制请求的处理方法，1等待，0是丢弃该次请求

	UpdateOrderIntervalMs int64 // 挂单的频率，单位ms

	// 套利配置
	VolatilityTimeMs int64                   // 计算Volatility需要的数据时长
	LocalOkxIP       string                  // OKX本地出口IP, 如果设置了就用这个IP来连接dex，不设置就用默认IP
	LocalBinanceIP   string                  // Binance本地出口IP, 如果设置了就用这个IP来连接dex，不设置就用默认IP
	TickerSource     TickerSource            // Ticker的消息源：WebSocket
	UseDepth         bool                    // 是否使用Okx Depth数据来更新Ticker（Depth数据本身也Merge了Ticker数据)
	DepthIPC         string                  // Okx Depth数据的IPC
	Colo             bool                    // 是否是co-location部署
	InstIDs          []string                // 要套利的交易对
	InstIDConfigs    map[string]InstIDConfig // 交易对的详细配置

	MaxErrorsPerMinute int64   // 每分钟允许出现的 Error 日志数量（超出数量之后退出程序）
	MinDeltaRate       float64 // 最小差比例， 价格变动超过这个才进行处理
	MinAccuracy        float64 // 价格最小精度

	// 订单撤销频率
	FarOrderCancelFreq time.Duration // 撤销较远订单的频率
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
