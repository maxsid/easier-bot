package easierbot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendMessage is a simple method of sending simple messages by chatId
func (bot *Bot) SendMessage(chatID int64, message string) {
	log.Printf("Send '%s' to %d chat id\n", message, chatID)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = bot.MessageParseMode
	_, err := bot.API.Send(msg)
	if err != nil {
		log.Panicf("Send Messsage Error: %v\n", err)
	}
}

// panicRecover logs a received panic message. Use with defer.
func panicRecover() {
	if rec := recover(); rec != nil {
		log.Printf("Panic: %v", rec)
	}
}
