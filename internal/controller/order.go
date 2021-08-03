package controller

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyOrder(c *gin.Context)  {
	var p model.PendingOrders
	id := c.Params.ByName("Userid")
	c.BindJSON(&p)
	p.OrderType="Buy"
	p.Userid=id
	err:=services.BuyOrder(&p)

	if err !=nil{
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else {
		c.JSON(http.StatusOK,p)
	}
}

func SellOrder(c *gin.Context)  {
	var p model.PendingOrders
	id := c.Params.ByName("Userid")
	c.BindJSON(&p)
	p.OrderType="Sell"
	p.Userid=id

	err:=services.SellOrder(&p)

	if err !=nil{
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else {
		c.JSON(http.StatusOK,p)
	}
}

func CancelOrder(c *gin.Context)  {
	id := c.Params.ByName("OrderId")
	err:=services.CancelOrder(id)

	if err !=nil{
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else {
		c.JSON(http.StatusOK,c)
	}
}