// auth-server/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// จำลองฐานข้อมูล
var (
	authCodes = make(map[string]string) // auth_code -> user_id
	tokens    = make(map[string]string) // access_token -> user_id
)

func main() {
	r := gin.Default()

	// 1. Endpoint สำหรับแสดงหน้า Login
	r.GET("/oauth/authorize", func(c *gin.Context) {
		clientID := c.Query("client_id")
		redirectURI := c.Query("redirect_uri")

		// แสดงหน้า login (ในตัวอย่างนี้ใช้ HTML อย่างง่าย)
		html := `
            <form method="post" action="/oauth/approve">
                <input type="hidden" name="client_id" value="%s">
                <input type="hidden" name="redirect_uri" value="%s">
                <input type="text" name="username" placeholder="Username">
                <input type="password" name="password" placeholder="Password">
                <button type="submit">Login & Approve</button>
            </form>
        `
		c.Header("Content-Type", "text/html")
		c.String(200, html, clientID, redirectURI)
	})

	// 2. Endpoint สำหรับรับ Login และออก Authorization Code
	r.POST("/oauth/approve", func(c *gin.Context) {
		// ในความเป็นจริงต้องตรวจสอบ credentials
		username := c.PostForm("username")
		redirectURI := c.PostForm("redirect_uri")

		// สร้าง authorization code
		authCode := uuid.New().String()
		authCodes[authCode] = username

		// redirect กลับไปที่ client พร้อม code
		redirectURL := redirectURI + "?code=" + authCode
		c.Redirect(302, redirectURL)
	})

	// 3. Endpoint สำหรับแลก Authorization Code เป็น Access Token
	r.POST("/oauth/token", func(c *gin.Context) {
		code := c.PostForm("code")

		// ตรวจสอบ code
		userID, exists := authCodes[code]
		if !exists {
			c.JSON(400, gin.H{"error": "invalid_code"})
			return
		}

		// สร้าง access token
		token := uuid.New().String()
		tokens[token] = userID
		delete(authCodes, code) // ใช้ code ได้ครั้งเดียว

		c.JSON(200, gin.H{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})

	r.Run(":8080")
}
