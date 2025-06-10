package telegram

import (
	"github.com/Ablyamitov/userbot-core/internal/domain/entity"
	"github.com/gotd/td/tg"
)

func MapTelegramMessage(msg *tg.Message) *entity.Message {
	if msg == nil {
		return nil
	}

	var replyToID *int
	if replyTo, ok := msg.ReplyTo.(*tg.MessageReplyHeader); ok {
		replyToID = &replyTo.ReplyToMsgID
	}

	entitiesMapped := make([]entity.MessageEntity, 0, len(msg.Entities))
	for _, e := range msg.Entities {
		switch ent := e.(type) {
		case *tg.MessageEntityTextURL:
			entitiesMapped = append(entitiesMapped, entity.MessageEntity{
				Url: ent.URL,
			})
		}
	}

	return &entity.Message{
		ID:        msg.ID,
		Text:      msg.Message,
		Date:      msg.Date,
		ReplyToID: replyToID,
		Entities:  entitiesMapped,
	}
}
