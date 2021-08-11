package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"log"
	"time"
)

//import (
//	"Tradeasy/config"
//	"encoding/json"
//	"github.com/gorilla/websocket"
//	"log"
//	"time"
//)

//orderUpdated is used by GetUpdates whenever matching algo updates the data for particular order
//var orderUpdated = make(chan OrderResponse)

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "PENDING"
	resp.OrderID = buyOrderBody.OrderID
	resp.StockName = buyOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.BuyOrderBook{
		OrderID:           buyOrderBody.OrderID,
		StockTickerSymbol: buyOrderBody.StockName,
		OrderQuantity:     buyOrderBody.Quantity,
		OrderStatus:       "PENDING",
		OrderPrice:        buyOrderBody.LimitPrice,
		OrderType:         buyOrderBody.OrderType,
		CreatedAt:         buyOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if buyOrderBody.OrderType != "Market" && buyOrderBody.OrderType != "Limit" {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if buyOrderBody.OrderType == "Limit" && buyOrderBody.LimitPrice == 0 {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(buyOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, nil
	}
	if buyOrderBody.OrderType == "Market" {
		newEntry.OrderPrice = ltp
	}
	err = config.DB.Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, err
	}
	go BuyOrderMatchingAlgo(buyOrderBody)
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "PENDING"
	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.SellOrderBook{
		OrderID:           sellOrderBody.OrderID,
		StockTickerSymbol: sellOrderBody.StockName,
		OrderQuantity:     sellOrderBody.Quantity,
		OrderStatus:       "PENDING",
		OrderPrice:        sellOrderBody.LimitPrice,
		OrderType:         sellOrderBody.OrderType,
		CreatedAt:         sellOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if sellOrderBody.OrderType != "Market" && sellOrderBody.OrderType != "Limit" {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if sellOrderBody.OrderType == "Limit" && sellOrderBody.LimitPrice == 0 {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(sellOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, nil
	}
	if sellOrderBody.OrderType == "Market" {
		newEntry.OrderPrice = ltp
	}
	err = config.DB.Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, err
	}
	go SellOrderMatchingAlgo(sellOrderBody)
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", orderId).Error
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
	err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", orderId).Error
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
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC "+" LIMIT 5", stock).Scan(&buyBook).Error
	vdRes.Message = "Error in fetching data"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC"+" LIMIT 5", stock).Scan(&sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}

//func GetOrderUpdates(conn *websocket.Conn) {
//	for {
//		select {
//		case orderMsg := <-orderUpdated:
//			orderJson, err := json.Marshal(orderMsg)
//			if err != nil {
//				log.Println("Error while converting stocks to bytes", err)
//				continue
//			}
//
//			if err := conn.WriteMessage(websocket.TextMessage, orderJson); err != nil {
//				log.Println("Error during writing stocks to websocket:", err)
//				continue
//			}
//		}
//	}
//}
//
//func GetStockUpdates(conn *websocket.Conn, timeInterval time.Duration) {
//	for range time.Tick(timeInterval) {
//		var stocks []StockDetails
//		if err := config.DB.Table("stocks").Find(&stocks).Error; err != nil {
//			log.Println("Error while pulling stocks from stock exchange:", err)
//			continue
//		}
//		stockJson, err := json.Marshal(&stocks)
//		if err != nil {
//			log.Println("Error while converting stocks to bytes", err)
//			continue
//		}
//
//		if err := conn.WriteMessage(websocket.TextMessage, stockJson); err != nil {
//			log.Println("Error during writing stocks to websocket:", err)
//		}
//	}
//}
