package main

import (
	"github.com/boltdb/bolt"
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/config"
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/repository"
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/repository/boltdb"
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/server"
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("config init: ", err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("bot init: ", err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal("pocket init: ", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("db init: ", err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(
		bot,
		pocketClient,
		cfg.AuthServerURL,
		tokenRepository,
		cfg.Messages,
	)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatalf("start err: %s", err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
