package api

import (
	shared "expence_management/Shared"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(c *gin.Context) {
	//extracting token from Authorization request header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" || authToken[1] == "null" {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// verify the token validity
	if err != nil || !token.Valid {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}
	// Extracting claims section from token to be checked later
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}
	// Check if the token is not expired yet
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("token", tokenString)
	c.Set("role", claims["role"]) //fixme: security issue
	c.Set("username", claims["username"])
	c.Next()
}

func IsAdmin(c *gin.Context) {

	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": shared.ErrMissedRoleProperty})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}

	// check if role is admin
	if role != shared.Roles[0] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": shared.ErrOnlyAdmin})
		c.AbortWithStatus(http.StatusUnauthorized)
		// c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Next()
}
