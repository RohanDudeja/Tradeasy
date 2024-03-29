package order

const (
	Pending      = "PENDING"
	Completed    = "COMPLETED"
	Partial      = "PARTIAL"
	Cancelled    = "CANCELLED"
	Failed       = "FAILED"
	Market       = "Market"
	Limit        = "Limit"
	Buy          = "Buy"
	Sell         = "Sell"
	BuyOrderURL  = "http://localhost:8080/buy_order_book/buy_order"
	SellOrderURL = "http://localhost:8080/sell_order_book/sell_order"
)

type BuyRequest struct {
	UserId     string `json:"user_id"`
	StockName  string `json:"stock_name"`
	BookType   string `json:"book_type"`
	LimitPrice int    `json:"limit_price"`
	Quantity   int    `json:"quantity"`
}

type BuyResponse struct {
	StockName  string `json:"stock_name"`
	BookType   string `json:"book_type"`
	LimitPrice int    `json:"limit_price"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
	Status     int    `json:"status"`
	OrderPrice int    `json:"order_price"`
	Message    string `json:"message"`
}

type SellResponse struct {
	StockName  string `json:"stock_name"`
	BookType   string `json:"book_type"`
	LimitPrice int    `json:"limit_price"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
	Status     int    `json:"status"`
	OrderPrice int    `json:"order_price"`
	Message    string `json:"message"`
}

type SellRequest struct {
	UserId     string `json:"user_id"`
	StockName  string `json:"stock_name"`
	BookType   string `json:"book_type"`
	LimitPrice int    `json:"limit_price"`
	Quantity   int    `json:"quantity"`
}

type CancelResponse struct {
	UserId    string `json:"user_id"`
	OrderId   string `json:"order_id"`
	StockName string `json:"stock_name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
