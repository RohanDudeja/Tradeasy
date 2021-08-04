package controller

import (
	"Tradeasy/internal/services/payments"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddAmount(c *gin.Context){
	var addReq payments.AddRequest
	id := c.Params.ByName("id")
	c.BindJSON(&addReq)
	addRes, err := payments.AddAmount(addReq,id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, addRes)
	}
}
func WithdrawAmount(c *gin.Context){
	var withdrawReq payments.WithdrawRequest
	id := c.Params.ByName("id")
	c.BindJSON(&withdrawReq)
	withdrawRes, err := payments.WithdrawAmount(withdrawReq,id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, withdrawRes)
	}
}
