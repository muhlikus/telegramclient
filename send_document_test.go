package telegramclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendDocument(t *testing.T) {

	type args struct {
		token  string
		chatID int
	}

	type mocks struct {
		StatusCode int
		Body       string
	}

	tests := []struct {
		name            string
		expectedError   assert.ErrorAssertionFunc
		expectedMessage *Message
		args
		mocks
	}{
		{
			name: "ValidResponse",
			args: args{
				token:  "SomeToken",
				chatID: 1,
			},
			mocks: mocks{
				StatusCode: http.StatusOK,
				Body:       `{ "ok": true, "result": { "message_id": 123, "document": {"file_id": "fileIdentificator"} } }`,
			},
			expectedMessage: &Message{MessageId: 123, Document: Document{FileID: "fileIdentificator"}},
			expectedError:   assert.NoError,
		},
	}

	// Create a temporary file to simulate the document to be sent
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the temporary file
	if _, err := tmpFile.Write([]byte("test content")); err != nil {
		t.Fatalf("writing to temp file: %v", err)
	}
	tmpFile.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tt.mocks.StatusCode)
						_, _ = w.Write([]byte(tt.mocks.Body))
					}))
			defer server.Close()

			urlMock, _ := url.Parse(server.URL)

			client, err := New(Config{
				Token:        tt.args.token,
				botApiScheme: urlMock.Scheme,
				botApiHost:   urlMock.Host,
			})
			require.NoError(t, err)

			message, err := client.SendDocument(tt.chatID, tmpFile.Name())
			tt.expectedError(t, err)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}
