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
		token string
	}

	type response struct {
		StatusCode int
		Body       string
	}

	tests := []struct {
		name            string
		errorExpected   assert.ErrorAssertionFunc
		expectedUpdates []Update
		args
		response
	}{
		{
			name: "Success",
			args: args{
				token: "SomeToken",
			},
			response: response{
				StatusCode: http.StatusOK,
				Body:       `{ "ok": true, "description": "OK", "result": [{ "update_id": 1 }] }`,
			},
			expectedUpdates: []Update{{UpdateID: 1}},
			errorExpected:   assert.NoError,
		},
		{
			name: "Error - UnexpectedStatusCode",
			args: args{
				token: "SomeToken",
			},
			response: response{
				StatusCode: http.StatusBadRequest, //400
				Body:       `{ "ok": true, "description": "Not OK", "result": [{ "update_id": 1 }] }`,
			},
			expectedUpdates: nil,
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "unexpected status code") && assert.ErrorContains(tt, err, "400")
			},
		},
		{
			name: "Error - InvalidResponseJSON",
			args: args{
				token: "SomeToken",
			},
			response: response{
				StatusCode: http.StatusOK,
				Body:       `{ "ok" = true, "description" = "Not OK", "result": [{ "updateId": 1 }] }`,
			},
			expectedUpdates: nil,
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "parsing response JSON")
			},
		},
		{
			name: "Error - ResponseNotOK",
			args: args{
				token: "SomeToken",
			},
			response: response{
				StatusCode: http.StatusOK,
				Body:       `{ "ok": false, "description": "Not OK"}`,
			},
			expectedUpdates: nil,
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "response not OK")
			},
		},
		{
			name: "Error - InvalidResultJSON",
			args: args{
				token: "SomeToken",
			},
			response: response{
				StatusCode: http.StatusOK,
				Body:       `{ "ok": true, "description": "OK", "result": "updateId := 1" }] }`,
			},
			expectedUpdates: nil,
			errorExpected: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "parsing updates JSON")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.response.StatusCode)
						_, _ = w.Write([]byte(tt.response.Body))
					}))
			defer server.Close()

			urlMock, _ := url.Parse(server.URL)

			client, err := New(Config{
				Token:        tt.args.token,
				BotApiScheme: urlMock.Scheme,
				BotApiHost:   urlMock.Host,
			})
			require.NoError(t, err)

			updates, err := client.GetUpdates()
			tt.errorExpected(t, err)
			assert.Equal(t, tt.expectedUpdates, updates)
		})
	}
}
