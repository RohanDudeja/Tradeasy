package controller

import (
	"Tradeasy/internal/services/order"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyOrder(c *gin.Context) {
	var BReq order.BuyRequest
	id := c.Params.ByName("Userid")
	c.BindJSON(&BReq)
	BReq.UserId = id
	BRes, err := order.BuyOrder(BReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &BRes)
	}
}

func SellOrder(c *gin.Context) {
	var SReq order.SellRequest
	id := c.Params.ByName("Userid")
	c.BindJSON(&SReq)
	SReq.UserId = id

	SRes, err := order.SellOrder(SReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &SRes)
	}
}

func CancelOrder(c *gin.Context) {
	id := c.Params.ByName("OrderId")
	CRes, err := order.CancelOrder(id)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &CRes)
	}
}
