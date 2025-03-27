package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message, email, password)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, email, password string) {
	text := message.Text
	chatID := message.Chat.ID

	switch text {
	case "/start":
		bot.Send(tgbotapi.NewMessage(chatID, "Привет! Напиши /code, чтобы прочитать последнее письмо."))
	case "/code":
		lastMsg, err := GetLastEmail(email, password)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Ошибка: "+err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Последнее письмо:\n"+lastMsg))
		}
	default:
		bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда. Попробуй /start"))
	}
}
