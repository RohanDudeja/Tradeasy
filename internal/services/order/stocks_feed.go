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
			LTP:       int(stock.LTP),
			Open:      int(stock.Open),
			High:      int(stock.High),
			Low:       int(stock.Low),
			TradedAt:  stock.UpdatedAt,
			CreatedAt: time.Now(),
		}
		if err=config.DB.Create(newStock).Error;err!=nil{
			return err
		}
	}
	return nil
}
