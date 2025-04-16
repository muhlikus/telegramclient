package telegramclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewClient tests the New function
func TestNewClient(t *testing.T) {
	type args struct {
		token        string
		botApiScheme string
		botApiHost   string
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
				token:        "SomeToken",
				botApiScheme: "https",
				botApiHost:   "api.telegram.org",
			},
			errorExpected: assert.NoError,
		},
		{
			name: "EmptyToken",
			args: args{
				token: "",
			},
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, errEmptyToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(Config{Token: tt.token, BotApiScheme: tt.botApiScheme, BotApiHost: tt.botApiHost})
			tt.errorExpected(t, err)
		})
	}

}
