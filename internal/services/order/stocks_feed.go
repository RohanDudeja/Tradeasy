package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"time"
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
			CreatedAt: time.Now(),
		}
		if err=config.DB.Table("stocks_feed").Create(newStock).Error;err!=nil{
			return err
		}
	}
	return nil
}
