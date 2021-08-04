package controller

import (
	"Tradeasy/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyOrder(c *gin.Context) {
	var breq services.BuyRequest
	var bres services.BuyResponse
	id := c.Params.ByName("Userid")
	c.BindJSON(&breq)
	breq.UserId = id
	err := services.BuyOrder(&breq, &bres)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &bres)
	}
}

func SellOrder(c *gin.Context) {
	var sreq services.SellRequest
	var sres services.SellResponse
	id := c.Params.ByName("Userid")
	c.BindJSON(&sreq)
	sreq.UserId = id

	err := services.SellOrder(&sreq, &sres)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &sres)
	}
}

func CancelOrder(c *gin.Context) {
	var cres services.CancelResponse
	id := c.Params.ByName("OrderId")
	err := services.CancelOrder(id, &cres)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &cres)
	}
}
