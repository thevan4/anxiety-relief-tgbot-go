// Package main is the entry point for the anxiety relief Telegram bot.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mymmrac/telego"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/cleanup"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
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

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST is empty")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("REDIS_PORT is empty")
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	sessionStorage, err := session.NewRedisStorage(redisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	statisticsInMemory := statistic.NewStatistics(nil)

	botHandler := telegram.MustNewBotHandler(
		ctx,
		botToken,
		statisticsInMemory,
		sessionStorage,
		telego.WithDefaultDebugLogger(),
	)

	// Start cleanup worker for background message cleanup
	cleanupWorker := cleanup.NewWorker(ctx, botHandler.GetBot(), sessionStorage)
	go cleanupWorker.Start()

	go botHandler.Start()

	<-stop
	log.Println("Shutting down gracefully...")
	cancel()
	done := make(chan struct{})
	go func() {
		botHandler.Stop()
		close(done)
	}()
	const shutdownTimeout = 5 * time.Second
	select {
	case <-done:
		log.Println("Bot handler stopped")
	case <-time.After(shutdownTimeout):
		log.Printf("WARNING: bot handler did not stop within %v, forcing exit", shutdownTimeout)
	}
	if closeErr := sessionStorage.Close(); closeErr != nil {
		log.Printf("WARNING: failed to close Redis connection: %v", closeErr)
	}
	log.Println("Bot stopped successfully")
}
