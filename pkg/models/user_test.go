package models

import "testing"

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		want     bool
	}{
		{"valid_username", true},
		{"valid-username", true},
		{"valid@username", true},
		{"INVALID_USERNAME", true},
		{"1234567890", true},
		{"invalid username", false},
		{"invalid/username", false},
		{"invalid.username", false},
		{"invalid,username", false},
		{"invalid;username", false},
		{"invalid:username", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := ValidateUsername(tt.username); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
