package controller

import (
	"Tradeasy/internal/model/stock_exchange"
	"Tradeasy/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StockExchangeActions interface {
	ExecuteBuyOrder(c *gin.Context)
	ExecuteSellOrder(c *gin.Context)
	DeleteBuyOrder(c *gin.Context)
	DeleteSellOrder(c *gin.Context)
}
type StockExchange struct {
	Stocks   stock_exchange.Stocks
	BuyBook  stock_exchange.BuyOrderBook
	SellBook stock_exchange.SellOrderBook
}

func (se *StockExchange)ExecuteBuyOrder(c *gin.Context) {
	var BuyOrder stock_exchange.BuyOrderBook
	c.BindJSON(&BuyOrder)
	err,msg := services.ExecuteBuyOrder(se,&BuyOrder)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK,[]byte(msg))
	}
}
func (se *StockExchange)ExecuteSellOrder(c *gin.Context) {
	var SellOrder stock_exchange.SellOrderBook
	c.BindJSON(&SellOrder)
	err,msg := services.ExecuteSellOrder(se,&SellOrder)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK,[]byte(msg))
	}
}
func (se *StockExchange)DeleteBuyOrder(c *gin.Context) {
	var BuyOrder stock_exchange.BuyOrderBook
	c.BindJSON(&BuyOrder)
	err,msg := services.DeleteBuyOrder(se,&BuyOrder)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK,[]byte(msg))
	}

}
func (se *StockExchange)DeleteSellOrder(c *gin.Context) {
	var SellOrder stock_exchange.SellOrderBook
	c.BindJSON(&SellOrder)
	err,msg := services.DeleteSellOrder(se,&SellOrder)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK,[]byte(msg))
	}
}