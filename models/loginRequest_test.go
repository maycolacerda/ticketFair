package models

import "testing"

func TestLoginRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request LoginRequest
		wantErr bool
	}{
		{
			name: "Valid Request",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			request: LoginRequest{
				Email:    "invalid-email",
				Password: "Password123!",
			},
			wantErr: true,
		},
		{
			name: "Missing Password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.Validate()
			if (len(got) > 0) != tt.wantErr {
				t.Errorf("LoginRequest.Validate() = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}
