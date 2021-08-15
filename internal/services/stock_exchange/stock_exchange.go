package stock_exchange

import (
	model "Tradeasy/internal/model/stock_exchange"
	"Tradeasy/internal/provider/database"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//orderUpdated is used by GetUpdates whenever matching algo updates the data for particular order
var orderUpdated = make(chan OrderResponse)

const (
	Pending   = "PENDING"
	Completed = "COMPLETED"
	Partial   = "PARTIAL"
	Cancelled = "CANCELLED"
	Failed    = "FAILED"
	Market    = "Market"
	Limit     = "Limit"
)

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.OrderID = buyOrderBody.OrderID
	resp.StockName = buyOrderBody.StockName
	currentTime := time.Now().Hour()
	if currentTime >= 3 && currentTime < 9 {
		resp.Status = Cancelled
		resp.Message = "Cannot place orders at this time. Market closed."
		return resp, nil
	}
	resp.Status = Pending
	resp.Message = "Order Received"
	newEntry := model.BuyOrderBook{
		OrderID:           buyOrderBody.OrderID,
		StockTickerSymbol: buyOrderBody.StockName,
		OrderQuantity:     buyOrderBody.Quantity,
		OrderStatus:       Pending,
		OrderPrice:        buyOrderBody.LimitPrice,
		OrderType:         buyOrderBody.OrderType,
		CreatedAt:         buyOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if buyOrderBody.OrderType != Market && buyOrderBody.OrderType != Limit {
		resp.Status = Cancelled
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if buyOrderBody.OrderType == Limit && buyOrderBody.LimitPrice == 0 {
		resp.Status = Cancelled
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(buyOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = Failed
		resp.Message = "Internal Error"
		return resp, nil
	}
	if buyOrderBody.OrderType == Market {
		newEntry.OrderPrice = ltp
	}
	err = database.GetDB().Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = Failed
		resp.Message = "Internal Error"
		return resp, err
	}
	go BuyOrderMatchingAlgo(buyOrderBody)
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	currentTime := time.Now().Hour()
	if currentTime >= 3 && currentTime < 9 {
		resp.Status = Cancelled
		resp.Message = "Cannot place orders at this time. Market closed."
		return resp, nil
	}
	resp.Status = Pending
	resp.Message = "Order Received"
	newEntry := model.SellOrderBook{
		OrderID:           sellOrderBody.OrderID,
		StockTickerSymbol: sellOrderBody.StockName,
		OrderQuantity:     sellOrderBody.Quantity,
		OrderStatus:       Pending,
		OrderPrice:        sellOrderBody.LimitPrice,
		OrderType:         sellOrderBody.OrderType,
		CreatedAt:         sellOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if sellOrderBody.OrderType != Market && sellOrderBody.OrderType != Limit {
		resp.Status = Cancelled
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if sellOrderBody.OrderType == Limit && sellOrderBody.LimitPrice == 0 {
		resp.Status = Cancelled
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(sellOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = Failed
		resp.Message = "Internal Error"
		return resp, nil
	}
	if sellOrderBody.OrderType == Market {
		newEntry.OrderPrice = ltp
	}
	err = database.GetDB().Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = Failed
		resp.Message = "Internal Error"
		return resp, err
	}
	go SellOrderMatchingAlgo(sellOrderBody)
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = database.GetDB().Table("buy_order_book").Where("order_id= ?", orderId).Delete(model.BuyOrderBook{}).Error
	if err != nil {
		deleteRes.Message = "Failed"
		deleteRes.Success = false
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	deleteRes.Success = true
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = database.GetDB().Table("sell_order_book").Where("order_id= ?", orderId).Delete(model.SellOrderBook{}).Error
	if err != nil {
		deleteRes.Success = false
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Success = true
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// ViewMarketDepth ...Returns 5 depth orders from order book
func ViewMarketDepth(stock string) (vdRes ViewDepthResponse, err error) {

	var buyBook []model.BuyOrderBook
	err = database.GetDB().Raw("SELECT * FROM buy_order_book WHERE deleted_at IS NULL AND stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC "+" LIMIT 5", stock).Scan(&buyBook).Error
	vdRes.Message = "Internal Error"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = database.GetDB().Raw("SELECT * FROM sell_order_book WHERE deleted_at IS NULL AND stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC"+" LIMIT 5", stock).Scan(&sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}

func GetOrderUpdates(conn *websocket.Conn) {
	for {
		select {
		case orderMsg := <-orderUpdated:
			orderJson, err := json.Marshal(orderMsg)
			if err != nil {
				log.Println("Error while converting stocks to bytes", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, orderJson); err != nil {
				log.Println("Error during writing stocks to websocket:", err)
				continue
			}
		}
	}
}

func GetStockUpdates(conn *websocket.Conn, timeInterval time.Duration) {
	for range time.Tick(timeInterval) {
		var stocks []StockDetails
		if err := database.GetDB().Table("stocks").Find(&stocks).Error; err != nil {
			log.Println("Error while pulling stocks from stock exchange:", err)
			continue
		}
		stockJson, err := json.Marshal(&stocks)
		if err != nil {
			log.Println("Error while converting stocks to bytes", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, stockJson); err != nil {
			log.Println("Error during writing stocks to websocket:", err)
		}
	}
}
