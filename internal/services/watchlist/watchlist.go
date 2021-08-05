package watchlist

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
)

func CreateWatchlist(wl* model.Watchlist,CrReq CreateRequest) (CrRes CreateResponse, err error) {

	CrRes.Message="Watchlist created"
	err = config.DB.Table("watchlist").Create(wl).Error
	return CrRes, err
}

func AddStockEntry(uwl* model.UserWatchlist,AddReq AddStockRequest, watchlistId string) (AddRes AddStockResponse, err error) {
	uwl.Userid=AddReq.UserId
	uwl.WatchlistId=watchlistId

	err = config.DB.Table("user_watchlist").Create(uwl).Error
	AddRes.Message="Stock added"
	return AddRes, err
}

func DeleteStockEntry(uwl* model.UserWatchlist,DelReq DeleteStockRequest, watchlistId string) (DelRes DeleteStockResponse, err error) {

	config.DB.Table("user_watchlist").Where("userid = ? AND watchlist_id = ? AND stock_name = ?", DelReq.UserId, watchlistId,DelReq.StockName).Delete(uwl)
	DelRes.Message="Stock deleted"
	return DelRes, err
}
func SortWatchlist(SortReq SortRequest) (SortRes SortResponse, err error) {
	var wl []model.UserWatchlist
	config.DB.Raw("SELECT * FROM user_watchlist ORDER BY stock_name").Scan(&wl)

	SortRes.SortedWatchlist=wl
	return SortRes, err
}
