package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/stock_exchange"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type HoldingsQuantity struct {
	StockName     string
	TotalQuantity int
}

func BuyOrder(bReq BuyRequest) (bRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed
	var account model.TradingAccount

	if bReq.BookType == Market {
		if err = config.DB.Table("stocks_feed").Where("stock_name=?", bReq.StockName).Last(&stocks).Error; err != nil {
			return bRes, err
		}
		bReq.LimitPrice = stocks.LTP
	}

	if err = config.DB.Table("trading_account").Where("user_id=?", bReq.UserId).First(&account).Error; err != nil {
		return bRes, err
	}

	if account.Balance < int64(bReq.Quantity*bReq.LimitPrice) {
		return bRes, errors.New("balance is insufficient for the placed order")
	}

	account.Balance = account.Balance - int64(bReq.Quantity*bReq.LimitPrice)
	if err = config.DB.Table("trading_account").Where("user_id=?", bReq.UserId).Updates(&account).Error; err != nil {
		return bRes, err
	}

	orderId := uuid.New().String()
	p := model.PendingOrders{
		UserId:    bReq.UserId,
		OrderId:   orderId,
		StockName: bReq.StockName,
		OrderType: Buy,
		BookType:  bReq.BookType,
		Quantity:  bReq.Quantity,
		Status:    Pending,
	}

	if bReq.BookType == Market {
		p.OrderPrice = bReq.LimitPrice
	} else if bReq.BookType == Limit {
		p.LimitPrice = bReq.LimitPrice
	}

	exeOrder := stock_exchange.OrderRequest{
		OrderID:    orderId,
		StockName:  bReq.StockName,
		Quantity:   bReq.Quantity,
		OrderType:  bReq.BookType,
		LimitPrice: bReq.LimitPrice,
	}

	request, err := json.Marshal(exeOrder)
	if err != nil {
		return bRes, err
	}
	response, err := http.Post(BuyOrderURL, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return bRes, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &bRes)
	if err != nil {
		return bRes, err
	}

	if bRes.Status == Pending {
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return bRes, err
		}
	}

	return bRes, nil
}

func SellOrder(sReq SellRequest) (sRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed

	var r HoldingsQuantity
	if err = config.DB.Table("holdings").Select("stock_name, sum(quantity) as total_quantity").Where("user_id=? AND stock_name=?", sReq.UserId, sReq.StockName).Group("stock_name").Scan(&r).Error; err != nil {
		return sRes, err
	} else if r.TotalQuantity < sReq.Quantity {
		return sRes, errors.New("sell Order quantity is higher than holdings quantity")
	}

	if sReq.BookType == Market {
		if err = config.DB.Table("stocks_feed").Where("stock_name=?", sReq.StockName).Last(&stocks).Error; err != nil {
			return sRes, err
		}
		sReq.LimitPrice = stocks.LTP
	}

	orderId := uuid.New().String()
	p := model.PendingOrders{
		UserId:    sReq.UserId,
		OrderId:   orderId,
		StockName: sReq.StockName,
		OrderType: Sell,
		BookType:  sReq.BookType,
		Quantity:  sReq.Quantity,
		Status:    Pending,
	}

	if sReq.BookType == Market {
		p.OrderPrice = sReq.LimitPrice
	} else if sReq.BookType == Limit {
		p.LimitPrice = sReq.LimitPrice
	}

	exeOrder := stock_exchange.OrderRequest{
		OrderID:    orderId,
		StockName:  sReq.StockName,
		Quantity:   sReq.Quantity,
		OrderType:  sReq.BookType,
		LimitPrice: sReq.LimitPrice,
	}

	request, err := json.Marshal(exeOrder)
	if err != nil {
		return sRes, err
	}
	response, err := http.Post(SellOrderURL, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return sRes, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &sRes)
	if err != nil {
		return sRes, err
	}

	if sRes.Status == Pending {
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return sRes, err
		}
	}

	return sRes, nil
}

func CancelOrder(id string) (cRes CancelResponse, err error) {
	var p model.PendingOrders
	if err = config.DB.Table("pending_orders").Where("order_id=?", id).First(&p).Error; err != nil {
		return cRes, err
	}
	url := ""
	if p.OrderType == Buy {
		url = BuyOrderURL + "/" + id
	} else {
		url = SellOrderURL + "/" + id
	}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return cRes, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return cRes, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	var dRes stock_exchange.DeleteResponse
	err = json.Unmarshal(body, &dRes)
	if err != nil {
		return cRes, err
	}

	cRes.UserId = p.UserId
	cRes.OrderId = p.OrderId
	cRes.StockName = p.StockName
	if dRes.Success == true {
		p.Status = Cancelled
		cRes.Status = Cancelled
		cRes.Message = dRes.Message
		if p.OrderType == Buy {
			var account model.TradingAccount
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return cRes, err
			}
			if p.BookType == Market {
				account.Balance = account.Balance + int64(p.Quantity*p.OrderPrice)
			} else if p.BookType == Limit {
				account.Balance = account.Balance + int64(p.Quantity*p.LimitPrice)
			}
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return cRes, err
			}
		}
		if err = config.DB.Table("pending_orders").Where("order_id=?", id).Updates(&p).Error; err != nil {
			return cRes, err
		}
		if err = config.DB.Table("pending_orders").Where("order_id=?", id).Delete(&p).Error; err != nil {
			return cRes, err
		}
	} else {
		cRes.Message = dRes.Message
		cRes.Status = Failed
	}
	return cRes, nil
}
