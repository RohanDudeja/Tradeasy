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
	"time"
)
type result struct {
	StockName     string
	TotalQuantity int
}

func BuyOrder(bReq BuyRequest) (bRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed
	var account model.TradingAccount

	if bReq.BookType == market {
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
	//account.UpdatedAt=time.Now()
	if err = config.DB.Table("trading_account").Where("user_id=?", bReq.UserId).Updates(&account).Error; err != nil {
		return bRes, err
	}

	orderId := uuid.New().String()
	p := model.PendingOrders{
		UserId:    bReq.UserId,
		OrderId:   orderId,
		StockName: bReq.StockName,
		OrderType: buy,
		BookType:  bReq.BookType,
		Quantity:  bReq.Quantity,
		Status:    pending,
	}

	if bReq.BookType == market {
		p.OrderPrice = bReq.LimitPrice
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return bRes, err
		}
	} else {
		p.LimitPrice = bReq.LimitPrice
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return bRes, err
		}
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
	response, err := http.Post("http://localhost:8080/buy_order_book/buy_order", "application/json", bytes.NewBuffer(request))
	if err != nil {
		return bRes, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &bRes)
	if err != nil {
		return bRes, err
	}
	return bRes, nil
}

func SellOrder(sReq SellRequest) (sRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed

	var r result
	if err = config.DB.Table("holdings").Select("stock_name, sum(quantity) as total_quantity").Where("user_id=? AND stock_name=?", sReq.UserId, sReq.StockName).Group("stock_name").Scan(&r).Error; err != nil {
		return sRes, err
	} else if r.TotalQuantity < sReq.Quantity {
		return sRes, errors.New("sell Order quantity is higher than holdings quantity")
	}

	if sReq.BookType == "Market" {
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
		OrderType: sell,
		BookType:  sReq.BookType,
		Quantity:  sReq.Quantity,
		Status:    pending,
	}

	if sReq.BookType == market {
		p.OrderPrice = sReq.LimitPrice
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return sRes, err
		}
	} else {
		p.LimitPrice = sReq.LimitPrice
		if err = config.DB.Table("pending_orders").Create(&p).Error; err != nil {
			return sRes, err
		}
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
	response, err := http.Post("http://localhost:8080/sell_order_book/sell_order", "application/json", bytes.NewBuffer(request))
	if err != nil {
		return sRes, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &sRes)
	if err != nil {
		return sRes, err
	}
	return sRes, nil
}

func CancelOrder(id string) (cRes CancelResponse, err error) {
	var p model.PendingOrders
	if err = config.DB.Table("pending_orders").Where("order_id=?", id).First(&p).Error; err != nil {
		return cRes, err
	}
	url := "http://localhost:8080/"
	if p.OrderType == "Buy" {
		url = url + "/buy_order_book/buy_order/" + id
	} else {
		url = url + "/sell_order_book/sell_order/" + id
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
		p.Status = cancelled
		cRes.Status = cancelled
		cRes.Message = dRes.Message
		if p.OrderType == buy {
			var account model.TradingAccount
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				return cRes, err
			}
			if p.BookType == market {
				account.Balance = account.Balance + int64(p.Quantity*p.OrderPrice)
			} else {
				account.Balance = account.Balance + int64(p.Quantity*p.LimitPrice)
			}
			account.UpdatedAt=time.Now()
			if err = config.DB.Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				return cRes, err
			}
		}
		p.UpdatedAt=time.Now()
		if err = config.DB.Table("pending_orders").Where("order_id=?", id).Updates(&p).Error; err != nil {
			return cRes, err
		}
		if err = config.DB.Table("pending_orders").Where("order_id=?", id).Delete(&p).Error; err != nil {
			return cRes, err
		}
	} else {
		cRes.Message=dRes.Message
		cRes.Status=failed
	}
	return cRes, nil
}
