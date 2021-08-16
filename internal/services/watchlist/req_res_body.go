package watchlist

import "Tradeasy/internal/model"

type CreateRequest struct {
	WatchlistName string `json:"watchlist_name"`
	UserId        string `json:"user_id"`
}

type CreateResponse struct {
	WatchlistId int    `json:"watchlist_id"`
	Message     string `json:"message"`
}

type AddStockRequest struct {
	UserId    string `json:"user_id"`
	StockName string `json:"stock_name"`
}

type AddStockResponse struct {
	Message string `json:"message"`
}
type DeleteStockRequest struct {
	UserId    string `json:"user_id"`
	StockName string `json:"stock_name"`
}

type DeleteStockResponse struct {
	Message string `json:"message"`
}
type SortRequest struct {
	UserId int `json:"user_id"`
}

type SortResponse struct {
	SortedWatchlist []model.UserWatchlist `json:"sorted_watchlist"`
}
