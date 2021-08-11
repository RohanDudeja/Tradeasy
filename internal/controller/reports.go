package controller

import (
	"Tradeasy/internal/services/reports"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DailyPendingOrders(c *gin.Context) {
	id := c.Params.ByName("user_id")
	pendingOrderResponse, err := reports.DailyPendingOrders(id)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, pendingOrderResponse)
	}
}
func Portfolio(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}
	PortfolioResponse, err := reports.Portfolio(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, PortfolioResponse)
	}
}
func OrdersHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}
	orderHisResponse, err := reports.OrdersHistory(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, orderHisResponse)
	}
}
func ProfitLossHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}
	profitLossResponse, err := reports.ProfitLossHistory(id, reportsParamRequest)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, profitLossResponse)
	}
}
