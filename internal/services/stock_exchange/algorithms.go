package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"log"
	"math"
	"time"
)

var orderResponse = make(chan OrderResponse, 2)

func UpdateLTP(ltp int, stock string) {
	currentStock := model.Stocks{}
	err := config.DB.Table("stocks").Where("stock_ticker_symbol = ?", stock).Find(&currentStock).Error
	if err != nil {
		log.Println(err.Error())
	}
	currentStock.HighPrice = int(math.Max(float64(ltp), float64(currentStock.HighPrice)))
	currentStock.LowPrice = int(math.Min(float64(ltp), float64(currentStock.LowPrice)))
	currentStock.LTP = ltp
	currentStock.UpdatedAt = time.Now()
	err = config.DB.Table("stocks").Where("stock_ticker_symbol = ?", stock).Updates(&currentStock).Error
	if err != nil {
		log.Println(err.Error())
	}
}

func GetLTP(stock string) (int, error) {
	var currStockLTP []model.Stocks
	err := config.DB.Raw("SELECT * FROM stocks WHERE stock_ticker_symbol = ?", stock).Scan(&currStockLTP).Error
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return currStockLTP[0].LTP, nil
}

func UpdateMarketOrderPrices(stock string) {

	ltp, err := GetLTP(stock)
	if err != nil {
		log.Println(err.Error())
	}
	err = config.DB.Exec("UPDATE buy_order_book SET order_price = ? WHERE stock_ticker_symbol = ? AND order_type =?", ltp, stock, "Market").Error
	if err != nil {
		log.Println(err.Error())
	}
	err = config.DB.Exec("UPDATE sell_order_book SET order_price = ? WHERE stock_ticker_symbol = ? AND order_type =?", ltp, stock, "Market").Error
	if err != nil {
		log.Println(err.Error())
	}
}

func SendResponse(stock string, status string, message string, orderId string, price int, quantity int) {
	resp := OrderResponse{}
	resp.StockName = stock
	resp.Status = status
	resp.Message = message
	resp.OrderID = orderId
	resp.AveragePrice = price
	resp.OrderExecutionTime = time.Now()
	resp.Quantity = quantity
	orderResponse <- resp
}

func CancelAtExpiry() {
	var buyOrders []model.BuyOrderBook
	err := config.DB.Raw("SELECT * FROM buy_order_book").Scan(&buyOrders).Error
	if err != nil {
		log.Println(err.Error())
	}
	for _, order := range buyOrders {
		SendResponse(order.StockTickerSymbol, "CANCELLED", "Order Executed", order.OrderID, order.OrderPrice, order.OrderQuantity)
		err = config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", order.ID).Error
		if err != nil {
			log.Println(err.Error())
		}
	}
	var sellOrders []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book").Scan(&sellOrders).Error
	if err != nil {
		log.Println(err.Error())
	}
	for _, order := range sellOrders {
		SendResponse(order.StockTickerSymbol, "CANCELLED", "Order Executed", order.OrderID, order.OrderPrice, order.OrderQuantity)
		err = config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", order.ID).Error
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func BuyLimitOrder(buyOrderBody OrderRequest, sellBook []model.SellOrderBook) {

	ltp := 0
	for _, elem := range sellBook {
		if elem.OrderPrice <= buyOrderBody.LimitPrice {
			if elem.OrderQuantity > buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				err := config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, buyOrderBody.StockName)
				UpdateMarketOrderPrices(buyOrderBody.StockName)
				SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(buyOrderBody.StockName, "PARTIAL", "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				buyOrderBody.Quantity -= buyOrderBody.Quantity
				break
			} else if elem.OrderQuantity < buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				SendResponse(buyOrderBody.StockName, "PARTIAL", "Order Executed Partially", buyOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
				buyOrderBody.Quantity -= elem.OrderQuantity
				err := config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", buyOrderBody.Quantity, time.Now(), buyOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, buyOrderBody.StockName)
				UpdateMarketOrderPrices(buyOrderBody.StockName)
			} else {
				ltp = elem.OrderPrice
				SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				buyOrderBody.Quantity -= elem.OrderQuantity
				err := config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, buyOrderBody.StockName)
				UpdateMarketOrderPrices(buyOrderBody.StockName)
			}
		}
	}
}

