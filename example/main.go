package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	easierbot "github.com/maxsid/easier-bot"
	"os"
)

func main() {
	config := NewConfig(true)
	bot := easierbot.NewBotViaWebhook(config.token, config.webhookSite, config.listenAddress, config.isDebug)
	// or
	// bot := easierbot.NewBot(config.token, config.isDebug)
	bot.Handlers.AddSeveralCommandsHandler([]string{"hello", "start", "hi"}, func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Hello World!")
	})
	bot.Handlers.AddCommandHandler("bye", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "Ok. Bye.")
	})
	bot.Handlers.AddRegexpHandler(".*ping.*", func(bot *easierbot.Bot, msg *tgbotapi.Message) {
		bot.SendMessage(msg.Chat.ID, "pong")
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
