package reports

type DailyPendingOrderResponse struct {
	Userid     string `json:"user_id"`
	OrderId    string `json:"order_id"`
	StockName  string `json:"stock_name"`
	OrderType  string `json:"order_type"`
	BookType   string `json:"book_type"`
	LimitPrice int    `json:"limit_price"`
	Quantity   int    `json:"quantity"`
	OrderPrice int    `json:"order_price"`
	Status     string `json:"status"`
}
type PortfolioResponse struct {
	Userid    string `json:"user_id"`
	OrderId   string `json:"order_id"`
	StockName string `json:"stock_name"`
	Quantity  int    `json:"quantity"`
	BuyPrice  int    `json:"buy_price"`
}

type OrderHistoryResponse struct {
	Userid    string `json:"user_id"`
	OrderId   string `json:"order_id"`
	StockName string `json:"stock_name"`
	Quantity  int    `json:"quantity"`
	BuyPrice  int    `json:"buy_price"`
	SellPrice int    `json:"sell_price"`
}

type ProfitLossHistoryResponse struct {
	Userid     string `json:"user_id"`
	OrderId    string `json:"order_id"`
	StockName  string `json:"stock_name"`
	Quantity   int    `json:"quantity"`
	BuyPrice   int    `json:"buy_price"`
	SellPrice  int    `json:"sell_price"`
	ProfitLoss int    `json:"profit_loss"`
}
