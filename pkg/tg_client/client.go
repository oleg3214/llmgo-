package tg_client

import (
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

const SessionPath = "session.json"

func NewClient(AppID int, AppHash string) *telegram.Client {
	return telegram.NewClient(AppID, AppHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: SessionPath},
		NoUpdates:      false,
	})
}
