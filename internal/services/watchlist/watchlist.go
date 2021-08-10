package watchlist

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"errors"
)

func CreateWatchlist(req CreateRequest) (res CreateResponse, err error) {
	var user model.Users
	err = config.DB.Table("users").Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	type Result struct {
		Userid      string `json:"user_id"`
		Name        string `json:"name"`
		WatchlistId int    `json:"watchlist_id"`
	}
	var result []Result
	err = config.DB.Raw("SELECT user_watchlist.user_id,user_watchlist.watchlist_id,watchlist.name FROM user_watchlist INNER JOIN watchlist ON user_watchlist.watchlist_id = watchlist.watchlist_id WHERE user_watchlist.user_id = ? AND watchlist.name = ?", req.UserId, req.WatchlistName).Scan(&result).Error
	if err != nil {
		return res, err
	}
	if len(result) != 0 {
		return res, errors.New("watchlist name already found")
	}
	var wl model.Watchlist
	wl.Name = req.WatchlistName
	err = config.DB.Table("watchlist").Create(&wl).Error
	if err != nil {
		return res, errors.New("failed to create watchlist")
	}
	var uwl model.UserWatchlist
	uwl.WatchlistId = wl.Id
	uwl.Userid = req.UserId
	err = config.DB.Table("user_watchlist").Create(&uwl).Error
	if err != nil {
		return res, errors.New("failed to create watchlist")
	}
	res.Message = "Watchlist created"
	res.WatchlistId = wl.Id

	return res, nil
}

func AddStockEntry(req AddStockRequest, watchlistId int) (res AddStockResponse, err error) {
	var uwl model.UserWatchlist
	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ?", req.UserId, watchlistId).First(&uwl).Error
	if err != nil {
		return res, errors.New("user not found")
	}

	err = config.DB.Table("user_watchlist").Where("user_id = ? AND stock_name = ? AND watchlist_id = ?", req.UserId, req.StockName, watchlistId).First(&uwl).Error
	if err == nil {
		return res, errors.New("stock name already exists")
	}

	var stockNames []string
	err = config.DB.Table("user_watchlist").Select("stock_name").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, "").Scan(&stockNames).Error
	if err != nil {
		return res, err
	}
	if len(stockNames) == 0 {
		uwl.Userid = req.UserId
		uwl.WatchlistId = watchlistId
		uwl.StockName = req.StockName
		err = config.DB.Table("user_watchlist").Create(&uwl).Error

		if err != nil {
			return res, errors.New("stock not added")
		}
	} else {
		config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, "").Update("stock_name", req.StockName)
	}
	res.Message = "Stock added"
	return res, nil
}

func DeleteStockEntry(req DeleteStockRequest, watchlistId int) (res DeleteStockResponse, err error) {
	var wl model.Watchlist
	err = config.DB.Table("watchlist").Where("watchlist_id = ?", watchlistId).First(&wl).Error
	if err != nil {
		return res, errors.New("watchlist not found")
	}
	var uwl model.UserWatchlist
	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ?", req.UserId, watchlistId).First(&uwl).Error
	if err != nil {
		return res, errors.New("user not found")
	}

	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, req.StockName).First(&uwl).Error
	if err != nil {
		return res, errors.New("stock not found")
	}

	err = config.DB.Table("user_watchlist").Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, req.StockName).Delete(&uwl).Error
	if err != nil {
		return res, errors.New("stock not deleted")
	}
	res.Message = "Stock deleted"
	return res, nil
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
