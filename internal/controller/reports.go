package controller

import (
	"Tradeasy/internal/services/reports"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DailyPendingOrders(c *gin.Context) {
	id := c.Params.ByName("user_id")
	penOrderRes, err := reports.DailyPendingOrders(id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, penOrderRes)
	}
}
func Portfolio(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	PortfolioRes, err := reports.Portfolio(id, reportsParamRequest.From, reportsParamRequest.To)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, PortfolioRes)
	}
}
func OrdersHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	ordHisRes, err := reports.OrdersHistory(id, reportsParamRequest.From, reportsParamRequest.To)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, ordHisRes)
	}
}
func ProfitLossHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	proLosRes, err := reports.ProfitLossHistory(id, reportsParamRequest.From, reportsParamRequest.To)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, proLosRes)
	}
}
