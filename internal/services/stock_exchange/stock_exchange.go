package stock_exchange

import (
	"Tradeasy/config"
	"log"
)

var OrderUpdated = make(chan OrderResponse)

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
