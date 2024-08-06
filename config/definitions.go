package config

type (
	RiskType       int
	Exchange       string
	InstrumentType string
	OrderSide      string
	OrderStatus    string
	TickerSource   string
)

const (
	// NoRisk 可以挂单，其他RiskType暂停挂单。1表示出错，2表示处于结算时间，3系统暂停等待价格更新
	NoRisk                = RiskType(iota)
	FatalErrorRisk        = NoRisk + 1
	SettlingRisk          = FatalErrorRisk + 1
	ExitRisk              = SettlingRisk + 1
	Price10sNotUpdateRisk = ExitRisk + 1
	FeeRateWrongRisk      = Price10sNotUpdateRisk + 1

	BinanceExchange = Exchange("Binance")
	OkxExchange     = Exchange("Okx")

	UnknownInstrument = InstrumentType("UNKNOWN")
	SpotInstrument    = InstrumentType("SPOT")
	FuturesInstrument = InstrumentType("FUTURES")

	BuyOrderSide  = OrderSide("buy")
	SellOrderSide = OrderSide("sell")

	OrderCreate          = OrderStatus("create")         // 发起创建请求，但不确定订单是否创建成功
	OrderLive            = OrderStatus("live")           // 新订单创建成功
	OrderCanceling       = OrderStatus("OrderCanceling") // 发起取消请求，但不确定订单是否取消成功
	OrderPartiallyFilled = OrderStatus("partially_filled")
	OrderFilled          = OrderStatus("filled")
	OrderCancel          = OrderStatus("canceled")

	RapidMarketTickerSource = TickerSource("RapidMarket")
	WebSocketTickerSource   = TickerSource("WebSocket")
)
