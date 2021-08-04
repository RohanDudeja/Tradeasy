package stock_exchange

import "time"

type OrderRequest struct {
	OrderID         string    `json:"order_id"`
	StockName       string    `json:"stock_name"`
	OrderPlacedTime time.Time `json:"order_placed_time"`
	OrderType       string    `json:"order_type"`
	LimitPrice      uint      `json:"limit_price"`
	Quantity        uint      `json:"quantity"`
}

type OrderResponse struct {
	OrderID            string    `json:"order_id"`
	StockName          string    `json:"stock_name"`
	AveragePrice       uint      `json:"average_price"`
	Status             string    `json:"status"`
	OrderExecutionTime time.Time `json:"order_execution_time"`
	Message            string    `json:"message"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type ViewDepthResponse struct {
	BuyOrders  []string //top 5 buy order details
	SellOrders []string
	Message    string
}
