package middleware

import (
	"Tradeasy/config"
	"github.com/gin-gonic/gin"
)

func ExchangeBasicAuth() gin.HandlerFunc {
	Auth := config.AuthDetails(config.Con)
	return func(c *gin.Context) {
		userName, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Use Basic Authentication to access this API"})
			return
		}
		if Auth.UserName == userName && Auth.Password == password {
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication Failed"})
			return
		}

	}
}
