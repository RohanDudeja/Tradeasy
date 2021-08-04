package services

type BuyRequest struct {
	UserId string `json:"user_id"`
	StockName string `json:"stock_name"`
	BookType string `json:"book_type"`
	LimitPrice int `json:"limit_price"`
	Quantity string `json:"quantity"`
}

type BuyResponse struct {
	StockName string `json:"stock_name"`
	BookType string `json:"book_type"`
	LimitPrice int `json:"limit_price"`
	Quantity string `json:"quantity"`
	TotalPrice int `json:"total_price"`
	Status int `json:"status"`
	OrderPrice int `json:"order_price"`
	Message string `json:"message"`
}

type SellResponse struct {
	StockName string `json:"stock_name"`
	BookType string `json:"book_type"`
	LimitPrice int `json:"limit_price"`
	Quantity string `json:"quantity"`
	TotalPrice int `json:"total_price"`
	Status int `json:"status"`
	OrderPrice int `json:"order_price"`
	Message string `json:"message"`
}

type SellRequest struct {
	UserId string `json:"user_id"`
	StockName string `json:"stock_name"`
	BookType string `json:"book_type"`
	LimitPrice int `json:"limit_price"`
	Quantity string `json:"quantity"`
}

type CancelResponse struct {
	UserId string `json:"user_id"`
	OrderId string `json:"order_id"`
	StockName string `json:"stock_name"`
	Status string `json:"status"`
	Message string `json:"message"`
}