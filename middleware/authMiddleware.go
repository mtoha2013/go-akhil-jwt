package middleware

import (
	"fmt"
	helper "github.com/mtoha/akhil/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Token Kosong")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		log.Print("claims:", claims)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user type", claims.User_type)
		c.Next()
	}
}
