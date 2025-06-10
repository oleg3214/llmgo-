package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Ablyamitov/userbot-core/internal/domain/entity"
	"github.com/Ablyamitov/userbot-core/internal/domain/service"
	"github.com/Ablyamitov/userbot-core/internal/infrastructure/cache"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type TgService struct {
	client    *telegram.Client
	processor *RequestProcessor
}

func NewTelegramService(c *telegram.Client) service.UserbotService {
	processor := NewRequestProcessor(100, 200*time.Millisecond)
	return &TgService{client: c, processor: processor}
}
func (s *TgService) ResolvePeer(ctx context.Context, username string) (*entity.Peer, error) {
	resolved, err := s.client.API().ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("resolve username: %w", err)
	}

	// Check Users
	for _, u := range resolved.Users {
		if user, ok := u.(*tg.User); ok {
			_, err := s.client.API().UsersGetFullUser(ctx, &tg.InputUser{
				UserID:     user.ID,
				AccessHash: user.AccessHash,
			})
			if err != nil {
				return nil, fmt.Errorf("get full user: %w", err)
			}
			tgPeer := tg.InputPeerUser{
				UserID:     user.ID,
				AccessHash: user.AccessHash,
			}
			peer := &entity.Peer{
				Username:   username,
				PeerObject: tgPeer,
				Type:       entity.User,
			}
			cache.ChatInfoMap.Store(username, peer)
			return peer, nil
		}
	}

	// Check Channels
	for _, ch := range resolved.Chats {
		if channel, ok := ch.(*tg.Channel); ok {
			_, err = s.client.API().ChannelsGetFullChannel(ctx, &tg.InputChannel{
				ChannelID:  channel.ID,
				AccessHash: channel.AccessHash,
			})

			if err != nil {
				return nil, fmt.Errorf("get full channel: %w", err)
			}

			tgPeer := tg.InputPeerChannel{
				ChannelID:  channel.ID,
				AccessHash: channel.AccessHash,
			}
			peer := &entity.Peer{
				Username:   username,
				PeerObject: tgPeer,
				Type:       entity.Channel,
			}

			cache.ChatInfoMap.Store(username, peer)
			return peer, nil
		}
		if chat, ok := ch.(*tg.Chat); ok {
			tgPeer := tg.InputPeerChat{
				ChatID: chat.ID,
			}
			peer := &entity.Peer{
				Username:   username,
				PeerObject: tgPeer,
				Type:       entity.Chat,
			}
			cache.ChatInfoMap.Store(username, peer)
			return &entity.Peer{
				Username:   username,
				PeerObject: peer,
				Type:       entity.Chat,
			}, nil
		}
	}
	return nil, fmt.Errorf("no chats: %w", err)
}

func (s *TgService) GetMessages(ctx context.Context, username string, limit int) ([]entity.Message, error) {

	type resultStruct struct {
		messages []entity.Message
		err      error
	}
	result := make(chan resultStruct, 1)
	s.processor.Enqueue(TgRequest{
		Do: func(ctx context.Context) error {
			var err error
			peer, ok := cache.ChatInfoMap.Load(username)
			if !ok {
				peer, err = s.ResolvePeer(ctx, username)
				if err != nil {
					log.Printf("error resolving: %v", err)
					result <- resultStruct{nil, fmt.Errorf("error resolving channel username")}
					return nil
				}
			}

			var messagesClass tg.MessagesMessagesClass
			if channelPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerChannel); ok {
				messagesClass, err = s.client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  &channelPeer,
					Limit: limit,
				})
				if err != nil {
					result <- resultStruct{nil, err}
					return nil
				}
			}

			if userPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerUser); ok {
				messagesClass, err = s.client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  &userPeer,
					Limit: limit,
				})
				if err != nil {
					result <- resultStruct{nil, err}
					return nil
				}
			}

			if chatPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerChat); ok {
				messagesClass, err = s.client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  &chatPeer,
					Limit: limit,
				})
				if err != nil {
					result <- resultStruct{nil, err}
					return nil
				}
			}

			if messagesClass == nil {
				result <- resultStruct{nil, fmt.Errorf("message class not found")}
				return nil
			}

			modifiedMessages, ok := messagesClass.AsModified()
			if !ok {
				result <- resultStruct{nil, fmt.Errorf("error to modify messages")}
				return nil
			}

			var messages []entity.Message
			for _, msg := range modifiedMessages.GetMessages() {
				if message, ok := msg.(*tg.Message); ok {

					/**/

					messages = append(messages, *MapTelegramMessage(message))

					/**/
					//messages = append(messages, entity.Message{
					//	ID:   message.ID,
					//	Text: message.Message,
					//	Date: message.Date,
					//})
				}
			}
			result <- resultStruct{messages, nil}
			return nil
		},
		Ctx:    ctx,
		Result: make(chan error, 1),
	})
	res := <-result
	return res.messages, res.err
}

func (s *TgService) SendMessage(ctx context.Context, username string, text string) error {

	result := make(chan error, 1)
	s.processor.Enqueue(TgRequest{
		Do: func(ctx context.Context) error {
			randomID := rand.Int63()
			var err error
			peer, ok := cache.ChatInfoMap.Load(username)
			if !ok {
				peer, err = s.ResolvePeer(ctx, username)
				if err != nil {
					log.Printf("error resolving: %v", err)
					return fmt.Errorf("error resolving channel username")
				}
			}

			if channelPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerChannel); ok {
				_, err := s.client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
					Message:  text,
					Peer:     &channelPeer,
					RandomID: randomID,
				})
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			}

			if userPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerUser); ok {
				_, err := s.client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
					Message:  text,
					Peer:     &userPeer,
					RandomID: randomID,
				})
				if err != nil {
					return err
				}
				return nil
			}

			if chatPeer, ok := peer.(*entity.Peer).PeerObject.(tg.InputPeerChat); ok {
				_, err := s.client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
					Message:  text,
					Peer:     &chatPeer,
					RandomID: randomID,
				})
				if err != nil {
					return err
				}
				return nil
			}
			return fmt.Errorf("chat not found")
		},
		Ctx:    ctx,
		Result: result,
	})
	return <-result

}
