package handlers

// pass รอแก้เป็น hash แล้วเก็บไว้
import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"onlinecourse/internal/models"

	"github.com/gin-gonic/gin"
)

// Gen APIKey
func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func Register(c *gin.Context) {
	var aff models.Affiliates
	if err := c.ShouldBindJSON(&aff); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		log.Println("error", err.Error())
		return
	}

	// aff.Affiliate_APIKey = generateAPIKey()

	// _, err := database.DB.Exec("Insert into affiliates (affiliate_name,affiliate_email,affiliate_password,affiliate_url, affiliate_api_key) value ($1, $2, $3, $4, $5)", aff.Affiliate_Name, aff.Affiliate_Email, aff.Affiliate_Password, aff.Affiliate_Url, aff.Affiliate_APIKey)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register client"})
	// 	log.Println("error", err)
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"message": "You are registered, Thanks to join Us!", "api_key": aff.Affiliate_APIKey})

	// Test without DB
	Affiliate_APIKey := generateAPIKey()
	c.JSON(http.StatusOK, gin.H{"message": "You are registered, Thanks to join Us!", "api_key": Affiliate_APIKey})
}
