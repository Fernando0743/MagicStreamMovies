package middleware

import (
	"net/http"

	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
)

// Gin handler function used to validate incoming access tokens and grant/prohibt access to protected endpoints
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Extract the token from the http cookies
		token, err := utils.GetAccessToken(c)

		//Use Abort because if user is not authorized, we will not continue to call the targeted endpoint
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		//Validate Token and Decode Claims
		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		//Set parameters once authenticated
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)
		//Continue to execute the targeted endpoint
		c.Next()

	}

}
