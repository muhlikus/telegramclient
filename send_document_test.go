package telegramclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendDocument(t *testing.T) {

	type args struct {
		token       string
		chatID      int
		fileName    string
		fileContent string
	}

	type response struct {
		StatusCode int
		Body       string
	}

	tests := []struct {
		name            string
		expectedError   assert.ErrorAssertionFunc
		expectedMessage *Message
		args
		response
	}{
		{
			name: "OK",
			args: args{
				token:       "SomeToken",
				chatID:      1,
				fileName:    "testfile.txt",
				fileContent: "test content",
			},
			response: response{
				StatusCode: http.StatusOK,
				Body:       `{ "ok": true, "result": { "message_id": 123, "document": {"file_id": "fileIdentificator"} } }`,
			},
			expectedMessage: &Message{MessageId: 123, Document: Document{FileID: "fileIdentificator"}},
			expectedError:   assert.NoError,
		},
		{
			name: "UnexpectedStatusCode",
			args: args{
				token:       "SomeToken",
				chatID:      1,
				fileName:    "testfile.txt",
				fileContent: "test content",
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body:       `{ "ok": true, "result": { "message_id": 123, "document": {"file_id": "fileIdentificator"} } }`,
			},
			expectedMessage: nil,
			expectedError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "unexpected status code") && assert.ErrorContains(tt, err, "400")
			},
		}}

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

			fileBuffer := bytes.NewBufferString(tt.fileContent)

			message, err := client.SendDocument(tt.chatID, tt.fileName, fileBuffer)
			tt.expectedError(t, err)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}
