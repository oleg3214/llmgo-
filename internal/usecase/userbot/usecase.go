package userbot

import (
	"context"

	"github.com/Ablyamitov/userbot-core/internal/domain/entity"
	"github.com/Ablyamitov/userbot-core/internal/domain/service"
)

type Usecase struct {
	userbotService service.UserbotService
}

func NewUserbotUsecase(s service.UserbotService) *Usecase {
	return &Usecase{userbotService: s}
}

func (u *Usecase) GetMessages(ctx context.Context, username string, limit int) ([]entity.Message, error) {
	return u.userbotService.GetMessages(ctx, username, limit)
}

func (u *Usecase) SendMessage(ctx context.Context, username, text string) error {
	return u.userbotService.SendMessage(ctx, username, text)
}
