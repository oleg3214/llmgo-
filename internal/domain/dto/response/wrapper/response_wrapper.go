package wrapper

import "github.com/Ablyamitov/userbot-core/internal/domain/entity"

type ResponseWrapper struct {
	Status bool `json:"status"`
	Code   int  `json:"code"`
	Data   any  `json:"data"`
	Error  any  `json:"error"`
}

type GetChannelResponse struct {
	ChatId     int64  `json:"chat_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Title      string `json:"title,omitempty"`
	AccessHash int64  `json:"access_hash,omitempty"`
}

type MessagesResponse struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	MessageId int64           `json:"message_id,omitempty"`
	Text      string          `json:"text,omitempty"`
	Date      int64           `json:"date,omitempty"`
	ReplyToID *int            `json:"reply_to_id,omitempty"`
	Entities  []MessageEntity `json:"entities,omitempty"`
}

type MessageEntity struct {
	Url string `json:"url"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

func FromEntity(m *entity.Message) *Message {
	if m == nil {
		return nil
	}
	entitiesDTO := make([]MessageEntity, len(m.Entities))
	for i, e := range m.Entities {
		entitiesDTO[i] = MessageEntity{
			Url: e.Url,
		}
	}
	return &Message{
		MessageId: int64(m.ID),
		Text:      m.Text,
		Date:      int64(m.Date),
		ReplyToID: m.ReplyToID,
		Entities:  entitiesDTO,
	}
}
