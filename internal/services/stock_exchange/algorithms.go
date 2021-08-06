package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"math"
	"time"
)

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	orderId := buyOrderBody.OrderID
	quantity := buyOrderBody.Quantity
	stock := buyOrderBody.StockName
	price := buyOrderBody.LimitPrice
	orderType := buyOrderBody.OrderType

	resp.Status = "Cancelled by engine"
	resp.OrderID = orderId
	resp.StockName = stock
	resp.Message = "Couldn't execute due to unavailability of shares"

	var sellBook []model.SellOrderBook
	// db lock
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", stock).Scan(sellBook).Error
	if err != nil {
		return resp, err
	}
	if len(sellBook) == 0 {
		// abort with message not enough shares
		return resp, nil
	}
	SharesAvailable := 0
	for _, elem := range sellBook {
		SharesAvailable += elem.OrderQuantity
	}

	if uint(SharesAvailable) < quantity {
		// abort with message not enough shares
		return resp, err
	}
	totalCost := 0
	ltp := 0
	if orderType == "Limit" {
		if uint(sellBook[0].OrderPrice) >= price {

			for _, elem := range sellBook {
				if uint(elem.OrderQuantity) > quantity {
					ltp = elem.OrderPrice
					config.DB.Raw("UPDATE sell_order_book SET order_quantity = ? WHERE id = ? ", uint(elem.OrderQuantity)-quantity, elem.ID)
					break
				} else {
					ltp = elem.OrderPrice
					totalCost += elem.OrderPrice * elem.OrderQuantity
					quantity -= uint(elem.OrderQuantity)
					config.DB.Delete(&model.SellOrderBook{}, elem.ID)
				}
			}

			//set values
			currentStock := model.Stocks{}
			config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
			currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
			currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
			config.DB.Save(&currentStock)
			//update response
			resp.Message = "Order Executed Successfully"
			resp.Status = "Completed"
			resp.OrderExecutionTime = time.Now()
			resp.AveragePrice = uint(totalCost / int(quantity))
			// order complete
			return resp, err
		} else {
			// add to pending order list in order book
			newEntry := model.BuyOrderBook{
				OrderID:           orderId,
				StockTickerSymbol: stock,
				OrderQuantity:     int(quantity),
				OrderStatus:       "Pending",
				OrderPrice:        int(price),
				//CreatedAt:         orderTime,
				//UpdatedAt:         time.Now(),
				//DeletedAt:         nil,
			}
			config.DB.Create(&newEntry)
			//unlock db
			time.Sleep(30 * time.Second)
			//Recursive call :
			BuyOrder(buyOrderBody)
			//update status and return
		}

	} else {
		//execute order
		for _, elem := range sellBook {
			if uint(elem.OrderQuantity) > quantity {
				ltp = elem.OrderPrice
				config.DB.Raw("UPDATE sell_order_book SET order_quantity = ? WHERE id = ", uint(elem.OrderQuantity)-quantity, elem.ID)
				break
			} else {
				ltp = elem.OrderPrice
				totalCost += elem.OrderPrice * elem.OrderQuantity
				quantity -= uint(elem.OrderQuantity)
				config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}

		//set values
		currentStock := model.Stocks{}
		config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
		currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
		currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
		config.DB.Save(&currentStock)
		//update response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(quantity))
		// order complete
		return resp, nil
	}
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	orderId := sellOrderBody.OrderID
	quantity := sellOrderBody.Quantity
	stock := sellOrderBody.StockName
	price := sellOrderBody.LimitPrice
	orderType := sellOrderBody.OrderType

	resp.Status = "Cancelled by engine"
	resp.OrderID = orderId
	resp.StockName = stock
	resp.Message = "Couldn't execute due to unavailability of shares"

	var buyBook []model.BuyOrderBook
	// db lock
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC ", stock).Scan(buyBook).Error
	if err != nil {
		return resp, err
	}
	if len(buyBook) == 0 {
		// abort with message not enough shares
		return resp, nil
	}
	SharesAvailable := 0
	for _, elem := range buyBook {
		SharesAvailable += elem.OrderQuantity
	}

	if uint(SharesAvailable) < quantity {
		// abort with message not enough shares
		return resp, nil
	}
	totalCost := 0
	ltp := 0
	if orderType == "Limit" {
		if uint(buyBook[0].OrderPrice) >= price {

			for _, elem := range buyBook {
				if uint(elem.OrderQuantity) > quantity {
					ltp = elem.OrderPrice
					config.DB.Raw("UPDATE buy_order_book SET order_quantity = ? WHERE id = ? ", uint(elem.OrderQuantity)-quantity, elem.ID)
					break
				} else {
					ltp = elem.OrderPrice
					totalCost += elem.OrderPrice * elem.OrderQuantity
					quantity -= uint(elem.OrderQuantity)
					config.DB.Delete(&model.BuyOrderBook{}, elem.ID)
				}
			}

			//set values
			currentStock := model.Stocks{}
			config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
			currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
			currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
			config.DB.Save(&currentStock)
			//config.DB.Raw("UPDATE stocks SET ltp = ? WHERE stock_ticker_symbol = ?", ltp, stock)
			//update response
			resp.Message = "Order Executed Successfully"
			resp.Status = "Completed"
			resp.OrderExecutionTime = time.Now()
			resp.AveragePrice = uint(totalCost / int(quantity))
			//order completed
			return resp, nil

		} else {
			// add to pending order list in order book
			newEntry := model.SellOrderBook{
				OrderID:           orderId,
				StockTickerSymbol: stock,
				OrderQuantity:     int(quantity),
				OrderStatus:       "Pending",
				OrderPrice:        int(price),
				//CreatedAt:         orderTime,
				//UpdatedAt:         time.Now(),
				//DeletedAt:         nil,
			}
			config.DB.Create(&newEntry)
			//unlock db
			time.Sleep(30 * time.Second)
			//Recursive call :
			SellOrder(sellOrderBody)
			//update status and return
		}

	} else {
		//execute order
		for _, elem := range buyBook {
			if uint(elem.OrderQuantity) > quantity {
				ltp = elem.OrderPrice
				config.DB.Raw("UPDATE sell_order_book SET order_quantity = ? WHERE id = ", uint(elem.OrderQuantity)-quantity, elem.ID)
				break
			} else {
				ltp = elem.OrderPrice
				totalCost += elem.OrderPrice * elem.OrderQuantity
				quantity -= uint(elem.OrderQuantity)
				config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}
		//update values
		currentStock := model.Stocks{}
		config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
		currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
		currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
		config.DB.Save(&currentStock)
		//config.DB.Raw("UPDATE stocks SET ltp = ? WHERE stock_ticker_symbol = ?", ltp, stock)
		//set response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(quantity))
		// order complete
		return resp, nil
	}
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Raw("DELETE FROM buy_order_book WHERE order_id = ?", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Raw("DELETE FROM buy_order_book WHERE order_id = ?", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// ViewMarketDepth ...Returns 5 depth orders from order book
func ViewMarketDepth(stock string) (vdRes ViewDepthResponse, err error) {

	var buyBook []model.BuyOrderBook
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", stock).Scan(buyBook).Error
	vdRes.Message = "Error in fetching data"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", stock).Scan(sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}
