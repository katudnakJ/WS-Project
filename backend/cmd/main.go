package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"onlinecourse/internal/handlers"
)

// Public Key ในรูปแบบ PEM (คัดลอกจาก Keycloak) ----> รอเปลี่ยนเป็น .env แทน
var publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvDi7UB9tY8lSbCMGZ3DuHDucEhar6P9Yt4KaHaC/d7rocrTWteM/i+0mpbaw5FgaIRztrgJ+llDuUhfdAPtea+pPtlQhYrSyZzSjN2okXiVwcXkFgQA2ks1ZMO98aUih6j7GOTuwy28OaP8Z6x3ZG+HBK9E4ZBENK3W64O2/XvdYOcz+nXoaifGdC5T4PlCB/ZM5irKyZafnd5p1DpAG6uBIxe/88WrxwGAmDwuwoNwE8qKwQ0OhU1It1UBsHYSbhKoAGBGNPR+N4TMn8R+PVAu5oDN+2f8Hif+IA1ker2RmO5gDbiWT6Jk6EtrqjAwExmpmnE41ELaM9N0g0bTD8wIDAQAB
-----END PUBLIC KEY-----`

// Global variable สำหรับเก็บ *rsa.PublicKey ที่แปลงแล้ว
var rsaPublicKey *rsa.PublicKey

// Casbin enforcer
var enforcer *casbin.Enforcer

func init() {
	var err error
	// สร้าง Casbin enforcer
	enforcer, err = casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		panic("failed to create casbin enforcer: " + err.Error())
	}

	// Decode PEM และ parse public key เพียงครั้งเดียวตอนเริ่มต้น
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	var ok bool
	rsaPublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		panic("key is not of type *rsa.PublicKey")
	}
}

func main() {
	r := gin.Default()
	r.POST("/register", handlers.Register)
	// เพิ่ม CORS middleware - ใช้ server ของ keycloak
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ข้อ 1 request ได้เฉพาะคนที่สมัคร ----> ยังไม่ทดลอง
	protected := r.Group("/api")

	protected.Use(JWTAuthMiddleware())
	{
		protected.GET("/data", handlers.RequestLogMiddleware(), handlers.GetData)

		//เอาไว้ดึงข้อมูล user
		protected.GET("/userdata", func(c *gin.Context) {
			username := c.GetString("username")
			roles := c.GetStringSlice("roles")
			c.JSON(200, gin.H{
				"username": username,
				"roles":    roles,
			})
		})
	}

	r.Run(":8081")
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Ur right!")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse JWT โดยใช้ rsaPublicKey ที่แปลงแล้วจากขั้นตอน initial setup
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return rsaPublicKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ดึง claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// ตรวจสอบเวลาหมดอายุของโทเค็น
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// ตรวจสอบผู้ออกโทเค็น
		if !claims.VerifyIssuer("http://localhost:8080/realms/auth101", true) {
			c.JSON(401, gin.H{"error": "Invalid token issuer"})
			c.Abort()
			return
		}

		// ดึง username จาก claims
		username, ok := claims["preferred_username"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Username not found in token"})
			c.Abort()
			return
		}

		// ดึง roles จาก realm_access.roles ซึ่งอาจมีหลาย role
		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			c.JSON(401, gin.H{"error": "Roles not found in token"})
			c.Abort()
			return
		}
		rawRoles, ok := realmAccess["roles"].([]interface{})
		if !ok || len(rawRoles) == 0 {
			c.JSON(401, gin.H{"error": "No roles found in token"})
			c.Abort()
			return
		}

		// ดึง role ทั้งหมดจาก payload
		var rolesList []string
		for _, r := range rawRoles {
			if roleStr, ok := r.(string); ok {
				rolesList = append(rolesList, roleStr)
			}
		}

		// ตรวจสอบสิทธิ์ด้วย Casbin: ให้ตรวจสอบว่ามี role ใดที่อนุญาตให้เข้าถึง resource ได้หรือไม่
		resource := c.Request.URL.Path // เช่น /api/data
		action := c.Request.Method     // เช่น GET
		allowed := false
		for _, role := range rolesList {
			permit, err := enforcer.Enforce(role, resource, action)
			if err != nil {
				c.JSON(500, gin.H{"error": "Error checking permission"})
				c.Abort()
				return
			}
			if permit {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(403, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		// ส่ง username และ roles ไปยัง handler
		c.Set("username", username)
		c.Set("roles", rolesList)
		c.Next()
	}
}
