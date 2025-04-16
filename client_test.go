package telegramclient

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewClient tests the New function
func TestNewClient(t *testing.T) {
	type args struct {
		token string
	}

	//create table tests
	tests := []struct {
		name          string
		errorExpected assert.ErrorAssertionFunc
		args
	}{
		{
			name: "ValidToken",
			args: args{
				token: "SomeToken",
			},
			errorExpected: assert.NoError,
		},
		{
			name: "EmptyToken",
			args: args{
				token: "",
			},
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "parsing config")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", tt.token)
			_, err := New(Config{Token: tt.token})
			tt.errorExpected(t, err)
		})
	}

}
