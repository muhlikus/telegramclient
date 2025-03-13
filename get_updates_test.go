package telegramclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUpdates(t *testing.T) {
	type args struct {
		token          string
		httpStatusCode int
		httpBody       string
	}

	tests := []struct {
		name          string
		errorExpected assert.ErrorAssertionFunc
		args
	}{
		{
			name: "ValidResponse",
			args: args{
				token:          "SomeToken",
				httpStatusCode: http.StatusOK,
				httpBody:       `{ "ok": true, "description": "OK", "result": [{ "update_id": 1 }] }`,
			},
			errorExpected: assert.NoError,
		},
		{
			name: "EnxpectedStatusCode",
			args: args{
				token:          "SomeToken",
				httpStatusCode: http.StatusBadRequest,
				httpBody:       `{ "ok": true, "description": "Not OK", "result": [{ "update_id": 1 }] }`,
			},
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, errEmptyToken)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.args.httpStatusCode)
						_, _ = w.Write([]byte(tt.args.httpBody))
					}))
			defer server.Close()

			urlMock, _ := url.Parse(server.URL)

			client, err := New(Config{
				Token:        tt.args.token,
				botApiScheme: urlMock.Scheme,
				botApiHost:   urlMock.Host,
			})
			require.NoError(t, err)

			updates, err := client.GetUpdates()
			tt.errorExpected(t, err)
			_ = updates
		})
	}
}
