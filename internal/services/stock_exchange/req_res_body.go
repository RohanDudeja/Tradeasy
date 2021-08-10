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
}

type OrderResponse struct {
	OrderID            string    `json:"order_id"`
	StockName          string    `json:"stock_name"`
	AveragePrice       int       `json:"average_price"`
	Status             string    `json:"status"`
	Quantity           int       `json:"quantity"`
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
	High      int       `json:"high"`
	Open      int       `json:"open"`
	Low       int       `json:"low"`
}
