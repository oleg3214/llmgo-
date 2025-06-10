package service

import (
	"context"

	"github.com/Ablyamitov/userbot-core/internal/domain/entity"
)

type UserbotService interface {
	GetMessages(ctx context.Context, username string, limit int) ([]entity.Message, error)
	SendMessage(ctx context.Context, username string, text string) error
	ResolvePeer(ctx context.Context, username string) (*entity.Peer, error)
}
