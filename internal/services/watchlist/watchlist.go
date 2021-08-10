package watchlist

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"errors"
)

func CreateWatchlist(CrReq CreateRequest) (CrRes CreateResponse, err error) {
	//var user model.Users
	//err = config.DB.Table("users").Where("user_id = ?", CrReq.UserId).First(&user).Error
	//if err != nil {
	//	return CrRes, errors.New("user not found")
	//}
	var wl model.Watchlist
	er := config.DB.Table("watchlist").Where("name = ?", CrReq.WatchlistName).First(&wl).Error
	if er == nil {
		return CrRes, errors.New("watchlist name already exists")
	}
	wl.Name = CrReq.WatchlistName
	wl.Userid = CrReq.UserId
	err = config.DB.Table("watchlist").Create(&wl).Error
	if err != nil {
		return CrRes, errors.New("failed to create watchlist")
	}
	CrRes.Message = "Watchlist created"
	CrRes.WatchlistId = wl.Id

	return CrRes, nil
}

func AddStockEntry(AddReq AddStockRequest, watchlistId int) (AddRes AddStockResponse, err error) {
	var wl model.Watchlist
	err = config.DB.Table("watchlist").Where("user_id = ? AND watchlist_id = ? ", AddReq.UserId, watchlistId).First(&wl).Error
	if err != nil {
		return AddRes, errors.New("user not found")
	}
	var uwl model.UserWatchlist

	er := config.DB.Table("user_watchlist").Where("stock_name = ? AND watchlist_id = ?", AddReq.StockName, watchlistId).First(&uwl).Error
	if er == nil {
		return AddRes, errors.New("stock name already exists")
	}
	uwl.Userid = AddReq.UserId
	uwl.WatchlistId = watchlistId
	uwl.StockName = AddReq.StockName
	err = config.DB.Table("user_watchlist").Create(&uwl).Error
	if err != nil {
		return AddRes, errors.New("stock not added")
	}
	AddRes.Message = "Stock added"
	return AddRes, nil
}

func DeleteStockEntry(DelReq DeleteStockRequest, watchlistId int) (DelRes DeleteStockResponse, err error) {
	var uwl model.UserWatchlist
	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", DelReq.UserId, watchlistId, DelReq.StockName).First(&uwl).Error
	if err != nil {
		return DelRes, errors.New("stock not found")
	}

	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", DelReq.UserId, watchlistId, DelReq.StockName).Delete(&uwl).Error
	if err != nil {
		return DelRes, errors.New("stock not deleted")
	}
	DelRes.Message = "Stock deleted"
	return DelRes, nil
}
func SortWatchlist(SortReq SortRequest) (SortRes SortResponse, err error) {
	var wl []model.UserWatchlist
	err = config.DB.Raw("SELECT * FROM user_watchlist ORDER BY stock_name ASC").Scan(&wl).Error
	if err != nil {
		return SortRes, errors.New("watchlist not sorted")
	}
	SortRes.SortedWatchlist = wl
	return SortRes, nil
}
