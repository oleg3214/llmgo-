package request

type GetChannelRequest struct {
	Username string `json:"username,omitempty"`
}

type MessagesRequest struct {
	Username string `json:"username,omitempty"`
	Limit    int32  `json:"limit,omitempty"`
}

type SendMessageRequest struct {
	Username string `json:"username,omitempty"`
	Text     string `json:"text,omitempty"`
}
