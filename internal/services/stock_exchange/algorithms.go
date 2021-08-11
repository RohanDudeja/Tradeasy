package stock_exchange

import (
	"Tradeasy/config"
	model "Tradeasy/internal/model/stock_exchange"
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"time"
)

var orderResponse = make(chan OrderResponse, 2)

// UpdateLTP ...Updates ltp, high and low price for a stock
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

// GetLTP ...FetchesLTP of a stock
func GetLTP(stock string) (int, error) {
	var currStockLTP []model.Stocks
	err := config.DB.Raw("SELECT * FROM stocks WHERE stock_ticker_symbol = ?", stock).Scan(&currStockLTP).Error
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return currStockLTP[0].LTP, nil
}

// UpdateMarketOrderPrices ...Updates prices of market orders
func UpdateMarketOrderPrices(stock string) {

	ltp, err := GetLTP(stock)
	if err != nil {
		log.Println(err.Error())
	}
	err = config.DB.Exec("UPDATE buy_order_book SET order_price = ? WHERE stock_ticker_symbol = ? AND order_type =?", ltp, stock, Market).Error
	if err != nil {
		log.Println(err.Error())
	}
	err = config.DB.Exec("UPDATE sell_order_book SET order_price = ? WHERE stock_ticker_symbol = ? AND order_type =?", ltp, stock, Market).Error
	if err != nil {
		log.Println(err.Error())
	}
}

// SendResponse ...Sends responseto response channel
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

