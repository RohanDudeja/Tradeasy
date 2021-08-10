package stock_exchange

import (
	"Tradeasy/config"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

//orderUpdated is used by GetUpdates whenever matching algo updates the data for particular order
var orderUpdated = make(chan OrderResponse)

func BuyOrder(buyReq OrderRequest) (buyRes OrderResponse, err error) {
	return buyRes, err
}

func SellOrder(sellReq OrderRequest) (sellRes OrderResponse, err error) {
	return sellRes, err
}

func DeleteBuyOrder(OrderId string) (delRes DeleteResponse, err error) {
	return delRes, err
}

func DeleteSellOrder(OrderId string) (delRes DeleteResponse, err error) {
	return delRes, err
}

func ViewMarketDepth(stockName string) (vdRes ViewDepthResponse, err error) {
	return vdRes, err
}
func StockWrite() (stocks []StockDetails, err error) {
	if err = config.DB.Table("stocks").Find(&stocks).Error; err != nil {
		log.Println("Error while pulling stocks from stock exchange:", err)
		return stocks, err
	}
	return stocks, nil
}
func GetUpdates(conn *websocket.Conn) {
	for {
		select {
		case orderMsg := <-orderUpdated:
			orderJson, err := json.Marshal(orderMsg)
			if err != nil {
				log.Println("Error while converting stocks to bytes", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, orderJson); err != nil {
				log.Println("Error during writing stocks to websocket:", err)
				return
			}

		}
	}
}
