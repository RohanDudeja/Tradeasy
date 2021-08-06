package watchlist

import "Tradeasy/internal/model"

type CreateRequest struct {
	WatchlistName string `json:"watchlist_name"`
	UserId        string `json:"userid"`
}

type CreateResponse struct {
	WatchlistId int    `json:"watchlist_id"`
	Message     string `json:"message"`
}

type AddStockRequest struct {
	UserId    string `json:"userid"`
	StockName string `json:"stock_name"`
}

type AddStockResponse struct {
	Message string `json:"message"`
}
type DeleteStockRequest struct {
	UserId    string `json:"userid"`
	StockName string `json:"stock_name"`
}

type DeleteStockResponse struct {
	Message string `json:"message"`
}
type SortRequest struct {
	WatchlistId int `json:"watchlist_id"`
}

type SortResponse struct {
	SortedWatchlist []model.UserWatchlist `json:"sorted_watchlist"`
}
