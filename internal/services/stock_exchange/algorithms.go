package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"fmt"
	"math"
	"time"
)

var orderResponse = make(chan OrderResponse)

func UpdateLTP(ltp int, stock string) {
	currentStock := model.Stocks{}
	config.DB.Raw("SELECT * FROM stocks WHERE stock_name = ?", stock).Scan(&currentStock)
	currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
	currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
	currentStock.LTP = ltp
	currentStock.UpdatedAt = time.Now()
	fmt.Println(currentStock)
	config.DB.Exec("UPDATE stocks SET ltp = ? ,high_price = ? ,low_price = ? , updated_at = ?  WHERE stock_ticker_symbol = ? ", ltp, currentStock.HighPrice, currentStock.LowPrice, currentStock.UpdatedAt, stock)
	//config.DB.Table("stocks").Save(&currentStock)
}

func BuyOrderMatchingAlgo(buyOrderBody OrderRequest, resp OrderResponse, id int) {

	var sellBook []model.SellOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC ", buyOrderBody.StockName).Scan(&sellBook).Error
	if err != nil {
		resp.Status = "Db Fetch failed, order cancelled"
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		orderResponse <- resp
		return
	}
	if len(sellBook) == 0 {
		// abort with message not enough shares
		resp.Status = "Db Fetch failed, order cancelled"
		orderResponse <- resp
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		return
	}
	SharesAvailable := 0
	for _, elem := range sellBook {
		SharesAvailable += elem.OrderQuantity
	}

	if SharesAvailable < buyOrderBody.Quantity {
		resp.Status = "Shares not available, order cancelled"
		resp.Message = "Shares not available"
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		orderResponse <- resp
		return
	}
	totalCost := 0
	ltp := 0
	orderedQuantity := buyOrderBody.Quantity
	if buyOrderBody.OrderType == "Limit" {
		//if int(sellBook[0].OrderPrice) <= price {
		for _, elem := range sellBook {

			if elem.OrderPrice <= buyOrderBody.LimitPrice {
				if elem.OrderQuantity > buyOrderBody.Quantity {
					ltp = elem.OrderPrice
					config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID)
					totalCost += elem.OrderPrice * buyOrderBody.Quantity
					buyOrderBody.Quantity -= buyOrderBody.Quantity
					UpdateLTP(ltp, buyOrderBody.StockName)
					break
				} else {
					ltp = elem.OrderPrice
					UpdateLTP(ltp, buyOrderBody.StockName)
					totalCost += elem.OrderPrice * elem.OrderQuantity
					buyOrderBody.Quantity -= elem.OrderQuantity
					config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
					//config.DB.Delete(&model.SellOrderBook{}, elem.ID)
				}
			} else {
				// add to pending order list in order book
				config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE id = ?", buyOrderBody.Quantity, time.Now(), id)
				time.Sleep(30 * time.Second)
				//Recursive call :
				BuyOrderMatchingAlgo(buyOrderBody, resp, id)
			}

		}

		//update response
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = totalCost / orderedQuantity
		// order complete
		orderResponse <- resp

	} else {
		//Market  order
		for _, elem := range sellBook {
			if elem.OrderQuantity > buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID)
				totalCost += elem.OrderPrice * buyOrderBody.Quantity
				buyOrderBody.Quantity -= buyOrderBody.Quantity
				UpdateLTP(ltp, buyOrderBody.StockName)
				break
			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				totalCost += elem.OrderPrice * elem.OrderQuantity
				buyOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
				//config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}
		// If the book ends but order is unfulfilled: Functionality not available yet
		/*
			if quantity!=0{
				//do something
			}
		*/
		//update response
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = totalCost / orderedQuantity
		// order complete
		orderResponse <- resp
	}
	return
}

