package stock_exchange

import (
	model "Tradeasy/internal/model/stock_exchange"
	"time"
)

type OrderRequest struct {
	OrderID         string    `json:"order_id"`
	StockName       string    `json:"stock_name"`
	OrderPlacedTime time.Time `json:"order_placed_time"`
	OrderType       string    `json:"order_type"`
	LimitPrice      int       `json:"limit_price"`
	Quantity        int       `json:"quantity"`
	IsDummy         bool      `json:"is_dummy"`
}

type OrderResponse struct {
	OrderID            string    `json:"order_id"`
	StockName          string    `json:"stock_name"`
	AveragePrice       int       `json:"average_price"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	OrderExecutionTime time.Time `json:"order_execution_time"`
	Message            string    `json:"message"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ViewDepthResponse struct {
	BuyOrders  []model.BuyOrderBook //top 5 buy order details
	SellOrders []model.SellOrderBook
	Message    string
}

type StockDetails struct {
	StockName string    `json:"stock_name"`
	LTP       int       `json:"ltp"`
	UpdatedAt time.Time `json:"updated_at"`
	HighPrice int       `json:"high_price"`
	OpenPrice int       `json:"open_price"`
	LowPrice  int       `json:"low_price"`
}
