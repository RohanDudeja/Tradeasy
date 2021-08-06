package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"math"
	"time"
)

func UpdateLTP(ltp int, stock string) {
	currentStock := model.Stocks{}
	config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
	currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
	currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
	currentStock.LTP = ltp
	config.DB.Save(&currentStock)
}

func BuyOrderMatchingAlgo(buyOrderBody OrderRequest, orderResponse chan OrderResponse, resp OrderResponse) {

	var sellBook []model.SellOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", buyOrderBody.StockName).Scan(&sellBook).Error
	if err != nil {
		resp.Status = "Db Fetch failed"
		orderResponse <- resp
	}
	if len(sellBook) == 0 {
		// abort with message not enough shares
		resp.Status = "Db Fetch failed"
		orderResponse <- resp
	}
	SharesAvailable := 0
	for _, elem := range sellBook {
		SharesAvailable += elem.OrderQuantity
	}

	totalCost := 0
	ltp := 0
	orderedQuantity := buyOrderBody.Quantity
	if buyOrderBody.OrderType == "Limit" {
		//if uint(sellBook[0].OrderPrice) <= price {
		for _, elem := range sellBook {

			if uint(elem.OrderPrice) <= buyOrderBody.LimitPrice {
				if uint(elem.OrderQuantity) > buyOrderBody.Quantity {
					ltp = elem.OrderPrice
					config.DB.Raw("UPDATE sell_order_book SET order_quantity = ? WHERE id = ? ", uint(elem.OrderQuantity)-buyOrderBody.Quantity, elem.ID)
					UpdateLTP(ltp, buyOrderBody.StockName)
					break
				} else {
					ltp = elem.OrderPrice
					UpdateLTP(ltp, buyOrderBody.StockName)
					totalCost += elem.OrderPrice * elem.OrderQuantity
					buyOrderBody.Quantity -= uint(elem.OrderQuantity)
					config.DB.Delete(&model.SellOrderBook{}, elem.ID)
				}
			} else {
				// add to pending order list in order book
				newEntry := model.BuyOrderBook{
					OrderID:           buyOrderBody.OrderID,
					StockTickerSymbol: buyOrderBody.StockName,
					OrderQuantity:     int(buyOrderBody.Quantity),
					OrderStatus:       "Pending",
					OrderPrice:        int(buyOrderBody.LimitPrice),
					//CreatedAt:         orderTime,
					//UpdatedAt:         time.Now(),
					//DeletedAt:         nil,
				}
				config.DB.Create(&newEntry)
				time.Sleep(30 * time.Second)
				//Recursive call :
				BuyOrderMatchingAlgo(buyOrderBody, orderResponse, resp)
			}

		}
		//update response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(orderedQuantity))
		// order complete
		orderResponse <- resp

	} else {
		//Market  order
		for _, elem := range sellBook {
			if uint(elem.OrderQuantity) > buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				config.DB.Raw("UPDATE sell_order_book SET order_quantity = ? WHERE id = ", uint(elem.OrderQuantity)-buyOrderBody.Quantity, elem.ID)
				break
			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				totalCost += elem.OrderPrice * elem.OrderQuantity
				buyOrderBody.Quantity -= uint(elem.OrderQuantity)
				config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}

		if buyOrderBody.Quantity != 0 {
			// add to pending order list in order book
			newEntry := model.BuyOrderBook{
				OrderID:           buyOrderBody.OrderID,
				StockTickerSymbol: buyOrderBody.StockName,
				OrderQuantity:     int(buyOrderBody.Quantity),
				OrderStatus:       "Pending",
				//OrderPrice:        int(buyOrderBody.LimitPrice),
				//CreatedAt:         orderTime,
				//UpdatedAt:         time.Now(),
				//DeletedAt:         nil,
			}
			config.DB.Create(&newEntry)
			time.Sleep(30 * time.Second)
			//Recursive call :
			BuyOrderMatchingAlgo(buyOrderBody, orderResponse, resp)
		}
		//update response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(orderedQuantity))
		// order complete
		orderResponse <- resp
	}

}

