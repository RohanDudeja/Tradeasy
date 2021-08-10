package order

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/services/stock_exchange"
	"log"
)

func UpdateStocksFeed(res []stock_exchange.StockDetails) (err error) {
	for _, stock := range res {
		newStock := model.StocksFeed{
			StockName: stock.StockName,
			LTP:       stock.LTP,
			Open:      stock.Open,
			High:      stock.High,
			Low:       stock.Low,
			TradedAt:  stock.UpdatedAt,
		}
		if err = database.GetDB().Table("stocks_feed").Create(&newStock).Error; err != nil {
			log.Printf("Error in Creating %s stock field\n", newStock.StockName)
			continue
		}
	}
	return nil
}
