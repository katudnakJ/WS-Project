package models

// รวม request GET ต่างๆ

type Click_logs struct {
	Affiliate_ID string `json:"affiliate_id"`
	Course_ID    string `json:"course_id"`
	Action       string `json:"action"`
	Clicks       int    `json:"clicks"`
	Timestamp    string `json:"click_date"`
}

// log การ request

type RequestLog struct {
	AffiliateID string `json:"affiliate_id"`
	Action      string `json:"action"`
	Parameter   string `json:"parameter"` // คำค้นหา
	Timestamp   string `json:"timestamp"` // วันและเวลา
}
