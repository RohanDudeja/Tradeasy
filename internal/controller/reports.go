package controller

import (
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/reports"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DailyPendingOrders(c *gin.Context){
	var rep []model.PendingOrders
	id := c.Params.ByName("id")
	if err := c.BindJSON(&rep); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	penOrderRes, err := reports.DailyPendingOrders(rep, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, penOrderRes)
	}
}
func Portfolio(c *gin.Context){
	var rep []model.Holdings
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	if err := c.BindJSON(&rep); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	PortfolioRes, err := reports.Portfolio(rep, id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, PortfolioRes)
	}
}
func OrdersHistory(c *gin.Context){
	var rep1 []model.OrderHistory
	var rep2 []model.Holdings
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	if err := c.BindJSON(&rep1); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if err := c.BindJSON(&rep2); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	ordHisRes ,err := reports.OrdersHistory(rep1, rep2, id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, ordHisRes)
	}
}
func ProfitLossHistory(c *gin.Context){
	var rep []model.OrderHistory
	id := c.Params.ByName("id")
	from := c.Params.ByName("From")
	to := c.Params.ByName("To")
	if err := c.BindJSON(&rep); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	proLosRes, err := reports.ProfitLossHistory(rep, id, from, to)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, proLosRes)
	}
}

