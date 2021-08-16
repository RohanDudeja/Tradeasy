package order

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/services/stock_exchange"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HoldingsQuantity struct {
	StockName     string
	TotalQuantity int
}

func BuyOrder(bReq BuyRequest) (bRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed
	var account model.TradingAccount

	if bReq.BookType == Market {
		if err = database.GetDB().Table("stocks_feed").Where("stock_name=?", bReq.StockName).Last(&stocks).Error; err != nil {
			log.Println("Error in Fetching stock feed", err)
			return bRes, errors.New("given stock name is not found")
		}
		bReq.LimitPrice = stocks.LTP
	}

	if err = database.GetDB().Table("trading_account").Where("user_id=?", bReq.UserId).First(&account).Error; err != nil {
		log.Println("Error in Fetching Trading account", err)
		return bRes, errors.New("users trading account is not found")
	}

	if account.Balance < int64(bReq.Quantity*bReq.LimitPrice) {
		log.Println("Balance in Trading Account is not sufficient")
		return bRes, errors.New("balance is insufficient for the placed order")
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
		log.Println("Error in Marshalling the Executive buy order", err)
		return bRes, err
	}
	req, err := http.NewRequest("POST", BuyOrderURL, bytes.NewBuffer(request))
	if err != nil {
		log.Println("Error while making a request", err)
		return bRes, errors.New("error in making a request with Stock Exchange engine")
	}
	userName := config.GetConfig().StockExchange.Authentication.UserName
	password := config.GetConfig().StockExchange.Authentication.Password
	req.SetBasicAuth(userName, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in getting response for buying order", err)
		return bRes, errors.New("error in getting response for placing buy order")
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &bRes)
	if err != nil {
		log.Println("Error in Unmarshalling the Executive buy order response", err)
		return bRes, err
	}

	if bRes.Status == Pending {
		account.Balance = account.Balance - int64(bReq.Quantity*bReq.LimitPrice)
		account.UpdatedAt = time.Now()
		if err = database.GetDB().Table("trading_account").Where("user_id=?", bReq.UserId).Updates(&account).Error; err != nil {
			log.Println("Error in Updating Balance in Trading Account", err)
			return bRes, errors.New("error in Updating Balance in Trading Account")
		}
		if err = database.GetDB().Table("pending_orders").Create(&p).Error; err != nil {
			log.Println("Error in Creating pending orders", err)
			return bRes, errors.New("error while adding data into database")
		}
	}
	return bRes, nil
}

func SellOrder(sReq SellRequest) (sRes stock_exchange.OrderResponse, err error) {
	var stocks model.StocksFeed

	var r HoldingsQuantity
	if err = database.GetDB().Table("holdings").Select("stock_name, sum(quantity) as total_quantity").
		Where("user_id=? AND stock_name=?", sReq.UserId, sReq.StockName).Group("stock_name").
		Scan(&r).Error; err != nil {

		log.Println("Error in Fetching Total Quantities from Holdings", err)
		return sRes, err
	} else if r.TotalQuantity < sReq.Quantity {
		return sRes, errors.New("sell Order quantity is higher than holdings quantity")
	}

	if sReq.BookType == Market {
		if err = database.GetDB().Table("stocks_feed").Where("stock_name=?", sReq.StockName).Last(&stocks).Error; err != nil {
			log.Println("Error in Fetching Stocks feed", err)
			return sRes, errors.New("stock name not found")
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
		log.Println("Error in Marshalling the Executive sell order", err)
		return sRes, err
	}
	req, err := http.NewRequest("POST", SellOrderURL, bytes.NewBuffer(request))
	if err != nil {
		log.Println("Error while making request", err)
		return sRes, errors.New("error while making request with Stock Exchange engine")
	}
	userName := config.GetConfig().StockExchange.Authentication.UserName
	password := config.GetConfig().StockExchange.Authentication.Password
	req.SetBasicAuth(userName, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in getting response for selling order", err)
		return sRes, errors.New("error in getting response for placing sell order")
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &sRes)
	if err != nil {
		log.Println("Error in Unmarshalling the Executive sell order response", err)
		return sRes, err
	}

	if sRes.Status == Pending {
		if err = database.GetDB().Table("pending_orders").Create(&p).Error; err != nil {
			log.Println("Error in Creating Pending Orders", err)
			return sRes, errors.New("error while adding data into database")
		}
	}

	return sRes, nil
}

func CancelOrder(id string) (cRes CancelResponse, err error) {
	var p model.PendingOrders
	if err = database.GetDB().Table("pending_orders").Where("order_id=?", id).First(&p).Error; err != nil {
		log.Println("Error in Fetching Pending Orders", err)
		return cRes, errors.New("order id is not found for cancelling the order")
	}
	url := ""
	if p.OrderType == Buy {
		url = BuyOrderURL + "/" + id
	} else {
		url = SellOrderURL + "/" + id
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("Error in Making request for cancelling order", err)
		return cRes, errors.New("error in Making request for cancelling order")
	}
	userName := config.GetConfig().StockExchange.Authentication.UserName
	password := config.GetConfig().StockExchange.Authentication.Password
	req.SetBasicAuth(userName, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in getting response for cancelling order", err)
		return cRes, errors.New("error in getting response for cancelling order")
	}
	body, _ := ioutil.ReadAll(res.Body)
	var dRes stock_exchange.DeleteResponse
	err = json.Unmarshal(body, &dRes)
	if err != nil {
		log.Println("Error in Unmarshalling the Executive cancel order", err)
		return cRes, err
	}

	cRes.UserId = p.UserId
	cRes.OrderId = p.OrderId
	cRes.StockName = p.StockName
	if dRes.Success == true {
		p.Status = Cancelled
		p.UpdatedAt = time.Now()
		cRes.Status = Cancelled
		cRes.Message = dRes.Message
		if p.OrderType == Buy {
			var account model.TradingAccount
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).First(&account).Error; err != nil {
				log.Println("Error in fetching trading account", err)
				return cRes, errors.New("users trading account not found - incorrect user id")
			}
			if p.BookType == Market {
				account.Balance = account.Balance + int64(p.Quantity*p.OrderPrice)
			} else if p.BookType == Limit {
				account.Balance = account.Balance + int64(p.Quantity*p.LimitPrice)
			}
			account.UpdatedAt = time.Now()
			if err = database.GetDB().Table("trading_account").Where("user_id=?", p.UserId).Updates(&account).Error; err != nil {
				log.Println("Error in updating balance in trading account", err)
				return cRes, errors.New("error in updating the balance for trading account")
			}
		}
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", id).Updates(&p).Error; err != nil {
			log.Println("Error in updating status in pending orders", err)
			return cRes, errors.New("error in updating the database")
		}
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", id).Delete(&p).Error; err != nil {
			log.Println("Error in deleting order in pending orders", err)
			return cRes, errors.New("error in updating the database")
		}
	} else {
		cRes.Message = dRes.Message
		cRes.Status = Failed
	}
	return cRes, nil
}
