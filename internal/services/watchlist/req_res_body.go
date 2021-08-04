package watchlist

type CreateRequest struct {
	WatchlistName string `json:"watchlist_name"`
}

type CreateResponse struct {
	WatchlistId int    `json:"watchlist_id"`
	Message     string `json:"message"`
}

type AddStockRequest struct {
	StockName string `json:"stock_name"`
}

type AddStockResponse struct {
	Message string `json:"message"`
}
type DeleteStockRequest struct {
	StockName string `json:"stock_name"`
}

type DeleteStockResponse struct {
	Message string `json:"message"`
}
type SortRequest struct {
	WatchlistId int `json:"watchlist_id"`
}

type SortResponse struct {
	Message string `json:"message"`
}
