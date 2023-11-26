package models

type CompanyUser struct {
	CompanyUserID       int    `json:"company_user_id"`
	CompanyID           int    `json:"company_id"`
	CompanyUserLogin    string `json:"company_user_login"`
	CompanyUserPassword string `json:"company_user_password"`
	ComanyUserRole      string `json:"company_user_role"`
}
