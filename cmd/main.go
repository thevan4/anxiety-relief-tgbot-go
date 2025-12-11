package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mymmrac/telego"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is empty")
	}

	statisticsInMemory := statistic.NewStatistics(nil)
	rateLimiter := rate_limiter.NewRateLimiter(time.Minute, 10)

	botHandler := telegram.MustNewBotHandler(
		ctx,
		botToken,
		nil,
		rateLimiter,
		statisticsInMemory,
		telego.WithDefaultDebugLogger(),
	)

	go botHandler.Start()

	<-stop
	log.Println("Shutting down gracefully...")
	botHandler.Stop()
	cancel()
	log.Println("Bot stopped successfully")
}
