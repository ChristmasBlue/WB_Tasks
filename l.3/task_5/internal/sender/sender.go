package sender

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wb-go/wbf/zlog"
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
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("could not connect to telegram api: ", err)
	}

	botApi.Debug = false

	return &TelegramSender{
		botApi: botApi,
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

	zlog.Logger.Info().Msgf("message to telegram user with id: %d was sent successfylly: ", telegramId)
	return nil
}
