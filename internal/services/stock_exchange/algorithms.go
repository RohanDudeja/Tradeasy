package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"math"
	"time"
)

var orderResponse = make(chan OrderResponse)

func UpdateLTP(ltp int, stock string) {
	currentStock := model.Stocks{}
	config.DB.Table("stocks").Where("stock_name = ?", stock).Find(&currentStock)
	currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
	currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
	currentStock.LTP = ltp
	currentStock.UpdatedAt = time.Now()
	config.DB.Table("stocks").Where("stock_ticker_symbol = ?", stock).Updates(&currentStock)

}
func SendResponse(resp OrderResponse, status string, message string, orderId string, price int, quantity int) {
	resp.Status = status
	resp.Message = message
	resp.OrderID = orderId
	resp.AveragePrice = price
	resp.OrderExecutionTime = time.Now()
	resp.Quantity = quantity
	orderResponse <- resp
}
func BuyLimitOrder(buyOrderBody OrderRequest, sellBook []model.SellOrderBook, resp OrderResponse) {

	ltp := 0
	for _, elem := range sellBook {
		if elem.OrderPrice <= buyOrderBody.LimitPrice {
			if elem.OrderQuantity > buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID)
				config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
				SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(resp, "Partial", "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				buyOrderBody.Quantity -= buyOrderBody.Quantity
				break
			} else if elem.OrderQuantity < buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				SendResponse(resp, "Partial", "Order Executed Partially", resp.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
				buyOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
				config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", buyOrderBody.Quantity, time.Now(), buyOrderBody.OrderID)
			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, buyOrderBody.StockName)
				SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				buyOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
				config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
			}
		} else {
			// add to pending order list in order book
			time.Sleep(30 * time.Second)
			//Recursive call :
			BuyOrderMatchingAlgo(buyOrderBody, resp)
		}

	}
}

func BuyMarketOrder(buyOrderBody OrderRequest, sellBook []model.SellOrderBook, resp OrderResponse) {

	ltp := 0
	for _, elem := range sellBook {
		if elem.OrderQuantity > buyOrderBody.Quantity {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, buyOrderBody.StockName)
			config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID)
			config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
			SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(resp, "Partial", "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			buyOrderBody.Quantity -= buyOrderBody.Quantity
			break
		} else if elem.OrderQuantity < buyOrderBody.Quantity {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, buyOrderBody.StockName)
			SendResponse(resp, "Partial", "Order Executed Partially", resp.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
			buyOrderBody.Quantity -= elem.OrderQuantity
			config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
			config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", buyOrderBody.Quantity, time.Now(), buyOrderBody.OrderID)
		} else {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, buyOrderBody.StockName)
			SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			buyOrderBody.Quantity -= elem.OrderQuantity
			config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID)
			config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
		}
	}
}
func SellLimitOrder(sellOrderBody OrderRequest, buyBook []model.BuyOrderBook, resp OrderResponse) {

	ltp := 0
	for _, elem := range buyBook {

		if elem.OrderPrice >= sellOrderBody.LimitPrice {
			if elem.OrderQuantity > sellOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? ,updated_at = ?  WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID)
				config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID)
				SendResponse(resp, "Completed", "Order Executed Partially", resp.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(resp, "Partial", "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				sellOrderBody.Quantity -= sellOrderBody.Quantity
				break
			} else if elem.OrderQuantity < sellOrderBody.Quantity {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				SendResponse(resp, "Completed", "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(resp, "Partial", "Order Executed Partially", resp.OrderID, elem.OrderPrice, elem.OrderQuantity)
				sellOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID)
				config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", sellOrderBody.Quantity, time.Now(), sellOrderBody.OrderID)

			} else {
				ltp = elem.OrderPrice
				UpdateLTP(ltp, sellOrderBody.StockName)
				SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				sellOrderBody.Quantity -= elem.OrderQuantity
				config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID)
				config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID)
			}
		} else {
			time.Sleep(30 * time.Second)
			//Recursive call :
			SellOrderMatchingAlgo(sellOrderBody, resp)
		}
	}

}
func SellMarketOrder(sellOrderBody OrderRequest, buyBook []model.BuyOrderBook, resp OrderResponse) {

	ltp := 0
	for _, elem := range buyBook {
		if elem.OrderQuantity > sellOrderBody.Quantity {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, sellOrderBody.StockName)
			config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID)
			config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", resp.OrderID)
			SendResponse(resp, "Completed", "Order Executed Partially", resp.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(resp, "Partial", "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			sellOrderBody.Quantity -= sellOrderBody.Quantity
			break
		} else if elem.OrderQuantity < sellOrderBody.Quantity {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, sellOrderBody.StockName)
			SendResponse(resp, "Completed", "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(resp, "Partial", "Order Executed Partially", resp.OrderID, elem.OrderPrice, elem.OrderQuantity)
			sellOrderBody.Quantity -= elem.OrderQuantity
			config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID)
			config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", sellOrderBody.Quantity, time.Now(), sellOrderBody.OrderID)

		} else {
			ltp = elem.OrderPrice
			UpdateLTP(ltp, sellOrderBody.StockName)
			SendResponse(resp, "Completed", "Order Executed", resp.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(resp, "Completed", "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			sellOrderBody.Quantity -= elem.OrderQuantity
			config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID)
			config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID)
		}
	}
}

