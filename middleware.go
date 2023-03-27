package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Origin") == "http://127.0.0.1:3500" {
			c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:3500")
		} else {
			c.Header("Access-Control-Allow-Origin", "https://wirtschaft.codelix.de")
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
