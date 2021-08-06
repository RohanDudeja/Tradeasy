package controller

import (
	"Tradeasy/internal/services/reports"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DailyPendingOrders(c *gin.Context) {
	id := c.Params.ByName("id")
	penOrderRes, err := reports.DailyPendingOrders(id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, penOrderRes)
	}
}
func Portfolio(c *gin.Context) {
	id := c.Params.ByName("id")
	from := c.Params.ByName("from")
	to := c.Params.ByName("to")
	PortfolioRes, err := reports.Portfolio(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, PortfolioRes)
	}
}
func OrdersHistory(c *gin.Context) {
	id := c.Params.ByName("id")
	from := c.Params.ByName("from")
	to := c.Params.ByName("to")
	ordHisRes, err := reports.OrdersHistory(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, ordHisRes)
	}
}
func ProfitLossHistory(c *gin.Context) {
	id := c.Params.ByName("id")
	from := c.Params.ByName("from")
	to := c.Params.ByName("to")
	proLosRes, err := reports.ProfitLossHistory(id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, proLosRes)
	}
}
