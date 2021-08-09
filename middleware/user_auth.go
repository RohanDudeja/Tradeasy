package middleware

import (
	"Tradeasy/config"
	"Tradeasy/internal/services/user_management"
	"github.com/gin-gonic/gin"
	"log"
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
		use User SignIn API
		req := user_management.SignInRequest{
			userId,
			password,
		}
		_, err := user_management.UserSignIn(req)
		if err != nil {
			log.Fatalf("%s",err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Credentials didn't matched"})
		}
		c.Next()
	}
}