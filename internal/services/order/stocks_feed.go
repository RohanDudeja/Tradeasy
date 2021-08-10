package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"log"
)

func UpdateStocksFeed(res []stock_exchange.StockDetails) (err error){
	for _,stock:=range res{
		newStock:=model.StocksFeed{
			StockName: stock.StockName,
			LTP:       stock.LTP,
			Open:      stock.Open,
			High:      stock.High,
			Low:       stock.Low,
			TradedAt:  stock.UpdatedAt,
		}
		if err=config.DB.Table("stocks_feed").Create(&newStock).Error;err!=nil{
			log.Println("Error in Creating a stock field")
			continue
		}
	}
	return nil
}