func SellOrderMatchingAlgo(sellOrderBody OrderRequest, orderResponse chan OrderResponse, resp OrderResponse) {
	var buyBook []model.BuyOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", sellOrderBody.StockName).Scan(&buyBook).Error
	if err != nil {
		resp.Status = "Db fetch failed"
		orderResponse <- resp
	}
	if len(buyBook) == 0 {
		// abort with message not enough shares
		resp.Status = "Db fetch failed"
		orderResponse <- resp
	}
	SharesAvailable := 0
	for _, elem := range buyBook {
		SharesAvailable += elem.OrderQuantity
	}

	totalCost := 0
	ltp := 0
	orderedQuantity := sellOrderBody.Quantity
	if sellOrderBody.OrderType == "Limit" {
		//if uint(sellBook[0].OrderPrice) <= price {
		for _, elem := range buyBook {

			if uint(elem.OrderPrice) >= sellOrderBody.LimitPrice {
				if uint(elem.OrderQuantity) > sellOrderBody.Quantity {
					ltp = elem.OrderPrice
					config.DB.Raw("UPDATE buy_order_book SET order_quantity = ? WHERE id = ? ", uint(elem.OrderQuantity)-sellOrderBody.Quantity, elem.ID)
					UpdateLTP(ltp, sellOrderBody.StockName)
					break
				} else {
					ltp = elem.OrderPrice
					UpdateLTP(ltp, sellOrderBody.StockName)
					totalCost += elem.OrderPrice * elem.OrderQuantity
					sellOrderBody.Quantity -= uint(elem.OrderQuantity)
					config.DB.Delete(&model.SellOrderBook{}, elem.ID)
				}
			} else {
				// add to pending order list in order book
				newEntry := model.BuyOrderBook{
					OrderID:           sellOrderBody.OrderID,
					StockTickerSymbol: sellOrderBody.StockName,
					OrderQuantity:     int(sellOrderBody.Quantity),
					OrderStatus:       "Pending",
					OrderPrice:        int(sellOrderBody.LimitPrice),
					//CreatedAt:         orderTime,
					//UpdatedAt:         time.Now(),
					//DeletedAt:         nil,
				}
				config.DB.Create(&newEntry)
				time.Sleep(30 * time.Second)
				//Recursive call :
				SellOrderMatchingAlgo(sellOrderBody, orderResponse, resp)
			}

		}
		//update response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(orderedQuantity))
		// order complete
		orderResponse <- resp

	} else {
		//Market  order
		for _, elem := range buyBook {
			if uint(elem.OrderQuantity) > sellOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				config.DB.Raw("UPDATE buy_order_book SET order_quantity = ? WHERE id = ", uint(elem.OrderQuantity)-sellOrderBody.Quantity, elem.ID)
				break
			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				totalCost += elem.OrderPrice * elem.OrderQuantity
				sellOrderBody.Quantity -= uint(elem.OrderQuantity)
				config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}

		if sellOrderBody.Quantity != 0 {
			// add to pending order list in order book
			newEntry := model.BuyOrderBook{
				OrderID:           sellOrderBody.OrderID,
				StockTickerSymbol: sellOrderBody.StockName,
				OrderQuantity:     int(sellOrderBody.Quantity),
				OrderStatus:       "Pending",
				//OrderPrice:        int(buyOrderBody.LimitPrice),
				//CreatedAt:         orderTime,
				//UpdatedAt:         time.Now(),
				//DeletedAt:         nil,
			}
			config.DB.Create(&newEntry)
			time.Sleep(30 * time.Second)
			//Recursive call :
			SellOrderMatchingAlgo(sellOrderBody, orderResponse, resp)
		}
		//update response
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = uint(totalCost / int(orderedQuantity))
		// order complete
		orderResponse <- resp
	}
}

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "Pending"
	resp.OrderID = buyOrderBody.OrderID
	resp.StockName = buyOrderBody.StockName
	resp.Message = "Order Received, Pending"
	//orderUpdate := make(chan string)
	orderResponse := make(chan OrderResponse)
	go BuyOrderMatchingAlgo(buyOrderBody, orderResponse, resp)

	response := <-orderResponse
	return response, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "Cancelled by engine"
	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	resp.Message = "Couldn't execute due to unavailability of shares"
	orderResponse := make(chan OrderResponse)
	go SellOrderMatchingAlgo(sellOrderBody, orderResponse, resp)

	response := <-orderResponse
	return response, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Raw("DELETE FROM buy_order_book WHERE order_id = ? AND DeletedAt IS NOT NULL", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Raw("DELETE FROM sell_order_book WHERE order_id = ? AND DeletedAt IS NOT NULL", orderId).Error
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
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC LIMIT 5", stock).Scan(buyBook).Error
	vdRes.Message = "Error in fetching data"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC LIMIT 5", stock).Scan(sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}
