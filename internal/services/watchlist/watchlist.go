package watchlist

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"errors"
	"log"
)

func CreateWatchlist(req CreateRequest) (res CreateResponse, err error) {
	var (
		user          model.Users
		watchlist     model.Watchlist
		userWatchlist model.UserWatchlist
	)

	err = database.GetDB().Table("users").Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	err = database.GetDB().Table("watchlist").Where("name = ?", req.WatchlistName).First(&watchlist).Error
	if err != nil {
		watchlist.Name = req.WatchlistName
		err = database.GetDB().Table("watchlist").Create(&watchlist).Error
		if err != nil {
			return res, errors.New("watchlist not created")
		}
	}
	err = database.GetDB().Table("user_watchlist").Where("user_id = ? AND watchlist_id = ?", req.UserId, watchlist.Id).
		First(&userWatchlist).Error
	if err == nil {
		return res, errors.New("watchlist name already exists")
	}
	userWatchlist.Userid = req.UserId
	userWatchlist.WatchlistId = watchlist.Id

	err = database.GetDB().Table("user_watchlist").Create(&userWatchlist).Error
	if err != nil {
		return res, errors.New("user watchlist not created")
	}

	res.Message = "Watchlist created"
	res.WatchlistId = watchlist.Id

	return res, nil
}

func AddStockEntry(req AddStockRequest, watchlistId int) (res AddStockResponse, err error) {
	var userWatchlist model.UserWatchlist
	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND watchlist_id = ?", req.UserId, watchlistId).
		First(&userWatchlist).Error
	if err != nil {
		return res, errors.New("user not found")
	}

	userWatchlist = model.UserWatchlist{}
	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND stock_name = ? AND watchlist_id = ?", req.UserId, req.StockName, watchlistId).
		First(&userWatchlist).Error
	if err == nil {
		return res, errors.New("stock name already exists")
	}

	userWatchlist = model.UserWatchlist{}
	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, "").
		First(&userWatchlist).Error
	if err != nil && userWatchlist.StockName != "" {
		var newUserWatchlist model.UserWatchlist
		newUserWatchlist.Userid = req.UserId
		newUserWatchlist.WatchlistId = watchlistId
		newUserWatchlist.StockName = req.StockName
		err = database.GetDB().Table("user_watchlist").Create(&newUserWatchlist).Error
		if err != nil {
			return res, errors.New("stock not added")
		}
	} else if err == nil {
		userWatchlist.StockName = req.StockName
		err = database.GetDB().Table("user_watchlist").
			Update(&userWatchlist).Error
		if err != nil {
			return res, errors.New("stock not added")
		}
	} else {
		log.Println(err)
		return res, err
	}
	res.Message = "Stock added"
	return res, nil
}

func DeleteStockEntry(req DeleteStockRequest, watchlistId int) (res DeleteStockResponse, err error) {

	var userWatchlist model.UserWatchlist
	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND watchlist_id = ?", req.UserId, watchlistId).First(&userWatchlist).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	userWatchlist = model.UserWatchlist{}
	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, req.StockName).
		First(&userWatchlist).Error
	if err != nil {
		return res, errors.New("stock not found")
	}

	err = database.GetDB().Table("user_watchlist").
		Where("user_id = ? AND watchlist_id = ? AND stock_name = ?", req.UserId, watchlistId, req.StockName).
		Delete(&userWatchlist).Error
	if err != nil {
		return res, errors.New("stock not deleted")
	}
	res.Message = "Stock deleted"
	return res, nil
}
func SortWatchlist(SortReq SortRequest, watchlistId int) (SortRes SortResponse, err error) {
	var wl []model.UserWatchlist
	err = database.GetDB().Raw("SELECT * FROM user_watchlist WHERE user_id = ? AND watchlist_id = ? ORDER BY stock_name ASC", SortReq.UserId, watchlistId).Scan(&wl).Error
	if err != nil {
		return SortRes, errors.New("watchlist not sorted")
	}
	SortRes.SortedWatchlist = wl
	return SortRes, nil
}
