package bootstrap

import (
	"context"
	"fmt"
	"net"

	"github.com/Ablyamitov/userbot-core/config"
	tgservice "github.com/Ablyamitov/userbot-core/internal/infrastructure/telegram"
	"github.com/Ablyamitov/userbot-core/internal/interface/grpc"
	"github.com/Ablyamitov/userbot-core/internal/interface/http"
	"github.com/Ablyamitov/userbot-core/internal/usecase/userbot"
	"github.com/Ablyamitov/userbot-core/pkg/tg_client"
	"github.com/gofiber/fiber/v2"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	grpcTransport "google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config) error {
	tgClient := tg_client.NewClient(cfg.Client.AppID, cfg.Client.AppHash)

	svc := tgservice.NewTelegramService(tgClient)
	uc := userbot.NewUserbotUsecase(svc)

	app := fiber.New()
	httpHandler := http.NewHandler(uc)
	httpHandler.RegisterRoutes(app)

	grpcServer := grpcTransport.NewServer()
	grpc.RegisterServices(grpcServer, uc)

	errCh := make(chan error, 3)
	clientRun := make(chan bool)

	go func(chan bool) {
		if err := tgClient.Run(ctx, func(ctx context.Context) error {
			if err := tgClient.Auth().IfNecessary(ctx, auth.NewFlow(
				auth.Constant(cfg.Client.Phone, cfg.Client.Password, auth.CodeAuthenticatorFunc(codePrompt)),
				auth.SendCodeOptions{},
			)); err != nil {
				errCh <- err
				return fmt.Errorf("authorization error: %w", err)
			}
			clientRun <- true

			<-ctx.Done()
			return ctx.Err()
		}); err != nil {
			errCh <- err
			return
		}
		return
	}(clientRun)

	<-clientRun
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
		errCh <- app.Listen(addr)
	}()

	go func() {
		lis, err := net.Listen(cfg.GRPCServer.Network, cfg.GRPCServer.Address)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		_ = app.Shutdown()
		return nil
	case err := <-errCh:
		return err
	}
}

func codePrompt(_ context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Printf("Code sent via: %s\n", sentCode.Type)
	fmt.Print("Enter the verification code from Telegram: ")

	var code string
	_, err := fmt.Scanln(&code)
	if err != nil {
		return "", fmt.Errorf("error entering code: %w", err)
	}
	return code, nil
}
