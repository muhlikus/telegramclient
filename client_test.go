package telegramclient

import (
	"testing"
)

func TestNewClientOK(t *testing.T) {

	_, err := New(Config{Token: "SomeToken"})

	if err != nil {
		t.Errorf("Unexpected error = %v", err)
	}
}

func TestNewClientEmptyToken(t *testing.T) {

	_, err := New(Config{Token: ""})

	if err == nil {
		t.Errorf("Expected error \"the token cannot be empty\"")
	}
}

func TestNewClient(t *testing.T) {

	//create table tests
	tests := []struct {
		name          string
		token         string
		errorExpected bool
	}{
		{name: "ValidToken", token: "SomeToken", errorExpected: false},
		{name: "EmptyToken", token: "", errorExpected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(Config{Token: tt.token})

			if tt.errorExpected && err == nil {
				t.Errorf("Expected error \"the token cannot be empty\"")
			}
			if !tt.errorExpected && err != nil {
				t.Errorf("Unexpected error = %v", err)
			}
		})
	}

}
