package models

type Merchant struct {
	MerchantID  string `json:"merchant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
