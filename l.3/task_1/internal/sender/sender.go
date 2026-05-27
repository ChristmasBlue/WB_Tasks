package sender

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSender struct {
	botApi  *tgbotapi.BotAPI
	enabled bool
}

func New() *TelegramSender {
	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Println("BOT_TOKEN is empty, Telegram sender is disabled")
		return &TelegramSender{
			botApi:  nil,
			enabled: false,
		}
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("could not connect to telegram api: ", err)
		return &TelegramSender{
			botApi:  nil,
			enabled: false,
		}
	}

	botApi.Debug = false

	return &TelegramSender{
		botApi:  botApi,
		enabled: true,
	}
}

func (t *TelegramSender) SendToTelegram(telegramId int, text string) error {
	if !t.enabled {
		log.Println("Telegram sender is disabled")
		return nil
	}

	msg := tgbotapi.NewMessage(int64(telegramId), text)
	_, err := t.botApi.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send message to telegram user: %s", err.Error())
	}

	return nil
}
