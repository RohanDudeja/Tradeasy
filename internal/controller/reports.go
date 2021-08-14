package controller

import (
	"Tradeasy/internal/services/reports"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DailyPendingOrders(c *gin.Context) {
	id := c.Params.ByName("user_id")
	response, err := reports.DailyPendingOrders(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
func Portfolio(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in binding query parameters"})
	}
	response, err := reports.Portfolio(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
func OrdersHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in binding query parameters"})
	}
	response, err := reports.OrdersHistory(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
func ProfitLossHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in binding query parameters"})
	}
	response, err := reports.ProfitLossHistory(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
