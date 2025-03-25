package telegramclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
)

func (c *Client) SendDocument(chatID int, fileName string, fileBuff *bytes.Buffer) (*Message, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField("chat_id", strconv.Itoa(chatID))

	part, err := writer.CreateFormFile("document", filepath.Base(fileName))
	if err != nil {
		return nil, fmt.Errorf("form file creation: %w", err)
	}

	_, err = part.Write(fileBuff.Bytes())
	if err != nil {
		return nil, fmt.Errorf("write file to mime object: %w", err)
	}

	// we need to close writer for:
	// "finishes the multipart message and writes the trailing boundary end line to the output."
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("closing writer: %w", err)
	}

	reqURL := url.URL{
		Scheme: c.cfg.botApiScheme,
		Host:   c.cfg.botApiHost,
		Path:   path.Join(c.cfg.botApiPath, sendDocumentMethod),
	}

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("parsing response JSON: %d", err)
	}

	// TODO проверить соответсвие response.OK и resp.StatusCode
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
