package middleware

import (
	"Tradeasy/config"
	"Tradeasy/internal/services/user_management"
	"fmt"
	"github.com/gin-gonic/gin"
)

func UserBasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		userId, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Use Basic Authentication to access this API"})
			return
		}
		if userId == "" || password == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Enter details in Basic Authentication"})
			return
		}
		// use User SignIn API
		req := user_management.SignInRequest{
			userId,
			password,
		}
		_, err := user_management.UserSignIn(req)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatusJSON(401, gin.H{"error": "Credentials didn't matched"})
		}
		c.Next()
	}
}

func TradingAuth() gin.HandlerFunc {
	databaseName, databasePass := config.TradingDetails(config.BuildConfig())
	return func(c *gin.Context) {
		user, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Use Basic Authentication to access this API"})
			return
		}
		if user == databaseName && password == databasePass {
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication Failed"})
			return
		}

	}
}
