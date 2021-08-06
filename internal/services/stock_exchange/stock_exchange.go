package stock_exchange

import (
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var OrderUpdated chan bool

func GetAllStocks() (stocks *[]StockDetails, err error) {
	if err := DB.Find(stocks).Error; err != nil {
		return stocks, err
	}
	return stocks, nil
}
func UpdateOrder() (order OrderResponse, err error) {
	return order, nil
}
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
