package telegramclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func (c *Client) SendDocument(chatID int, filePath string) (*Message, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reqURL := url.URL{
		Scheme: c.cfg.botApiScheme,
		Host:   c.cfg.botApiHost,
		Path:   path.Join(c.cfg.botApiPath, sendDocumentMethod),
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField("chat_id", strconv.Itoa(chatID))

	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}

	// we need to close writer for:
	// "finishes the multipart message and writes the trailing boundary end line to the output."
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("closing writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("parsing response JSON: %d", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("response not OK: %s", response.Description)
	}

	var message Message
	err = json.Unmarshal(response.Result, &message)
	if err != nil {
		return nil, fmt.Errorf("parsing message JSON: %d", err)
	}
	return &message, nil
}
