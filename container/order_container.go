package container

import (
	"fmt"
	"sort"
	"sync"
	"wstrade/config"
	"wstrade/utils/logger"
)

type UmCancelOrder struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
}

type Order struct {
	Exchange      config.Exchange
	InstID        string
	OrderSide     config.OrderSide // 订单方向: buy sell
	OrderID       string
	OrderPrice    float64
	OrderVolume   float64 // Okx合约是张数
	CreateAt      int64
	ClientOrderID string             // 自定义订单ID
	Precision     [2]int             // //  [4, 2], 以BTCBUSD为例，BTC的精度是4，BUSD的精度是2
	Status        config.OrderStatus // 订单状态
}

func (order *Order) FormatString() string {
	str := fmt.Sprintf("OrderID=%s, ClientOrderID=%s, OrderSide=%s, Price=%v, Volume=%v, InstID=%s",
		order.OrderID, order.ClientOrderID, order.OrderSide, order.OrderPrice,
		order.OrderVolume, order.InstID)
	return str
}

type OrderList []*Order

func (lst OrderList) Len() int           { return len(lst) }
func (lst OrderList) Less(i, j int) bool { return lst[i].OrderPrice < lst[j].OrderPrice }
func (lst OrderList) Swap(i, j int)      { lst[i], lst[j] = lst[j], lst[i] }

type OrderBook struct {
	Data   OrderList
	RwLock sync.RWMutex
}

// DeleteByClientOrderID 通过ClientOrderID删除对应的Order
func (orderBook *OrderBook) DeleteByClientOrderID(clientOrderID string) *Order {
	orderBook.RwLock.Lock()
	defer orderBook.RwLock.Unlock()
	var order *Order = nil
	for i := 0; i < len(orderBook.Data); {
		if orderBook.Data[i].ClientOrderID == clientOrderID {
			order = orderBook.Data[i]
			orderBook.Data = append(orderBook.Data[:i], orderBook.Data[i+1:]...)
			break
		} else {
			i++
		}
	}
	return order
}

// UpdateStatus 更新order状态
func (orderBook *OrderBook) UpdateStatus(clientOrderID string, status config.OrderStatus) *Order {
	orderBook.RwLock.Lock()
	defer orderBook.RwLock.Unlock()
	for i := 0; i < len(orderBook.Data); i++ {
		if orderBook.Data[i].ClientOrderID == clientOrderID {
			orderBook.Data[i].Status = status
			return orderBook.Data[i]
		}
	}
	return nil
}

func (orderBook *OrderBook) Size() int {
	orderBook.RwLock.RLock()
	defer orderBook.RwLock.RUnlock()
	return len(orderBook.Data)
}

func (orderBook *OrderBook) SizeNoLock() int {
	return len(orderBook.Data)
}

// Add 添加一个Order到OrderList中
func (orderBook *OrderBook) Add(order *Order) {
	orderBook.RwLock.Lock()
	defer orderBook.RwLock.Unlock()

	orderBook.Data = append(orderBook.Data, order)
}

// Update 修改orderbook中order 的信息
func (orderBook *OrderBook) Update(clientOrderID string, orderID string, orderStatus config.OrderStatus, createAt int64) {
	orderBook.RwLock.Lock()
	defer orderBook.RwLock.Unlock()
	for i := 0; i < len(orderBook.Data); i++ {
		if orderBook.Data[i].ClientOrderID == clientOrderID {
			orderBook.Data[i].OrderID = orderID
			orderBook.Data[i].Status = orderStatus
			orderBook.Data[i].CreateAt = createAt
		}
	}
}

// Sort 将订单按照价格从小到大排序
func (orderBook *OrderBook) Sort() {
	orderBook.RwLock.Lock()
	defer orderBook.RwLock.Unlock()
	sort.Sort(orderBook.Data)
}

func (orderBook *OrderBook) SortNoLock() {
	sort.Sort(orderBook.Data)
}

type OrderBookComposite struct {
	BuyOrderBook  map[string]*OrderBook
	SellOrderBook map[string]*OrderBook
	rwLock        *sync.RWMutex
}

func (o *OrderBookComposite) Init(instIDs []string) {
	o.BuyOrderBook = map[string]*OrderBook{}
	o.SellOrderBook = map[string]*OrderBook{}
	for _, instID := range instIDs {
		o.BuyOrderBook[instID] = &OrderBook{}
		o.SellOrderBook[instID] = &OrderBook{}
	}
	o.rwLock = new(sync.RWMutex)
}

func (o *OrderBookComposite) GetOrderBook(instID string, orderSide config.OrderSide) *OrderBook {
	o.rwLock.RLock()
	defer o.rwLock.RUnlock()

	if orderSide == config.BuyOrderSide {
		wrapper, has := o.BuyOrderBook[instID]
		logger.Debug("===%s %s wrapper=%+v", instID, orderSide, wrapper)
		if has {
			return wrapper
		} else {
			return nil
		}
	} else {
		wrapper, has := o.SellOrderBook[instID]
		logger.Debug("===%s %s wrapper=%+v", instID, orderSide, wrapper)
		if has {
			return wrapper
		} else {
			return nil
		}
	}
}