func SellOrderMatchingAlgo(sellOrderBody OrderRequest, resp OrderResponse, id int) {
	var buyBook []model.BuyOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC ", sellOrderBody.StockName).Scan(&buyBook).Error
	if err != nil {
		resp.Status = "Db fetch failed, Order cancelled"
		config.DB.Exec("DELETE sell_order_book WHERE id = ?", id)
		orderResponse <- resp
		return
	}
	if len(buyBook) == 0 {
		// abort with message not enough shares
		resp.Status = "Db fetch failed, Order Cancelled"
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		orderResponse <- resp
		config.DB.Exec("DELETE sell_order_book WHERE id = ?", id)
		return
	}
	SharesAvailable := 0
	for _, elem := range buyBook {
		SharesAvailable += elem.OrderQuantity
	}

	if SharesAvailable < sellOrderBody.Quantity {
		resp.Status = "Shares not available, order cancelled"
		resp.Message = "Shares not available"
		config.DB.Exec("DELETE sell_order_book WHERE id = ?", id)
		orderResponse <- resp
		return
	}
	totalCost := 0
	ltp := 0
	orderedQuantity := sellOrderBody.Quantity
	if sellOrderBody.OrderType == "Limit" {
		//if int(sellBook[0].OrderPrice) <= price {
		for _, elem := range buyBook {

			if elem.OrderPrice >= sellOrderBody.LimitPrice {
				if elem.OrderQuantity > sellOrderBody.Quantity {
					ltp = elem.OrderPrice
					config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? ,updated_at = ?  WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID)
					totalCost += elem.OrderPrice * sellOrderBody.Quantity
					sellOrderBody.Quantity -= sellOrderBody.Quantity
					UpdateLTP(ltp, sellOrderBody.StockName)
					break
				} else {
					ltp = elem.OrderPrice
					UpdateLTP(ltp, sellOrderBody.StockName)
					totalCost += elem.OrderPrice * elem.OrderQuantity
					sellOrderBody.Quantity -= elem.OrderQuantity
					config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
					//config.DB.Delete(&model.SellOrderBook{}, elem.ID)
				}
			} else {
				// update the pending order list in order book
				config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ?", sellOrderBody.Quantity, time.Now(), id)
				time.Sleep(30 * time.Second)
				//Recursive call :
				SellOrderMatchingAlgo(sellOrderBody, resp, id)
			}
		}

		//update response
		config.DB.Exec("DELETE buy_order_book WHERE id = ?", id)
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = totalCost / orderedQuantity
		// order complete
		orderResponse <- resp

	} else {
		//Market  order
		for _, elem := range buyBook {
			if elem.OrderQuantity > sellOrderBody.Quantity {
				ltp = elem.OrderPrice
				config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID)
				totalCost += elem.OrderPrice * sellOrderBody.Quantity
				sellOrderBody.Quantity -= sellOrderBody.Quantity
				UpdateLTP(ltp, sellOrderBody.StockName)
				break
			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				totalCost += elem.OrderPrice * elem.OrderQuantity
				sellOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID)
				//config.DB.Delete(&model.SellOrderBook{}, elem.ID)
			}
		}

		//update response
		config.DB.Exec("DELETE sell_order_book WHERE id = ?", id)
		resp.Message = "Order Executed Successfully"
		resp.Status = "Completed"
		resp.OrderExecutionTime = time.Now()
		resp.AveragePrice = totalCost / orderedQuantity
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
	newEntry := model.BuyOrderBook{
		OrderID:           buyOrderBody.OrderID,
		StockTickerSymbol: buyOrderBody.StockName,
		OrderQuantity:     buyOrderBody.Quantity,
		OrderStatus:       "Pending",
		OrderPrice:        buyOrderBody.LimitPrice,
		CreatedAt:         buyOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
		//DeletedAt:         nil,
	}
	config.DB.Create(&newEntry)
	var id int
	err = config.DB.Raw("SELECT id FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Scan(&id).Error
	if err != nil {
		resp.Status = "Db fetch failed, order cancelled"
		return resp, nil
	}
	go BuyOrderMatchingAlgo(buyOrderBody, resp, id)
	fmt.Println(<-orderResponse)
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "Pending"
	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	resp.Message = "Order Received, Pending"
	newEntry := model.SellOrderBook{
		OrderID:           sellOrderBody.OrderID,
		StockTickerSymbol: sellOrderBody.StockName,
		OrderQuantity:     sellOrderBody.Quantity,
		OrderStatus:       "Pending",
		OrderPrice:        sellOrderBody.LimitPrice,
		CreatedAt:         sellOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
		//DeletedAt:         nil,
	}
	config.DB.Create(&newEntry)
	var id int
	err = config.DB.Raw("SELECT id FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Scan(&id).Error
	if err != nil {
		resp.Status = "Db fetch failed, order cancelled"
		return resp, nil
	}
	go SellOrderMatchingAlgo(sellOrderBody, resp, id)
	fmt.Println(<-orderResponse)
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ? AND DeletedAt IS NOT NULL", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ? AND DeletedAt IS NOT NULL", orderId).Error
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
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC LIMIT 5", stock).Scan(&buyBook).Error
	vdRes.Message = "Error in fetching data"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC LIMIT 5", stock).Scan(&sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}
