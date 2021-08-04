package controller

import (
	"Tradeasy/internal/services/Order"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyOrder(c *gin.Context) {
	var BReq Order.BuyRequest
	id := c.Params.ByName("Userid")
	c.BindJSON(&BReq)
	BReq.UserId = id
	err, bres := Order.BuyOrder(BReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &bres)
	}
}

func SellOrder(c *gin.Context) {
	var SReq Order.SellRequest
	id := c.Params.ByName("Userid")
	c.BindJSON(&SReq)
	SReq.UserId = id

	err, sres := Order.SellOrder(SReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &sres)
	}
}

func CancelOrder(c *gin.Context) {
	id := c.Params.ByName("OrderId")
	err, cres := Order.CancelOrder(id)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &cres)
	}
}
