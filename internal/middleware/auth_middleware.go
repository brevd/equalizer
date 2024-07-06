package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add your authentication logic here

		fmt.Println("middle ware hit")

		c.Next()
	}
}
