package main

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware function definition
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	// var allowedOrigins = []string{
	// 	"http://localhost:5173",
	// 	"http://localhost:5174",
	// 	"http://localhost:3000",
	// 	"https://localhost:5173",
	// 	"https://localhost:3000",
	// }
	// Return the actual middleware handler function
	return func(c *gin.Context) {
		// Get the Origin header from the request
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			// If the origin is allowed, set CORS headers in the response
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PUT")
		}

		// Handle preflight OPTIONS requests by aborting with status 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Call the next handler
		c.Next()
	}
}

// Function to check if a given origin is allowed
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
