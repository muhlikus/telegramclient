package telegramclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	type args struct {
		token              string
		chatID             int
		text               string
		expectedStatusCode int
		expectedBody       string
	}

	tests := []struct {
		name            string
		expectedError   assert.ErrorAssertionFunc
		expectedMessage *Message
		args
	}{
		{
			name: "ValidResponse",
			args: args{
				token:              "SomeToken",
				chatID:             1,
				text:               "Hello",
				expectedStatusCode: http.StatusOK,
				expectedBody:       `{ "ok": true, "result": { "message_id": 123, "text": "Hello", "chat": {"id": 1} } }`,
			},
			expectedMessage: &Message{MessageId: 123, Text: "Hello", Chat: Chat{Id: 1}},
			expectedError:   assert.NoError,
		},
		{
			name: "UnexpectedStatusCode",
			args: args{
				token:              "SomeToken",
				chatID:             1,
				text:               "Hello",
				expectedStatusCode: http.StatusBadRequest, //400
				expectedBody:       `{ "ok": true, "result": { "message_id": 123, "text": "Hello", "chat": {"id": 1} } }`,
			},
			expectedMessage: nil,
			expectedError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "unexpected status code") && assert.ErrorContains(tt, err, "400")
			},
		},
		{
			name: "InvalidResponseJSON",
			args: args{
				token:              "SomeToken",
				chatID:             1,
				text:               "Hello",
				expectedStatusCode: http.StatusOK,
				expectedBody:       `{ "ok" = true, "result": { "message_id": 123, "text": "Hello", "chat": {"id": 1} } }`,
			},
			expectedMessage: nil,
			expectedError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "parsing response JSON")
			},
		},
		{
			name: "ResponseNotOK",
			args: args{
				token:              "SomeToken",
				chatID:             1,
				text:               "Hello",
				expectedStatusCode: http.StatusOK,
				expectedBody:       `{ "ok": false, "description": "Some error" }`,
			},
			expectedMessage: nil,
			expectedError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "response not OK")
			},
		},
		{
			name: "InvalidResultJSON",
			args: args{
				token:              "SomeToken",
				chatID:             1,
				text:               "Hello",
				expectedStatusCode: http.StatusOK,
				expectedBody:       `{ "ok": true, "result": "{ unknown_key = unknown_value }" }`,
			},
			expectedMessage: nil,
			expectedError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "parsing message JSON")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.args.expectedStatusCode)
						_, _ = w.Write([]byte(tt.args.expectedBody))
					}))
			defer server.Close()

			urlMock, _ := url.Parse(server.URL)

			client, err := New(Config{
				Token:        tt.args.token,
				botApiScheme: urlMock.Scheme,
				botApiHost:   urlMock.Host,
			})
			require.NoError(t, err)

			message, err := client.SendMessage(tt.chatID, tt.text)
			tt.expectedError(t, err)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}
