// models/user_test.go
package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// O validador deve ser inicializado uma vez, geralmente no init() ou injetado.
// Para testes, podemos tê-lo aqui.
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Este método ValidateUser deve estar em models/user.go
// func (u *User) ValidateUser() error {
//     return validate.Struct(u)
// }

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "Valid User",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
				Password: "Paww0rd1@",
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			user: User{
				Username: "invaliduser",
				Email:    "invalidemail",
				Password: "Paww0rd1@",
			},
			wantErr: true,
		},
		{
			name: "Short Password",
			user: User{
				Username: "shortuser",
				Email:    "short@example.com",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "Missing Symbol",
			user: User{
				Username: "missinguser",
				Email:    "missing@example.com",
				Password: "Paww0rd1",
			},
			wantErr: true,
		},
		{
			name: "Missing Uppercase",
			user: User{
				Username: "missinguser",
				Email:    "missing@example.com",
				Password: "password1@",
			},
			wantErr: true,
		},
		{
			name: "Missing Lowercase",
			user: User{
				Username: "missinguser",
				Email:    "missing@example.com",
				Password: "PASSWORD1@",
			},
			wantErr: true,
		},
		{
			name: "Missing Number",
			user: User{
				Username: "missinguser",
				Email:    "missing@example.com",
				Password: "Password@",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
