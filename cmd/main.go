package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ablyamitov/userbot-core/pkg/bootstrap"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()
	}()
	if err := bootstrap.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
