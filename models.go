package telegramclient

type Response struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      []byte `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}

//type Updates []Updates

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message,omitempty"`
}

type Message struct {
	MessageId int  `json:"message_id"`
	Date      int  `json:"date"`
	Chat      Chat `json:"chat"`
}

type Chat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"` //Type of the chat, can be either “private”, “group”, “supergroup” or “channel”
}
