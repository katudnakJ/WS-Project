package models

// ใส่หลัง json ให้เหมือนใน DB
type Affiliates struct {
	Affiliate_ID       string `json:"affiliate_id"`
	Affiliate_Name     string `json:"affiliate_name"`
	Affiliate_Email    string `json:"affiliate_email"`
	Affiliate_Password string `json:"affiliate_password"`
	Affiliate_Url      string `json:"affiliate_url"`
	Affiliate_APIKey   string `json:"affiliate_api_key"`
}

// type Response_AffRegister struct {
// 	Status  string
// 	Message string
// }

// 	"clicks": 150,
// 	"sales": 45,
// 	"commission_earned": 75.5,
// 	"conversion_rate": 0.3
