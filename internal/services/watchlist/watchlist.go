package watchlist

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"errors"
)

func CreateWatchlist(Req CreateRequest) (Res CreateResponse, err error) {
	var user model.Users
	err = config.DB.Table("users").Where("user_id = ?", Req.UserId).First(&user).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	var wl model.Watchlist
	err = config.DB.Table("watchlist").Where("name = ?", Req.WatchlistName).First(&wl).Error
	if err == nil {
		return Res, errors.New("watchlist name already exists")
	}
	wl.Name = Req.WatchlistName
	wl.Userid = Req.UserId
	err = config.DB.Table("watchlist").Create(&wl).Error
	if err != nil {
		return Res, errors.New("failed to create watchlist")
	}
	Res.Message = "Watchlist created"
	Res.WatchlistId = wl.Id

	return Res, nil
}

func AddStockEntry(Req AddStockRequest, watchlistId int) (Res AddStockResponse, err error) {
	var wl model.Watchlist
	err = config.DB.Table("watchlist").Where("user_id = ?", Req.UserId).First(&wl).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	err = config.DB.Table("watchlist").Where("watchlist_id = ? ", watchlistId).First(&wl).Error
	if err != nil {
		return Res, errors.New("watchlist not found")
	}
	var uwl model.UserWatchlist

	err = config.DB.Table("user_watchlist").Where("stock_name = ? AND watchlist_id = ?", Req.StockName, watchlistId).First(&uwl).Error
	if err == nil {
		return Res, errors.New("stock name already exists")
	}
	uwl.Userid = Req.UserId
	uwl.WatchlistId = watchlistId
	uwl.StockName = Req.StockName
	err = config.DB.Table("user_watchlist").Create(&uwl).Error
	if err != nil {
		return Res, errors.New("stock not added")
	}
	Res.Message = "Stock added"
	return Res, nil
}

func DeleteStockEntry(Req DeleteStockRequest, watchlistId int) (Res DeleteStockResponse, err error) {
	var uwl model.UserWatchlist
	err = config.DB.Table("user_watchlist").Where("user_id = ?", Req.UserId).First(&uwl).Error
	if err != nil {
		return Res, errors.New("user not found")
	}
	err = config.DB.Table("user_watchlist").Where("watchlist_id = ? ", watchlistId).First(&uwl).Error
	if err != nil {
		return Res, errors.New("watchlist not found")
	}
	err = config.DB.Table("user_watchlist").Where("stock_name = ?", Req.StockName).First(&uwl).Error
	if err != nil {
		return Res, errors.New("stock not found")
	}

	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", Req.UserId, watchlistId, Req.StockName).Delete(&uwl).Error
	if err != nil {
		return Res, errors.New("stock not deleted")
	}
	Res.Message = "Stock deleted"
	return Res, nil
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
