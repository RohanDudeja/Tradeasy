package middleware

import (
	"Tradeasy/internal/services/user_management"
	"github.com/gin-gonic/gin"
	"log"
)

//UserBasicAuth is just checking if account exists for user
func UserBasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		userID, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Use Basic Authentication to access this API"})
			return
		}
		if userID == "" || password == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Enter details in Basic Authentication"})
			return
		}
		//use User SignIn API
		req := user_management.SignInRequest{
			UserId:   userID,
			Password: password,
		}
		_, err := user_management.UserSignIn(req)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Credentials didn't matched"})
		}

		c.Next()
	}
}

//UserVerificationAuth checks if user had added PAN details or not
func UserVerificationAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		userID, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Use Basic Authentication to access this API"})
			return
		}
		if userID == "" || password == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Enter details in Basic Authentication"})
			return
		}
		//use User SignIn API
		req := user_management.SignInRequest{
			UserId:   userID,
			Password: password,
		}
		_, err := user_management.UserSignIn(req)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Credentials didn't matched"})
		}
		err = user_management.CheckUserDetails(userID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": err})
		}
		c.Next()
	}
}