func BuyMarketOrder(buyOrderBody OrderRequest, sellBook []model.SellOrderBook) {

	ltp := 0
	for _, elem := range sellBook {
		if elem.OrderQuantity > buyOrderBody.Quantity {
			ltp = elem.OrderPrice
			err := config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-buyOrderBody.Quantity, time.Now(), elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, buyOrderBody.StockName)
			UpdateMarketOrderPrices(buyOrderBody.StockName)
			SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(buyOrderBody.StockName, "PARTIAL", "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			buyOrderBody.Quantity -= buyOrderBody.Quantity
			break
		} else if elem.OrderQuantity < buyOrderBody.Quantity {
			ltp = elem.OrderPrice
			SendResponse(buyOrderBody.StockName, "PARTIAL", "Order Executed Partially", buyOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
			buyOrderBody.Quantity -= elem.OrderQuantity
			err := config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", buyOrderBody.Quantity, time.Now(), buyOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, buyOrderBody.StockName)
			UpdateMarketOrderPrices(buyOrderBody.StockName)
		} else {
			ltp = elem.OrderPrice
			SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(buyOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			buyOrderBody.Quantity -= elem.OrderQuantity
			err := config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, buyOrderBody.StockName)
			UpdateMarketOrderPrices(buyOrderBody.StockName)
		}
	}
}
func SellLimitOrder(sellOrderBody OrderRequest, buyBook []model.BuyOrderBook) {

	ltp := 0
	for _, elem := range buyBook {

		if elem.OrderPrice >= sellOrderBody.LimitPrice {
			if elem.OrderQuantity > sellOrderBody.Quantity {
				ltp = elem.OrderPrice
				err := config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? ,updated_at = ?  WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, sellOrderBody.StockName)
				UpdateMarketOrderPrices(sellOrderBody.StockName)
				SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(sellOrderBody.StockName, "PARTIAL", "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				sellOrderBody.Quantity -= sellOrderBody.Quantity
				break
			} else if elem.OrderQuantity < sellOrderBody.Quantity {
				ltp = elem.OrderPrice

				SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(sellOrderBody.StockName, "PARTIAL", "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
				sellOrderBody.Quantity -= elem.OrderQuantity
				err := config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", sellOrderBody.Quantity, time.Now(), sellOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, sellOrderBody.StockName)
				UpdateMarketOrderPrices(sellOrderBody.StockName)
			} else {
				ltp = elem.OrderPrice
				SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				sellOrderBody.Quantity -= elem.OrderQuantity
				err := config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID).Error
				if err != nil {
					log.Println(err.Error())
				}
				err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
				if err != nil {
					log.Println(err.Error())
				}
				UpdateLTP(ltp, sellOrderBody.StockName)
				UpdateMarketOrderPrices(sellOrderBody.StockName)
			}
		}
	}

}
func SellMarketOrder(sellOrderBody OrderRequest, buyBook []model.BuyOrderBook) {

	ltp := 0
	for _, elem := range buyBook {
		if elem.OrderQuantity > sellOrderBody.Quantity {
			ltp = elem.OrderPrice
			err := config.DB.Exec("UPDATE buy_order_book SET order_quantity = ? , updated_at = ? WHERE id = ? ", elem.OrderQuantity-sellOrderBody.Quantity, time.Now(), elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, sellOrderBody.StockName)
			UpdateMarketOrderPrices(sellOrderBody.StockName)
			SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(sellOrderBody.StockName, "PARTIAL", "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			sellOrderBody.Quantity -= sellOrderBody.Quantity
			break
		} else if elem.OrderQuantity < sellOrderBody.Quantity {
			ltp = elem.OrderPrice
			SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(sellOrderBody.StockName, "PARTIAL", "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
			sellOrderBody.Quantity -= elem.OrderQuantity
			err := config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("UPDATE sell_order_book SET order_quantity = ? , updated_at = ? WHERE order_id = ?", sellOrderBody.Quantity, time.Now(), sellOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, sellOrderBody.StockName)
			UpdateMarketOrderPrices(sellOrderBody.StockName)
		} else {
			ltp = elem.OrderPrice
			SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(sellOrderBody.StockName, "COMPLETED", "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			sellOrderBody.Quantity -= elem.OrderQuantity
			err := config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", elem.ID).Error
			if err != nil {
				log.Println(err.Error())
			}
			err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
			if err != nil {
				log.Println(err.Error())
			}
			UpdateLTP(ltp, sellOrderBody.StockName)
			UpdateMarketOrderPrices(sellOrderBody.StockName)
		}
	}
}

func BuyOrderMatchingAlgo(buyOrderBody OrderRequest) {

	var sellBook []model.SellOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC", buyOrderBody.StockName).Scan(&sellBook).Error

	if err != nil || len(sellBook) == 0 {
		// abort
		err := config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
		if err != nil {
			log.Println(err.Error())
		}
		SendResponse(buyOrderBody.StockName, "FAILED", "Couldn't execute order", buyOrderBody.OrderID, buyOrderBody.LimitPrice, buyOrderBody.Quantity)
		return
	}

	if buyOrderBody.OrderType == "Limit" {
		BuyLimitOrder(buyOrderBody, sellBook)
	} else {
		//Market  order
		BuyMarketOrder(buyOrderBody, sellBook)
	}
	//order complete
	return
}

func SellOrderMatchingAlgo(sellOrderBody OrderRequest) {
	var buyBook []model.BuyOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC ", sellOrderBody.StockName).Scan(&buyBook).Error

	if err != nil || len(buyBook) == 0 {
		// abort
		err := config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
		if err != nil {
			log.Println(err.Error())
		}
		SendResponse(sellOrderBody.StockName, "FAILED", "Couldn't execute order", sellOrderBody.OrderID, sellOrderBody.LimitPrice, sellOrderBody.Quantity)
		return
	}

	if sellOrderBody.OrderType == "Limit" {
		SellLimitOrder(sellOrderBody, buyBook)
	} else {
		//Market  order
		SellMarketOrder(sellOrderBody, buyBook)
	}
	//order complete
	return
}

// BuyOrder ...Update Buy Order actions on the StockExchange database
func BuyOrder(buyOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "PENDING"
	resp.OrderID = buyOrderBody.OrderID
	resp.StockName = buyOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.BuyOrderBook{
		OrderID:           buyOrderBody.OrderID,
		StockTickerSymbol: buyOrderBody.StockName,
		OrderQuantity:     buyOrderBody.Quantity,
		OrderStatus:       "PENDING",
		OrderPrice:        buyOrderBody.LimitPrice,
		OrderType:         buyOrderBody.OrderType,
		CreatedAt:         buyOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if buyOrderBody.OrderType != "Market" && buyOrderBody.OrderType != "Limit" {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if buyOrderBody.OrderType == "Limit" && buyOrderBody.LimitPrice == 0 {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(buyOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, nil
	}
	if buyOrderBody.OrderType == "Market" {
		newEntry.OrderPrice = ltp
	}
	err = config.DB.Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, err
	}
	go BuyOrderMatchingAlgo(buyOrderBody)
	return resp, nil
}

// SellOrder ...Update Sell Order actions on the StockExchange database
func SellOrder(sellOrderBody OrderRequest) (resp OrderResponse, err error) {

	resp.Status = "PENDING"
	resp.OrderID = sellOrderBody.OrderID
	resp.StockName = sellOrderBody.StockName
	resp.Message = "Order Received"
	newEntry := model.SellOrderBook{
		OrderID:           sellOrderBody.OrderID,
		StockTickerSymbol: sellOrderBody.StockName,
		OrderQuantity:     sellOrderBody.Quantity,
		OrderStatus:       "PENDING",
		OrderPrice:        sellOrderBody.LimitPrice,
		OrderType:         sellOrderBody.OrderType,
		CreatedAt:         sellOrderBody.OrderPlacedTime,
		UpdatedAt:         time.Now(),
	}
	if sellOrderBody.OrderType != "Market" && sellOrderBody.OrderType != "Limit" {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order type"
		return resp, nil
	}
	if sellOrderBody.OrderType == "Limit" && sellOrderBody.LimitPrice == 0 {
		resp.Status = "CANCELLED"
		resp.Message = "Incorrect order price"
		return resp, nil
	}
	ltp, err := GetLTP(sellOrderBody.StockName)
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, nil
	}
	if sellOrderBody.OrderType == "Market" {
		newEntry.OrderPrice = ltp
	}
	err = config.DB.Create(&newEntry).Error
	if err != nil {
		log.Println(err.Error())
		resp.Status = "FAILED"
		resp.Message = "Error in db fetch"
		return resp, err
	}
	go SellOrderMatchingAlgo(sellOrderBody)
	return resp, nil
}

// DeleteBuyOrder ...Update Delete Buy Order actions on the StockExchange database
func DeleteBuyOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", orderId).Error
	if err != nil {
		deleteRes.Message = "Failed"
		deleteRes.Success = false
		return deleteRes, err
	}
	deleteRes.Message = "Success"
	deleteRes.Success = true
	return deleteRes, nil
}

// DeleteSellOrder ...Update Delete Sell Order actions on the StockExchange database
func DeleteSellOrder(orderId string) (deleteRes DeleteResponse, err error) {
	err = config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", orderId).Error
	if err != nil {
		deleteRes.Success = false
		deleteRes.Message = "Failed"
		return deleteRes, err
	}
	deleteRes.Success = true
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
