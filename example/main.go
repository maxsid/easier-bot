package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	easierbot "github.com/maxsid/easier-bot"
	"github.com/maxsid/easier-bot/contentType"
	"os"
)

func main() {
	config := NewConfig(true)
	bot := easierbot.NewBotViaWebhook(config.token, config.webhookSite, config.listenAddress, config.isDebug)
	// or
	// bot := easierbot.NewBot(config.token, config.isDebug)
	bot.Handlers.AddSeveralCommandsHandler([]string{"hello", "start", "hi"}, func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Hello World!")
		msg.CommandWithAt()
	})
	bot.Handlers.AddCommandHandler("bye", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Ok. Bye.")
	})
	bot.Handlers.AddRegexpHandler(".*ping.*", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "I'll receive any text contains a regex expression")
	})
	bot.Handlers.SetContentHandler(contentType.Document, func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "I'll receive any your document")
	})
	bot.Handlers.SetContentHandler(contentType.Any, func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "I'll receive any your not text content")
	})
	bot.Handlers.SetDefaultHandler(func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "I answer to a message in any case")
	})
	bot.RunBotServer()
}

type Config struct {
	token         string
	isDebug       bool
	webhookSite   string
	listenAddress string
}

func NewConfig(isDebug bool) *Config {
	return &Config{token: os.Getenv("TOKEN"), isDebug: isDebug,
		webhookSite: os.Getenv("WEBHOOK_SITE"), listenAddress: os.Getenv("LISTEN_ADDRESS")}
}
