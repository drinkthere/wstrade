package trading

import (
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"time"
	"wstrade/client"
	"wstrade/config"
	mmcontext "wstrade/context"
	"wstrade/utils/logger"
)

func StartBinanceTradingWebsocket(globalConfig *config.Config, ctxt *mmcontext.GlobalContext) {
	StartWsOrderConnection(globalConfig, ctxt)
}

func StartWsOrderConnection(globalConfig *config.Config, ctxt *mmcontext.GlobalContext) {
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[BinanceTradingWs] Recovered from panic: %v", rc)
			}

			logger.Warn("[BinanceTradingWs] Create Trading Connection and Logon")
		}()
		for {
		ReConnect:
			errChan := make(chan *binanceFutures.Error)
			loginCh := make(chan *binanceFutures.LoginResp)
			orderCh := make(chan *binanceFutures.OrderResp)

			var bnClient = client.BinanceClient{}
			bnClient.Init(globalConfig)

			bnClient.FutresWsClient.SetChannels(errChan, loginCh, orderCh)
			bnClient.FutresWsClient.Connect()
			bnClient.FutresWsClient.Login()

			for {
				select {
				case err := <-errChan:
					logger.Error("[OrderWebSocket] Futures-Order Occur Some Error \t%+v", err)
				case s := <-loginCh:
					logger.Info("[OrderWebSocket] Login Info: %+v", s)
				case b := <-bnClient.FutresWsClient.DoneChan:
					logger.Info("[OrderWebSocket] Futures-Order End\t%v", b)
					// 暂停一秒再跳出，避免异常时频繁发起重连
					logger.Warn("[OrderWebSocket] Will Reconnect Futures-Order-WebSocket After 1 Second")
					time.Sleep(time.Second * 10)
					goto ReConnect
				case r := <-orderCh:
					logger.Info("order op result %+v", r)
				case pOrder := <-ctxt.PlaceOrderChan:
					bnClient.FutresWsClient.PlaceOrder(pOrder)
				case cOrder := <-ctxt.CancelOrderChan:
					bnClient.FutresWsClient.CancelOrder(cOrder)
				}
			}
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
