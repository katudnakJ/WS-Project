package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	// ดึง Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no_token"})
		c.Abort()
		return
	}

	// ตรวจสอบรูปแบบ "Bearer <token>"
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token_format"})
		c.Abort()
		return
	}

	// แยก token ออกมา (ตัด Bearer ออก)
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty_token"})
		c.Abort()
		return
	}

	// เก็บ token ไว้ใช้ใน handler ถ้าต้องการ
	c.Set("token", token)
	c.Next()
}

func getUserData(c *gin.Context) {
	// ดึง token ที่เก็บไว้มาใช้ (ถ้าต้องการ)
	token, _ := c.Get("token")

	c.JSON(http.StatusOK, gin.H{
		"message": "This is protected data",
		"token":   token,
	})
}

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return	
		}

		c.Next()
	})

	// กำหนด route และใช้ middleware
	r.GET("/api/user-data", authMiddleware, getUserData)

	// สร้าง route สำหรับทดสอบว่า server ทำงาน
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	fmt.Println("Server is running on :8081")
	r.Run(":8081")
}
