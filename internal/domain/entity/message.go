package entity

type Message struct {
	ID        int
	Text      string
	Date      int
	ReplyToID *int            `json:"reply_to_id,omitempty"`
	Entities  []MessageEntity `json:"entities,omitempty"`
}

type MessageEntity struct {
	Url string `json:"url"`
}
