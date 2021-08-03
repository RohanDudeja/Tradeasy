package controller

import (
	"Tradeasy/internal/model/transaction_engine"
	"Tradeasy/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)
func AbortMsg(code int, err error, c *gin.Context) {
	c.String(code, "Oops! Please retry.")
	c.Error(err)
	c.Abort()
}
func BuyOrder(c *gin.Context)  {
	var p transaction_engine.PendingOrders
	id := c.Params.ByName("Userid")
	c.BindJSON(&p)
	p.OrderType="Buy"
	p.Userid=id

	err:=services.BuyOrder(&p)

	if err !=nil{
		AbortMsg(500,err,c)
	}else {
		c.JSON(http.StatusOK,p)
	}
}

func SellOrder(c *gin.Context)  {
	var p transaction_engine.PendingOrders
	id := c.Params.ByName("Userid")
	c.BindJSON(&p)
	p.OrderType="Sell"
	p.Userid=id

	err:=services.SellOrder(&p)

	if err !=nil{
		AbortMsg(500,err,c)
	}else {
		c.JSON(http.StatusOK,p)
	}
}

func CancelOrder(c *gin.Context)  {
	id := c.Params.ByName("OrderId")
	err:=services.CancelOrder(id)

	if err !=nil{
		AbortMsg(500,err,c)
	}else {
		c.JSON(http.StatusOK,c)
	}
}