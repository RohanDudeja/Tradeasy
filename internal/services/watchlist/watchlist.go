package watchlist

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"errors"
)

func CreateWatchlist(CrReq CreateRequest) (CrRes CreateResponse, err error) {
	var wl model.Watchlist
	wl.Name=CrReq.WatchlistName
	err = config.DB.Table("watchlist").Create(&wl).Error
	if err !=nil {
		return CrRes,errors.New("failed to create watchlist")
	}
	CrRes.Message="Watchlist created"
	CrRes.WatchlistId=wl.Id

	return CrRes, nil
}

func AddStockEntry(AddReq AddStockRequest, watchlistId int) (AddRes AddStockResponse, err error) {
	var uwl model.UserWatchlist
	uwl.Userid=AddReq.UserId
	uwl.WatchlistId=watchlistId
	uwl.StockName=AddReq.StockName

	err = config.DB.Table("user_watchlist").Create(uwl).Error
	if err !=nil {
		return AddRes,errors.New("stock not added")
	}
	AddRes.Message="Stock added"
	return AddRes, nil
}

func DeleteStockEntry(DelReq DeleteStockRequest, watchlistId int) (DelRes DeleteStockResponse, err error) {
	var uwl model.UserWatchlist
	err=config.DB.Table("user_watchlist").Where("userid = ? AND watchlist_id = ? AND stock_name = ?", DelReq.UserId, watchlistId,DelReq.StockName).Delete(uwl).Error

	if err !=nil {
		return DelRes,errors.New("stock not deleted")
	}
	DelRes.Message="Stock deleted"
	return DelRes, nil
}
func SortWatchlist(SortReq SortRequest) (SortRes SortResponse, err error) {
	var wl []model.UserWatchlist
	err=config.DB.Raw("SELECT * FROM user_watchlist ORDER BY stock_name").Scan(&wl).Error
	if err !=nil {
		return SortRes,errors.New("watchlist not sorted")
	}
	SortRes.SortedWatchlist=wl
	return SortRes, nil
}
