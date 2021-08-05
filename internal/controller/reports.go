package controller

import (
	"Tradeasy/internal/services/reports"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DailyPendingOrders(c *gin.Context){
	id := c.Params.ByName("id")
	penOrderRes, err := reports.DailyPendingOrders(id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, penOrderRes)
	}
}
func Portfolio(c *gin.Context){
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	PortfolioRes, err := reports.Portfolio(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, PortfolioRes)
	}
}
func OrdersHistory(c *gin.Context){
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	ordHisRes ,err := reports.OrdersHistory(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, ordHisRes)
	}
}
func ProfitLossHistory(c *gin.Context){
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	proLosRes, err := reports.ProfitLossHistory(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, proLosRes)
	}
}

