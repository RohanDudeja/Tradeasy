package controller

import (
	"Tradeasy/internal/services/order"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyOrder(c *gin.Context) {
	var bReq order.BuyRequest
	id := c.Params.ByName("user_id")
	c.BindJSON(&bReq)
	bReq.UserId = id
	bRes, err := order.BuyOrder(bReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, bRes)
	}
}

func SellOrder(c *gin.Context) {
	var sReq order.SellRequest
	id := c.Params.ByName("user_id")
	c.BindJSON(&sReq)
	sReq.UserId = id

	sRes, err := order.SellOrder(sReq)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, sRes)
	}
}

func CancelOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	cRes, err := order.CancelOrder(id)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, cRes)
	}
}