func BuyOrderMatchingAlgo(buyOrderBody OrderRequest, resp OrderResponse) {

	var sellBook []model.SellOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC", buyOrderBody.StockName).Scan(&sellBook).Error

	if err != nil || len(sellBook) == 0 {
		// abort
		resp.Status = "Failed"
		orderResponse <- resp
		config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
		return
	}
	SharesAvailable := 0
	for _, elem := range sellBook {
		SharesAvailable += elem.OrderQuantity
	}

	if SharesAvailable < buyOrderBody.Quantity {
		resp.Status = "Cancelled"
		resp.Message = "Shares not available"
		config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID)
		orderResponse <- resp
		return
	}

	if buyOrderBody.OrderType == "Limit" {
		BuyLimitOrder(buyOrderBody, sellBook, resp)
	} else {
		//Market  order
		BuyMarketOrder(buyOrderBody, sellBook, resp)
	}
	//order complete
	orderResponse <- resp
	return
}

func SellOrderMatchingAlgo(sellOrderBody OrderRequest, resp OrderResponse) {
	var buyBook []model.BuyOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC ", sellOrderBody.StockName).Scan(&buyBook).Error

	if err != nil || len(buyBook) == 0 {
		// abort
		resp.Status = "Failed"
		config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID)
		orderResponse <- resp
		return
	}
	SharesAvailable := 0
	for _, elem := range buyBook {
		SharesAvailable += elem.OrderQuantity
	}

	if SharesAvailable < sellOrderBody.Quantity {
		resp.Status = "Cancelled"
		resp.Message = "Shares not available"
		config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID)
		orderResponse <- resp
		return
	}

	if sellOrderBody.OrderType == "Limit" {
		SellLimitOrder(sellOrderBody, buyBook, resp)
	} else {
		//Market  order
		SellMarketOrder(sellOrderBody, buyBook, resp)
	}
}

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "Pending"
	resp.OrderID = buyOrderBody.OrderID
	resp.StockName = buyOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.BuyOrderBook{
		OrderID:           buyOrderBody.OrderID,
		StockTickerSymbol: buyOrderBody.StockName,
		OrderQuantity:     buyOrderBody.Quantity,
		OrderStatus:       "Pending",
		OrderPrice:        buyOrderBody.LimitPrice,
		CreatedAt:         buyOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	config.DB.Create(&newEntry)
	go BuyOrderMatchingAlgo(buyOrderBody, resp)
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "Pending"
	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.SellOrderBook{
		OrderID:           sellOrderBody.OrderID,
		StockTickerSymbol: sellOrderBody.StockName,
		OrderQuantity:     sellOrderBody.Quantity,
		OrderStatus:       "Pending",
		OrderPrice:        sellOrderBody.LimitPrice,
		CreatedAt:         sellOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	config.DB.Create(&newEntry)
	go SellOrderMatchingAlgo(sellOrderBody, resp)
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", orderId).Error
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
	err = config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC "+" LIMIT 5", stock).Scan(&buyBook).Error
	vdRes.Message = "Error in fetching data"
	if err != nil {
		return vdRes, err
	}
	var sellBook []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC"+" LIMIT 5", stock).Scan(&sellBook).Error
	if err != nil {
		return vdRes, err
	}

	vdRes.SellOrders = sellBook
	vdRes.BuyOrders = buyBook
	vdRes.Message = "Success"
	return vdRes, nil
}
