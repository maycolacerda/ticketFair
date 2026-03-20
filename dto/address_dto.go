// dto/address_dto.go
package dto

type CreateAddressRequest struct {
	Street  string `json:"street"   validate:"required"`
	City    string `json:"city"     validate:"required"`
	State   string `json:"state"    validate:"required"`
	Country string `json:"country"  validate:"required,len=2"` // ISO 3166-1 alpha-2 e.g. "BR"
	ZipCode string `json:"zip_code" validate:"required"`
}

type UpdateAddressRequest struct {
	Street  string `json:"street"   validate:"omitempty"`
	City    string `json:"city"     validate:"omitempty"`
	State   string `json:"state"    validate:"omitempty"`
	Country string `json:"country"  validate:"omitempty,len=2"`
	ZipCode string `json:"zip_code" validate:"omitempty"`
}

type AddressResponse struct {
	AddressID string `json:"address_id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	ZipCode   string `json:"zip_code"`
}
