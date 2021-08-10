package controller

import (
	"Tradeasy/internal/services/reports"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func DailyPendingOrders(c *gin.Context) {
	id := c.Params.ByName("user_id")
	penOrderRes, err := reports.DailyPendingOrders(id)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		c.JSON(http.StatusOK, penOrderRes)
	}
}
func Portfolio(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Fatalln(err.Error())
	}
	PortfolioRes, err := reports.Portfolio(id, reportsParamRequest)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		c.JSON(http.StatusOK, PortfolioRes)
	}
}
func OrdersHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Fatalln(err.Error())
	}
	ordHisRes, err := reports.OrdersHistory(id, reportsParamRequest)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		c.JSON(http.StatusOK, ordHisRes)
	}
}
func ProfitLossHistory(c *gin.Context) {
	var reportsParamRequest reports.ReportsParamRequest
	id := c.Params.ByName("user_id")
	if err := c.BindQuery(&reportsParamRequest); err != nil {
		log.Fatalln(err.Error())
	}
	proLosRes, err := reports.ProfitLossHistory(id, reportsParamRequest)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		c.JSON(http.StatusOK, proLosRes)
	}
}
