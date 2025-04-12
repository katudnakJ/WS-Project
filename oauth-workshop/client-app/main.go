// client-app/main.go
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. หน้าแรกที่มีปุ่มเชื่อมต่อ
	r.GET("/", func(c *gin.Context) {
		html := `
            <h1>Welcome to Demo App</h1>
            <a href="/connect">Connect with OAuth</a>
        `
		c.Header("Content-Type", "text/html")
		c.String(200, html)
	})

	// 2. เริ่มต้น OAuth flow
	r.GET("/connect", func(c *gin.Context) {
		authServerURL := "http://localhost:8080/oauth/authorize"
		clientID := "demo-client"
		redirectURI := "http://localhost:3000/callback"

		// ส่ง user ไปที่ auth server
		authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s",
			authServerURL, clientID, redirectURI)
		c.Redirect(302, authURL)
	})

	// 3. รับ Authorization Code และแลกเป็น Token
	r.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")

		// แลก code เป็น token
		// ในความเป็นจริงควรใช้ HTTP client library
		// แต่ในตัวอย่างนี้แสดงแค่แนวคิด
		c.String(200, "Got code: "+code)
	})

	r.Run(":3000")
}
