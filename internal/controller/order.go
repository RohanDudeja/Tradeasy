package controller

import (
	"Tradeasy/internal/services/order"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BuyOrder(c *gin.Context) {
	var bReq order.BuyRequest
	id := c.Params.ByName("user_id")
	err:=c.BindJSON(&bReq)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	bReq.UserId = id
	bRes, err := order.BuyOrder(bReq)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, bRes)
	}
}

func SellOrder(c *gin.Context) {
	var sReq order.SellRequest
	id := c.Params.ByName("user_id")
	err:=c.BindJSON(&sReq)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	sReq.UserId = id

	sRes, err := order.SellOrder(sReq)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, sRes)
	}
}

func CancelOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	cRes, err := order.CancelOrder(id)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, cRes)
	}
}