// CancelAtExpiry ...Cancels all book orders at expiry
func CancelAtExpiry() {
	var buyOrders []model.BuyOrderBook
	err := config.DB.Raw("SELECT * FROM buy_order_book").Scan(&buyOrders).Error
	if err != nil {
		log.Println(err.Error())
	}
	for _, order := range buyOrders {
		err = config.DB.Exec("DELETE FROM buy_order_book WHERE id = ?", order.ID).Error
		if err != nil {
			log.Println(err.Error())
		} else {
			SendResponse(order.StockTickerSymbol, Cancelled, "Expiry Time Reached", order.OrderID, order.OrderPrice, order.OrderQuantity)
		}
	}
	var sellOrders []model.SellOrderBook
	err = config.DB.Raw("SELECT * FROM sell_order_book").Scan(&sellOrders).Error
	if err != nil {
		log.Println(err.Error())
	}
	for _, order := range sellOrders {
		err = config.DB.Exec("DELETE FROM sell_order_book WHERE id = ?", order.ID).Error
		if err != nil {
			log.Println(err.Error())
		} else {
			SendResponse(order.StockTickerSymbol, Cancelled, "Expiry Time Reached", order.OrderID, order.OrderPrice, order.OrderQuantity)
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
				SendResponse(buyOrderBody.StockName, Completed, "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(buyOrderBody.StockName, Partial, "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				buyOrderBody.Quantity -= buyOrderBody.Quantity
				break
			} else if elem.OrderQuantity < buyOrderBody.Quantity {
				ltp = elem.OrderPrice
				SendResponse(buyOrderBody.StockName, Partial, "Order Executed Partially", buyOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(buyOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
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
				SendResponse(buyOrderBody.StockName, Completed, "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
				SendResponse(buyOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
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
			SendResponse(buyOrderBody.StockName, Completed, "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(buyOrderBody.StockName, Partial, "Order Executed Partially", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			buyOrderBody.Quantity -= buyOrderBody.Quantity
			break
		} else if elem.OrderQuantity < buyOrderBody.Quantity {
			ltp = elem.OrderPrice
			SendResponse(buyOrderBody.StockName, Partial, "Order Executed Partially", buyOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(buyOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
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
			SendResponse(buyOrderBody.StockName, Completed, "Order Executed", buyOrderBody.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
			SendResponse(buyOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, buyOrderBody.Quantity)
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
				SendResponse(sellOrderBody.StockName, Completed, "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(sellOrderBody.StockName, Partial, "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				sellOrderBody.Quantity -= sellOrderBody.Quantity
				break
			} else if elem.OrderQuantity < sellOrderBody.Quantity {
				ltp = elem.OrderPrice

				SendResponse(sellOrderBody.StockName, Completed, "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
				SendResponse(sellOrderBody.StockName, Partial, "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
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
				SendResponse(sellOrderBody.StockName, Completed, "Order Executed", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
				SendResponse(sellOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
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
			SendResponse(sellOrderBody.StockName, Completed, "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(sellOrderBody.StockName, Partial, "Order Executed Partially", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			sellOrderBody.Quantity -= sellOrderBody.Quantity
			break
		} else if elem.OrderQuantity < sellOrderBody.Quantity {
			ltp = elem.OrderPrice
			SendResponse(sellOrderBody.StockName, Completed, "Order Executed Partially", elem.OrderID, elem.OrderPrice, elem.OrderQuantity)
			SendResponse(sellOrderBody.StockName, Partial, "Order Executed Partially", sellOrderBody.OrderID, elem.OrderPrice, elem.OrderQuantity)
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
			SendResponse(sellOrderBody.StockName, Completed, "Order Executed", sellOrderBody.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
			SendResponse(sellOrderBody.StockName, Completed, "Order Executed", elem.OrderID, elem.OrderPrice, sellOrderBody.Quantity)
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

// BuyOrderMatchingAlgo  ...Matching algorithm for buy orders
func BuyOrderMatchingAlgo(buyOrderBody OrderRequest) {

	var sellBook []model.SellOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM sell_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price ASC,created_at ASC", buyOrderBody.StockName).Scan(&sellBook).Error

	if err != nil {
		// abort
		err := config.DB.Exec("DELETE FROM buy_order_book WHERE order_id = ?", buyOrderBody.OrderID).Error
		if err != nil {
			log.Println(err.Error())
		}
		SendResponse(buyOrderBody.StockName, Failed, "Couldn't execute order", buyOrderBody.OrderID, buyOrderBody.LimitPrice, buyOrderBody.Quantity)
		return
	}

	if buyOrderBody.OrderType == Limit {
		BuyLimitOrder(buyOrderBody, sellBook)
	} else if buyOrderBody.OrderType == Market {
		BuyMarketOrder(buyOrderBody, sellBook)
	}
	//order complete
	return
}

// SellOrderMatchingAlgo  ...Matching algorithm for sell orders
func SellOrderMatchingAlgo(sellOrderBody OrderRequest) {
	var buyBook []model.BuyOrderBook
	// db lock
	err := config.DB.Raw("SELECT * FROM buy_order_book WHERE stock_ticker_symbol = ?  ORDER BY order_price DESC,created_at ASC ", sellOrderBody.StockName).Scan(&buyBook).Error

	if err != nil {
		// abort
		err := config.DB.Exec("DELETE FROM sell_order_book WHERE order_id = ?", sellOrderBody.OrderID).Error
		if err != nil {
			log.Println(err.Error())
		}
		SendResponse(sellOrderBody.StockName, Failed, "Couldn't execute order", sellOrderBody.OrderID, sellOrderBody.LimitPrice, sellOrderBody.Quantity)
		return
	}

	if sellOrderBody.OrderType == Limit {
		SellLimitOrder(sellOrderBody, buyBook)
	} else if sellOrderBody.OrderType == Market {
		SellMarketOrder(sellOrderBody, buyBook)
	}
	//order complete
	return
}

// RandomizerAlgo ...Generates Random traffic to fluctuate ltp of stocks
func RandomizerAlgo() {

	for {
		var allStocks []model.Stocks
		err := config.DB.Table("stocks").Find(&allStocks).Error
		if err != nil {
			log.Println(err.Error())
		}
		orderType := []string{Limit, Market}
		for _, stock := range allStocks {

			//placing buy order
			orderID := uuid.New().String()
			rand.Seed(time.Now().UnixNano())
			idx := rand.Intn(2)
			order := orderType[idx]
			min := stock.LTP - int(float64(stock.LTP)*PercentChange)
			max := stock.LTP + int(float64(stock.LTP)*PercentChange)
			buyOrderBody := OrderRequest{
				OrderID:         orderID,
				StockName:       stock.StockName,
				OrderPlacedTime: time.Time{},
				OrderType:       order,
				LimitPrice:      rand.Intn(max-min+1) + min,
				Quantity:        rand.Intn(OrdersQuantityRange) + 1,
			}
			_, err := BuyOrder(buyOrderBody)
			if err != nil {
				log.Println(err.Error())
				return
			}

			//placing sell order
			orderID = uuid.New().String()
			rand.Seed(time.Now().UnixNano())
			idx = rand.Intn(2)
			order = orderType[idx]
			min = stock.LTP - int(float64(stock.LTP)*PercentChange)
			max = stock.LTP + int(float64(stock.LTP)*PercentChange)
			time.Sleep(1 * time.Second)
			sellOrderBody := OrderRequest{
				OrderID:         orderID,
				StockName:       stock.StockName,
				OrderPlacedTime: time.Time{},
				OrderType:       order,
				LimitPrice:      rand.Intn(max-min+1) + min,
				Quantity:        rand.Intn(OrdersQuantityRange) + 1,
			}
			_, err = SellOrder(sellOrderBody)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
		// sleep and run again
		time.Sleep(5 * time.Second)
	}
}
