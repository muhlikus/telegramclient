package telegramclient

import "encoding/json"

type Response struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int      `json:"message_id"`
	Date      int      `json:"date"`
	Text      string   `json:"text,omitempty"`
	Chat      Chat     `json:"chat"`
	Document  Document `json:"document,omitempty"`
}

type Chat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"` //Type of the chat, can be either “private”, “group”, “supergroup” or “channel”
}

type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
}
