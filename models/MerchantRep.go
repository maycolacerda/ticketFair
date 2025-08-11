package models

type MerchantRep struct {
	MerchantRepID string `json:"id"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	MerchantID    string `json:"merchant_id"`
}
