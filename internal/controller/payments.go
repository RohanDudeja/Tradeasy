package controller

import (
	"Tradeasy/internal/services/payments"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddAmount(c *gin.Context) {
	var addReq payments.AddRequest
	userId := c.Params.ByName("user_id")
	if err := c.BindJSON(&addReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	addRes, err := payments.AddAmount(addReq, userId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, addRes)
	}
}
func WithdrawAmount(c *gin.Context) {
	var withdrawReq payments.WithdrawRequest
	userId := c.Params.ByName("user_id")
	if err := c.BindJSON(&withdrawReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	withdrawRes, err := payments.WithdrawAmount(withdrawReq, userId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, withdrawRes)
	}
}

func Callback(c *gin.Context) {
	var callbackParamRequest payments.CallbackParamRequest
	if err := c.BindQuery(&callbackParamRequest); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	callBackResponse, err := payments.Callback(callbackParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, callBackResponse)
	}
}
