package trading

import (
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"github.com/gorilla/websocket"
	"time"
	"wstrade/config"
	mmcontext "wstrade/context"
	"wstrade/utils"
	"wstrade/utils/logger"
)

type BinanceTradingWebSocket struct {
	conn      *websocket.Conn
	isStopped bool
	stopChan  chan struct{}
}

func newBinanceTradingWebSocket() *BinanceTradingWebSocket {
	return &BinanceTradingWebSocket{
		conn:      nil,
		isStopped: true,
		stopChan:  make(chan struct{}),
	}
}

func StartBinanceTradingWebsocket(globalConfig *config.Config, ctxt *mmcontext.GlobalContext) {
	tradingWs := newBinanceTradingWebSocket()
	tradingWs.ConnectAndLogon(globalConfig)
}

func (tw *BinanceTradingWebSocket) handleTradingConnection(resp *binanceFutures.WsTradingConnectionResp) {
	logger.Info("[BinanceTradingWs] Trading Ws Connection Response: %+v", resp)
}

func (tw *BinanceTradingWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BinanceTradingWs] Binance Connection Handler Error And Reconnect Ws: %s", err.Error())
	tw.stopChan <- struct{}{}
	tw.isStopped = true
}

func (tw *BinanceTradingWebSocket) ConnectAndLogon(globalConfig *config.Config) {
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[BinanceTradingWs] Recovered from panic: %v", rc)
			}

			logger.Warn("[BinanceTradingWs] Create Trading Connection and Logon")
		}()
		for {
			if !tw.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			conn, _, stopChan, err := binanceFutures.WsTradingConnect(globalConfig.LocalBinanceIP, utils.GetClientOrderID(), tw.handleTradingConnection, tw.handleError)
			if err != nil {
				logger.Error("[BinanceTradingWs] Connect to Trading Websocket Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BinanceTradingWs] Successfully Connect to Trading Websocket")
			// 重置channel和时间
			tw.conn = conn
			tw.stopChan = stopChan
			tw.isStopped = false

			// logon

		}
	}()
}

//
//func WaitingPlaceOrder(binanceClient *client.BinanceClient, ctxt *mmcontext.GlobalContext) {
//	go func() {
//		defer func() {
//			logger.Warn("[PlaceOrder] Waiting Place Order Exited.")
//		}()
//
//		logger.Info("[PlaceOrder] Start Waiting Place Order.")
//		createOrderService := binanceClient.PortfolioClient.NewUmCreateOrderService()
//		for {
//			select {
//			case order := <-ctxt.PlaceOrderChan:
//				if order != nil && binanceClient.CheckLimit(1) {
//					price := strconv.FormatFloat(order.OrderPrice, 'f', order.Precision[0], 64)
//					quantity := strconv.FormatFloat(order.OrderVolume, 'f', order.Precision[1], 64)
//
//					resp, err := createOrderService.
//						Symbol(order.InstID).
//						Side(portfolio.SideType(strings.ToUpper(string(order.OrderSide)))).
//						Type(portfolio.OrderTypeLimit).
//						TimeInForce(portfolio.TimeInForceTypeGTX).
//						Quantity(quantity).
//						Price(price).
//						NewClientOrderID(order.ClientOrderID).Do(context.Background())
//
//					if err != nil {
//						logger.Error("[PlaceOrder] PlaceOrderGTX place %s %s %s %s@%s order error: %s",
//							order.InstID, order.OrderSide, order.ClientOrderID, quantity, price, err.Error())
//					} else {
//						logger.Info("[PlaceOrder] PlaceOrderGTX place %s %s %s %s@%s %d order Success",
//							order.InstID, order.OrderSide, order.ClientOrderID, quantity, price, resp.OrderID)
//					}
//				} else {
//					logger.Info("[PlaceOrder] PlaceOrderGTX order is nil or reach API Limit")
//				}
//			}
//		}
//	}()
//}

//
//func WaitingCancelOrder(binanceClient *client.BinanceClient, ctxt *mmcontext.GlobalContext) {
//	go func() {
//		defer func() {
//			logger.Warn("[CancelOrder] Waiting Cancel Order Exited.")
//		}()
//
//		logger.Info("[CancelOrder] Start Waiting Cancel Order.")
//		cancelOrderService := binanceClient.PortfolioClient.NewUmCancelOrderService()
//		for {
//			select {
//
//			case order := <-ctxt.CancelOrderChan:
//				if order != nil {
//					_, err := cancelOrderService.Symbol(order.InstID).ClientOrderID(order.ClientOrderID).Do(context.Background())
//					if err != nil {
//						logger.Error("[CancelOrder] Error Occur While Cancel Signal Order: order=%+v, error=%s", order, err.Error())
//					} else {
//						logger.Info("[CancelOrder] Cancel %s %s %s  order Success",
//							order.InstID, order.OrderSide, order.ClientOrderID)
//					}
//				}
//			}
//		}
//	}()
//}
